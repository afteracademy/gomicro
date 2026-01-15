package auth

import (
	"github.com/afteracademy/gomicro/blog-service/api/auth/message"
	"github.com/afteracademy/goserve/arch/micro"
	"github.com/afteracademy/goserve/arch/network"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	Authenticate(token string) (*message.User, error)
	Authorize(user *message.User, roles ...string) error
	FindUserPublicProfile(userId primitive.ObjectID) (*message.User, error)
}

type service struct {
	network.BaseService
	authRequestBuilder  micro.RequestBuilder[message.User]
	authzRequestBuilder micro.RequestBuilder[message.User]
	userRequestBuilder  micro.RequestBuilder[message.User]
}

func NewService(natsClient micro.NatsClient) Service {
	return &service{
		BaseService:         network.NewBaseService(),
		authRequestBuilder:  micro.NewRequestBuilder[message.User](natsClient, "auth.authentication"),
		authzRequestBuilder: micro.NewRequestBuilder[message.User](natsClient, "auth.authorization"),
		userRequestBuilder:  micro.NewRequestBuilder[message.User](natsClient, "auth.profile.user"),
	}
}

func (s *service) Authenticate(token string) (*message.User, error) {
	msg := message.NewText(token)
	return s.authRequestBuilder.Request(msg).Nats()
}

func (s *service) Authorize(user *message.User, roles ...string) error {
	msg := message.NewUserRole(user, roles...)
	_, err := s.authzRequestBuilder.Request(msg).Nats()
	return err
}

func (s *service) FindUserPublicProfile(userId primitive.ObjectID) (*message.User, error) {
	msg := message.NewText(userId.Hex())
	return s.userRequestBuilder.Request(msg).Nats()
}
