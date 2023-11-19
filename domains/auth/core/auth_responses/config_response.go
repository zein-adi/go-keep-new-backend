package auth_responses

var Version = "1.0.2"

type ConfigResponse struct {
	Version string `json:"version"`
}
