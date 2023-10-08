package auth_requests

func NewGetRequest() GetRequest {
	return GetRequest{
		Skip:   0,
		Take:   10,
		Search: "",
	}
}

type GetRequest struct {
	Skip   int
	Take   int
	Search string
}
