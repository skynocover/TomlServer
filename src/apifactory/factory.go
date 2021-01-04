package apifactory

import (
	"net/http"
	"strings"
)

// API ...
type API interface {
	GetParam(*http.Request)
	Database() error
	Response() []byte
}

// NewAPI ...
func NewAPI(contentType, responseType, dataType, method string, errorCode int, errorMessage string, content []string) API {
	switch strings.ToLower(method) {
	case "get":
		return &get{
			responseType: responseType,
			dataType:     dataType,
			errorCode:    errorCode,
			errorMessage: errorMessage,
			content:      content,
		}
	case "post":
		return &post{
			contentType:  contentType,
			responseType: responseType,
			dataType:     dataType,
			errorCode:    errorCode,
			errorMessage: errorMessage,
			content:      content,
		}
	case "patch":
		return &patch{
			contentType:  contentType,
			responseType: responseType,
			dataType:     dataType,
			errorCode:    errorCode,
			errorMessage: errorMessage,
			content:      content,
		}
	case "put":
		return &put{
			contentType:  contentType,
			responseType: responseType,
			dataType:     dataType,
			errorCode:    errorCode,
			errorMessage: errorMessage,
			content:      content,
		}
	case "delete":
		return &delete{
			responseType: responseType,
			dataType:     dataType,
			errorCode:    errorCode,
			errorMessage: errorMessage,
			content:      content,
		}
	}
	return nil
}
