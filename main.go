package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type HTTPRequestMessage struct {
	Method        string
	Path          string
	HTTPVersion   string
	Headers       map[string]string
	ContentLength int
	Body          []byte
}

type HTTPResponseMessage struct {
	HTTPVersion           string
	StatusCode            int
	StatusMessage         string
	Body                  []byte
	ResponseHeaders       map[string]string
	RepresentationHeaders map[string]string
}

func main() {
	lister, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("Listening to port failed: %v ", err)
	}
	defer lister.Close()

	for {
		conn, err := lister.Accept()

		if err != nil {
			fmt.Printf("Accepting request failed: %v ", err)
		}

		go handleConn(conn)

	}

}

func handleConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// NOTE: Read string removes the read string from it's internal buffer
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	if len(requestLine) == 0 {
		fmt.Println("ERROR: request is empty")
		return
	}

	requestMessage := HTTPRequestMessage{}

	splitReq := strings.Split(requestLine, " ")

	requestMessage.Method = strings.TrimSpace(splitReq[0])
	requestMessage.Path = strings.TrimSpace(splitReq[1])
	requestMessage.HTTPVersion = strings.TrimSpace(splitReq[2])

	headers := make(map[string]string)
	contentLength := 0

	for {
		line, err := reader.ReadString('\n')

		line = strings.TrimSpace(line)

		if err != nil || line == "" {
			break
		}
		headerLine := strings.SplitN(line, ":", 2)
		if len(headerLine) != 2 {
			fmt.Println("Error splitting headers")
			break
		}

		key := strings.TrimSpace(headerLine[0])
		value := strings.TrimSpace(headerLine[1])

		headers[key] = value

		if key == "Content-Length" {
			contentLen, err := strconv.Atoi(value)
			if err != nil {
				fmt.Println("Error converting Content-Length to int")
				return
			}
			contentLength = contentLen
		}

	}

	requestMessage.Headers = headers

	// Read body if exists
	var body []byte
	if contentLength > 0 {
		body = make([]byte, contentLength)
		_, err := io.ReadFull(reader, body)
		if err != nil {
			fmt.Println(err)
			return
		}
		requestMessage.Body = body
	}

	responseMsg := handleRouting(&requestMessage)
	_, err = conn.Write([]byte(responseMsg))
	if err != nil {
		fmt.Println("Error writing response:", err)
		return
	}

}

func handleHome(message *HTTPRequestMessage) string {
	responseMessage := &HTTPResponseMessage{
		HTTPVersion:     message.HTTPVersion,
		StatusCode:      200,
		StatusMessage:   "OK",
		ResponseHeaders: map[string]string{},
	}
	return generateResponseMessage(responseMessage)
}

func handleRouting(message *HTTPRequestMessage) string {
	parts := strings.SplitN(message.Path, "?", 2)
	path := parts[0]

	switch path {
	case "/":
		return handleHome(message)
	case "/users":
		return handleUsers(message)
	default:
		return handle404(message)
	}
}

func generateTimestamp() string {
	return time.Now().UTC().Format(time.RFC1123)
}

func generateResponseMessage(responseMessage *HTTPResponseMessage) string {

	responseMessage.ResponseHeaders["Server"] = "MyServer"
	responseMessage.ResponseHeaders["Date"] = generateTimestamp()

	var response strings.Builder
	statusLine := fmt.Sprintf(
		"%s %d %s\r\n",
		responseMessage.HTTPVersion,
		responseMessage.StatusCode,
		responseMessage.StatusMessage,
	)
	response.WriteString(statusLine)

	var headers strings.Builder

	for key, value := range responseMessage.ResponseHeaders {
		headers.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	for key, value := range responseMessage.RepresentationHeaders {
		headers.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	response.WriteString(headers.String())

	response.WriteString("\r\n")
	response.Write(responseMessage.Body)

	return response.String()
}
