package types

type Response struct {
	Status  string `json:"status"`
	Result  any    `json:"result"`
	Meta    any    `json:"meta,omitempty"`
	Message string `json:"message"`
}
