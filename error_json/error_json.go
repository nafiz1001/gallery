package error_json

type ErrorJson struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e *ErrorJson) Error() string {
	return e.Message
}
