package response

type ErrorResponse struct {
	Message  string `json:"message"`
	Instance string `json:"instance"`
	Detail   string `json:"detail"`
}
