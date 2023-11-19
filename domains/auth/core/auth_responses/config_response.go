package auth_responses

var Version = "1.0.3"

type ConfigResponse struct {
	Version string `json:"version"`
}
