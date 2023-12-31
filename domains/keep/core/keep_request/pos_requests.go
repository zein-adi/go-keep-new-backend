package keep_request

type PosInputUpdate struct {
	Id       string `json:"id,omitempty" validate:""`
	Nama     string `json:"nama,omitempty" validate:"required"`
	Urutan   int    `json:"urutan,omitempty" validate:"required,number,min=1"`
	ParentId string `json:"parentId,omitempty" validate:""`
	IsShow   bool   `json:"isShow,omitempty" validate:"boolean"`
}
type PosUpdateUrutanItem struct {
	Id       string `json:"id,omitempty" validate:"required,number"`
	Urutan   int    `json:"urutan,omitempty" validate:"required,number,min=1"`
	ParentId string `json:"parentId,omitempty" validate:"omitempty,number"`
}
type PosUpdateVisibilityItem struct {
	Id     string `json:"id,omitempty" validate:"required,number"`
	IsShow bool   `json:"isShow,omitempty" validate:"boolean"`
}
