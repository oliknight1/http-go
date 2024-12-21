package main

func handle404(message *HTTPRequestMessage) string {
	responseMessage := &HTTPResponseMessage{
		HTTPVersion:   message.HTTPVersion,
		StatusCode:    404,
		StatusMessage: "Not Found",
	}
	return generateResponseMessage(responseMessage)
}
func handle405(message *HTTPRequestMessage) string {
	responseMessage := &HTTPResponseMessage{
		HTTPVersion:   message.HTTPVersion,
		StatusCode:    405,
		StatusMessage: "Method Not Allowed",
	}
	return generateResponseMessage(responseMessage)
}
