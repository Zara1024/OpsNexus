package model

type KubectlCommandRequest struct {
	Command        string `json:"command" binding:"required"`
	Namespace      string `json:"namespace,omitempty"`
	TimeoutSeconds int    `json:"timeoutSeconds,omitempty"`
}

type KubectlCommandResponse struct {
	Success      bool   `json:"success"`
	Command      string `json:"command"`
	Namespace    string `json:"namespace,omitempty"`
	Stdout       string `json:"stdout,omitempty"`
	Stderr       string `json:"stderr,omitempty"`
	ExitCode     int    `json:"exitCode"`
	DurationMs   int64  `json:"durationMs"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}
