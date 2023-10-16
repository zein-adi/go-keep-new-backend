package auth_responses

type ProfileResponse struct {
	Username    string   `json:"username"`
	Nama        string   `json:"nama"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}
