package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func handleUsers(message *HTTPRequestMessage) string {
	parts := strings.SplitN(message.Path, "?", 2)
	queries := ""
	if len(parts) > 1 {
		queries = parts[1]
	}

	switch message.Method {
	case "GET":
		if len(queries) == 0 {
			return getAllUsers(message)
		}
		return getUserById(message, queries)
	case "POST":
		return postUser(message)
	default:
		return handle405(message)

	}
}

func getAllUsers(message *HTTPRequestMessage) string {
	var body []byte
	allUsers, err := json.Marshal(&Data)
	body = allUsers

	if err != nil {
		panic(err)
	}
	responseMessage := &HTTPResponseMessage{
		HTTPVersion:     message.HTTPVersion,
		StatusCode:      200,
		StatusMessage:   "OK",
		Body:            body,
		ResponseHeaders: map[string]string{},
	}
	return generateResponseMessage(responseMessage)
}

func getUserById(message *HTTPRequestMessage, queries string) string {
	var body []byte
	splitQueries := strings.Split(queries, "&")

	queryMap := make(map[string]string)

	for _, query := range splitQueries {
		separatedQuery := strings.Split(query, "=")
		queryMap[separatedQuery[0]] = separatedQuery[1]
	}

	//NOTE: currently only support id query params
	idToSearch, err := strconv.Atoi(queryMap["id"])

	if err != nil {
		fmt.Println("id query param was not a number")
		return ""
	}

	var foundUser User

	for _, user := range Data.Users {
		if user.ID == idToSearch {
			foundUser = user
			break
		}
	}
	user, err := json.Marshal(&foundUser)
	if err != nil {
		panic(err)
	}
	body = user

	responseMessage := &HTTPResponseMessage{
		HTTPVersion:     message.HTTPVersion,
		StatusCode:      200,
		StatusMessage:   "OK",
		Body:            body,
		ResponseHeaders: map[string]string{},
	}
	return generateResponseMessage(responseMessage)
}

func postUser(message *HTTPRequestMessage) string {
	representationHeaders := map[string]string{}

	if val, exists := message.Headers["Content-Length"]; exists {
		representationHeaders["Content-Length"] = val
	} else {
		handle404(message)
	}
	if val, exists := message.Headers["Content-Type"]; exists {
		representationHeaders["Content-Type"] = val
	} else {
		handle404(message)
	}

	responseMessage := &HTTPResponseMessage{
		HTTPVersion:           message.HTTPVersion,
		StatusCode:            201,
		StatusMessage:         "OK",
		RepresentationHeaders: representationHeaders,
	}
	return generateResponseMessage(responseMessage)
}
