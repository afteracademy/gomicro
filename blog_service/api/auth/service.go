package auth

import (
	"github.com/afteracademy/gomicro/blog-service/api/auth/message"
	"github.com/afteracademy/goserve/v2/micro"
	"github.com/google/uuid"
)

const NATS_TOPIC_AUTH = "auth.authentication"
const NATS_TOPIC_AUTHZ = "auth.authorization"
const NATS_TOPIC_USERPROFILE = "auth.profile.user"

type Service interface {
	Authenticate(token string) (*message.User, error)
	Authorize(user *message.User, roles ...string) error
	FindUserPublicProfile(userId uuid.UUID) (*message.User, error)
}

type service struct {
	natsClient micro.NatsClient
}

func NewService(natsClient micro.NatsClient) Service {
	return &service{
		natsClient: natsClient,
	}
}

func (s *service) Authenticate(token string) (*message.User, error) {
	msg := message.NewText(token)
	return micro.RequestNats[message.Text, message.User](s.natsClient, NATS_TOPIC_AUTH, msg)
}

func (s *service) Authorize(user *message.User, roles ...string) error {
	msg := message.NewUserRole(user, roles...)
	_, err := micro.RequestNats[message.UserRole, message.User](s.natsClient, NATS_TOPIC_AUTHZ, msg)
	return err
}

func (s *service) FindUserPublicProfile(userId uuid.UUID) (*message.User, error) {
	msg := message.NewText(userId.String())
	return micro.RequestNats[message.Text, message.User](s.natsClient, NATS_TOPIC_USERPROFILE, msg)
}
