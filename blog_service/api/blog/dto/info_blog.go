package dto

import (
	"github.com/afteracademy/gomicro/blog-service/api/blog/model"
	"github.com/afteracademy/goserve/v2/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InfoBlog struct {
	ID          primitive.ObjectID `json:"_id" binding:"required" validate:"required"`
	Title       string             `json:"title" validate:"required,min=3,max=500"`
	Description string             `json:"description" validate:"required,min=3,max=2000"`
	Slug        string             `json:"slug" validate:"required,min=3,max=200"`
	ImgURL      *string            `json:"imgUrl,omitempty" validate:"omitempty,uri,max=200"`
	Score       float64            `json:"score," validate:"required,min=0,max=1"`
	Tags        []string           `json:"tags" validate:"required,dive,uppercase"`
}

func NewInfoBlog(blog *model.Blog) (*InfoBlog, error) {
	return utility.MapTo[InfoBlog](blog)
}
