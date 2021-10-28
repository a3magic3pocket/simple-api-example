package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// TestRequest : 테스트 요청
func TestRequest(baseURL string, method string, data map[string]string, headers map[string]string, cookies []*http.Cookie) (statusCode int, body []byte, err error) {
	// Check method is valid
	method = strings.ToUpper(method)
	switch method {
	case "GET":
	case "POST":
	case "PATCH":
	case "DELETE":
	default:
		errorMessage := fmt.Sprintf("Method is not allowed, method: %s\n", method)
		return 0, []byte{}, errors.New(errorMessage)
	}

	// Get Content Type
	var contentType *string
	for headerName, headerValue := range headers {
		headerName = strings.ToLower(headerName)
		headerValue = strings.ToLower(headerValue)
		if headerName == "content-type" {
			contentTypeValue := headerValue
			contentType = &contentTypeValue
		}

	}

	var reqBody io.Reader
	switch {
	case contentType == nil:
		break
	case *contentType == "application/json":
		marshaled, err := json.Marshal(data)
		if err != nil {
			return 0, []byte{}, err
		}
		reqBody = bytes.NewBuffer(marshaled)
	case *contentType == "application/x-www-form-urlencoded":
		urlData := url.Values{}
		for key, value := range data {
			urlData.Set(key, value)
		}
		reqBody = strings.NewReader(urlData.Encode())
	default:
		errorMessage := fmt.Sprintf(`
			Content-Type is not allowed, 
			input Content-Type: %s, 
			allowed Content-type: application/json, application/x-www-form-urlencoded\n`,
			*contentType,
		)
		return 0, []byte{}, errors.New(errorMessage)
	}

	// Create request
	client := &http.Client{}
	req, err := http.NewRequest(method, baseURL, reqBody)
	if err != nil {
		return 0, []byte{}, err
	}
	if len(cookies) > 0 {
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
	}

	// Set Headers
	for headerName, headerValue := range headers {
		req.Header.Add(headerName, headerValue)
	}

	// Get response
	resp, err := client.Do(req)
	if err != nil {
		return resp.StatusCode, []byte{}, err
	}

	// Stringify response.Body
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, body, err
	}

	return resp.StatusCode, body, nil
}
