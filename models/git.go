package models

type GitResponse struct {
	Action   string `json:"action"`
	Status   string `json:"status"`
	Branch   string `json:"branch,omitempty"`
	CommitID string `json:"commit_id,omitempty"`
	Message  string `json:"message,omitempty"`
	Error    string `json:"error,omitempty"`
}

type GitRequest struct {
	Project string `json:"project" validate:"required,alphanum" errmsg:"Project name must be alphanumeric"`
	RepoURL string `json:"repo_url" validate:"required,url" errmsg:"Invalid repository URL"`
	Pull    bool   `json:"pull" validate:"boolean" errmsg:"Pull must be a boolean value"`
}
