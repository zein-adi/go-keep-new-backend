package keep_request

type KantongHistoryInsertUpdate struct {
	Id        string `json:"id,omitempty"`
	KantongId string `json:"kantongId,omitempty" validate:"required,number"`
	Jumlah    int    `json:"jumlah,omitempty" validate:"required,number"`
	Uraian    string `json:"uraian,omitempty" validate:""`
}
