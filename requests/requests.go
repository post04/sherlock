package requests

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Jeffail/gabs"
)

// NewRequest builds a new request based off input arguments, returns http.Request (nil if err isn't nil)
func NewRequest(method, url, username, body string, headers [][]string) *http.Request {
	if body != "" {
		body = strings.ReplaceAll(body, "{}", username)
		req, err := http.NewRequest(method, strings.ReplaceAll(url, "{}", username), strings.NewReader(body))
		if err != nil {
			return nil
		}
		for _, header := range headers {
			if len(header) == 2 {
				req.Header.Set(header[0], header[1])
			}
		}
		return req
	}
	req, err := http.NewRequest(method, strings.ReplaceAll(url, "{}", username), nil)
	if err != nil {
		return nil
	}
	for _, header := range headers {
		if len(header) == 2 {
			req.Header.Set(header[0], header[1])
		}
	}

	return req
}

// MakeRequest takes in a http.Request and an http.Client then returns a http.Response (nil if err isn't nil)
func MakeRequest(req *http.Request, client *http.Client) *http.Response {
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	return resp
}

// CheckValidity takes in what type of thing we're checking for and does operations based off of that
func CheckValidity(resp *http.Response, arguments ...interface{}) bool {
	defer resp.Body.Close()
	switch arguments[0].(string) {
	case "status":
		return resp.StatusCode == int(arguments[1].(float64))
	case "body":
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return false
		}
		if strings.HasPrefix(arguments[1].(string), "!") {
			return !strings.Contains(string(b), strings.Replace(arguments[1].(string), "!", "", 1))
		}
		return strings.Contains(string(b), arguments[1].(string))
	case "struct":
		// W.I.P for now just try to use the body matching with structs

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return false
		}
		jsonParsed, err := gabs.ParseJSON(b)
		if err != nil {
			return false
		}
		toCheck := jsonParsed.Path(strings.Replace(arguments[1].(string), "!", "", 1)).String()
		if strings.HasPrefix(arguments[1].(string), "!") {
			return toCheck != strings.Replace(fmt.Sprintf("%v", arguments[2]), "!", "", 1)
		}
		return toCheck == fmt.Sprintf("%v", arguments[2])
	default:
		return false
	}
}
