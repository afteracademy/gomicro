package dto

import (
	"time"

	"github.com/afteracademy/gomicro/blog-service/api/auth/message"
	"github.com/afteracademy/gomicro/blog-service/api/author/dto"
	"github.com/afteracademy/gomicro/blog-service/api/blog/model"
	"github.com/afteracademy/goserve/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PublicBlog struct {
	ID          primitive.ObjectID `json:"_id" binding:"required" validate:"required"`
	Title       string             `json:"title" validate:"required,min=3,max=500"`
	Description string             `json:"description" validate:"required,min=3,max=2000"`
	Text        string             `json:"text" validate:"required,max=50000"`
	Slug        string             `json:"slug" validate:"required,min=3,max=200"`
	Author      *dto.InfoAuthor    `json:"author,omitempty" validate:"required,omitempty"`
	ImgURL      *string            `json:"imgUrl,omitempty" validate:"omitempty,uri,max=200"`
	Score       *float64           `json:"score,omitempty" validate:"omitempty,min=0,max=1"`
	Tags        *[]string          `json:"tags,omitempty" validate:"omitempty,dive,uppercase"`
	PublishedAt *time.Time         `json:"publishedAt,omitempty"`
}

func EmptyInfoPublicBlog() *PublicBlog {
	return &PublicBlog{}
}

func NewPublicBlog(blog *model.Blog, author *message.User) (*PublicBlog, error) {
	b, err := utils.MapTo[PublicBlog](blog)
	if err != nil {
		return nil, err
	}

	b.Author, err = utils.MapTo[dto.InfoAuthor](author)
	if err != nil {
		return nil, err
	}

	return b, err
}
