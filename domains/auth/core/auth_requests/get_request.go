package auth_requests

func NewGet() Get {
	return Get{
		Skip:   0,
		Take:   10,
		Search: "",
	}
}

type Get struct {
	Skip   int
	Take   int
	Search string
}
