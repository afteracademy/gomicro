package dto

import (
	"github.com/afteracademy/goserve/v2/utility"
	"github.com/go-playground/validator/v10"
)

type CreateBlog struct {
	Title       string   `json:"title" validate:"required,min=3,max=500"`
	Description string   `json:"description" validate:"required,min=3,max=2000"`
	DraftText   string   `json:"draftText" validate:"required,max=50000"`
	Slug        string   `json:"slug" validate:"required,min=3,max=200"`
	ImgURL      string   `json:"imgUrl" validate:"required,uri,max=200"`
	Tags        []string `json:"tags" validate:"required,min=1,dive,uppercase"`
}

func EmptyCreateBlog() *CreateBlog {
	return &CreateBlog{}
}

func (d *CreateBlog) GetValue() *CreateBlog {
	return d
}

func (b *CreateBlog) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utility.FormatValidationErrors(errs), nil
}
