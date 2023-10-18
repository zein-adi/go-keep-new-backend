package keep_request

type KantongHistoryInsertUpdate struct {
	Id        string `json:"id,omitempty"`
	KantongId string `json:"kantongId,omitempty" validate:"required,number"`
	Jumlah    int    `json:"jumlah,omitempty" validate:"number,min=0"`
	Uraian    string `json:"uraian,omitempty" validate:""`
}
