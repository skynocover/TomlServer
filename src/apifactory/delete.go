package apifactory

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"tomlserver/src/database"
	"tomlserver/src/resp"
)

type delete struct {
	//input
	responseType string
	errorCode    int
	errorMessage string
	content      []string
	dataType     string
	//generate
	urlKey   string
	urlValue string
}

func (d *delete) GetParam(r *http.Request) {

	params := strings.Split(r.RequestURI, "/")
	log.Println(params)

	switch len(params) {
	case 2:
		d.urlKey = params[1]
	case 3:
		d.urlKey = params[1]
		d.urlValue = params[2]
	default:
		return
	}

	return
}

func (d *delete) Database() error {
	database.Delete(d.urlKey, d.urlValue)
	return nil
}

func (d *delete) Response() []byte {
	switch d.dataType {
	case "text":
	case "hash":
		data := ""

		d.content[0] = data
	}

	switch d.responseType {
	case "application/json":

		var resp = resp.Response{
			ErrorCode:    d.errorCode,
			ErrorMessage: d.errorMessage,
			Data:         d.content[0],
		}

		return resp.ToBytes()
	case "application/x-www-form-urlencoded":

		value := url.Values{}
		value.Set("errorCode", strconv.Itoa(d.errorCode))
		value.Set("errorMessage", d.errorMessage)
		value.Set("data", d.content[0])

		return []byte(value.Encode())
	case "text/plain":
		switch d.dataType {
		case "db":
			jsondata, _ := json.Marshal(d.content[0])
			return jsondata
		default:
			return []byte(d.content[0])
		}
	}

	return nil
}
