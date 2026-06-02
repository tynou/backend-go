package response

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func OK(msg string) Response {
	return Response{
		Status:  "OK",
		Message: msg,
	}
}

func Error(msg string) Response {
	return Response{
		Status: "Error",
		Error:  msg,
	}
}
