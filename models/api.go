package models
type CommandRequest struct {
	Url   string `json:"url"`
	Mode  string `json:"mode"`
	Query string `json:"query"`
}
type OpenProjectRequest struct {
	Url     string `json:"url"`
	Project string `json:"project" validate:"required,alphanumeric"`
}
type Response struct {
	Data  any `json:"data"`
	Error any `json:"error"`
}
