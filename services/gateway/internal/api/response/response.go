package response

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func OK(msg string) Response {
	return Response{
		Status:  "ok",
		Message: msg,
	}
}

func Error(msg string) Response {
	return Response{
		Status: "error",
		Error:  msg,
	}
}
