package message

type Text struct {
	Value string `json:"value" validate:"required"`
}

func NewText(value string) *Text {
	return &Text{
		Value: value,
	}
}
