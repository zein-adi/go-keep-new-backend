package keep_request

func NewPosGetRequest() *PosGetRequest {
	return &PosGetRequest{
		IsLeafOnly: false,
	}
}

type PosGetRequest struct {
	IsLeafOnly bool
}

type PosInputUpdateRequest struct {
	Id       string `json:"id,omitempty" validate:""`
	Nama     string `json:"nama,omitempty" validate:"required"`
	Urutan   int    `json:"urutan,omitempty" validate:"required,number,min=1"`
	ParentId string `json:"parentId,omitempty" validate:""`
	IsShow   bool   `json:"isShow,omitempty" validate:"boolean"`
}
