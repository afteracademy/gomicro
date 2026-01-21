package dto

func EmptyTag() *Tag {
	return &Tag{}
}

type Tag struct {
	Tag string `uri:"tag" validate:"required,uppercase"`
}
