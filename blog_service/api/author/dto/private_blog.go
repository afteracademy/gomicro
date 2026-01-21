package dto

import (
	"time"

	"github.com/afteracademy/gomicro/blog-service/api/auth/message"
	"github.com/afteracademy/gomicro/blog-service/api/blog/model"
	"github.com/afteracademy/goserve/v2/utility"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PrivateBlog struct {
	ID          primitive.ObjectID `json:"_id" binding:"required" validate:"required"`
	Title       string             `json:"title" validate:"required,min=3,max=500"`
	Description string             `json:"description" validate:"required,min=3,max=2000"`
	Text        *string            `json:"text,omitempty" validate:"omitempty,max=50000"`
	DraftText   string             `json:"draftText" validate:"required"`
	Slug        string             `json:"slug" validate:"required,min=3,max=200"`
	Author      *InfoAuthor        `json:"author,omitempty" validate:"required,omitempty"`
	ImgURL      *string            `json:"imgUrl,omitempty" validate:"omitempty,uri,max=200"`
	Score       *float64           `json:"score,omitempty" validate:"omitempty,min=0,max=1"`
	Tags        *[]string          `json:"tags,omitempty" validate:"omitempty,dive,uppercase"`
	Submitted   bool               `json:"submitted" validate:"-"`
	Drafted     bool               `json:"drafted" validate:"-"`
	Published   bool               `json:"published" validate:"-"`
	PublishedAt *time.Time         `json:"publishedAt,omitempty"`
	CreatedAt   time.Time          `json:"createdAt" validate:"required"`
	UpdatedAt   time.Time          `json:"updatedAt" validate:"required"`
}

func EmptyInfoPrivateBlog() *PrivateBlog {
	return &PrivateBlog{}
}

func NewPrivateBlog(blog *model.Blog, author *message.User) (*PrivateBlog, error) {
	b, err := utility.MapTo[PrivateBlog](blog)
	if err != nil {
		return nil, err
	}

	b.Author, err = utility.MapTo[InfoAuthor](author)
	if err != nil {
		return nil, err
	}

	return b, err
}

func (d *PrivateBlog) GetValue() *PrivateBlog {
	return d
}

func (b *PrivateBlog) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utility.FormatValidationErrors(errs), nil
}
