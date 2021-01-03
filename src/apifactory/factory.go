package apifactory

import "strings"

// API ...
type API interface{}

// NewAPI ...
func NewAPI(method string) API {
	switch strings.ToLower(method) {
	case "get":
		return &get{}
	case "post":
		return &post{}
	}
	return nil
}
