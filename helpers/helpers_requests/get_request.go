package helpers_requests

func NewGet() *Get {
	return &Get{
		Skip:   0,
		Take:   10,
		Search: "",
	}
}

type Get struct {
	Skip   int    `json:"skip,omitempty" validate:""`
	Take   int    `json:"take,omitempty" validate:"min=10"`
	Search string `json:"search,omitempty" validate:""`
}
