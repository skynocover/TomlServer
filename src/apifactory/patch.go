package apifactory

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"tomlserver/src/database"
	"tomlserver/src/globals"
	"tomlserver/src/resp"
)

type patch struct {
	//input
	contentType  string
	responseType string
	errorCode    int
	errorMessage string
	content      []string
	dataType     string
	//generate
	urlKey   string
	urlValue string
	body     map[string]string
}

func (p *patch) Response() []byte {
	switch p.dataType {
	case "text":
	case "hash":
		data := ""
		for _, v := range p.body {
			data = fmt.Sprintf("%s%s", data, v)
		}

		switch p.content[0] {
		case "sha256":
			sum := []byte(fmt.Sprintf("%s%s", data, p.content[1]))
			data = fmt.Sprintf("%x", sha256.Sum256(sum))
		case "md5":
			sum := []byte(fmt.Sprintf("%s%s", data, p.content[1]))
			data = fmt.Sprintf("%x", md5.Sum(sum))
		}
		p.content[0] = data
	}

	switch p.responseType {
	case "application/json":

		var resp = resp.Response{
			ErrorCode:    p.errorCode,
			ErrorMessage: p.errorMessage,
		}

		return resp.ToBytesWithObject(p.body)
	case "application/x-www-form-urlencoded":

		value := url.Values{}
		value.Set("errorCode", strconv.Itoa(p.errorCode))
		value.Set("errorMessage", p.errorMessage)
		value.Set("data", p.content[0])
		return []byte(value.Encode())
	case "text/plain":
		return []byte(p.content[0])
	}
	return nil
}

func (p *patch) GetParam(r *http.Request) {

	params := strings.Split(r.RequestURI, "/")
	log.Println(params)

	switch len(params) {
	case 3:
		p.urlKey = params[1]
		p.urlValue = params[2]
	default:
		return
	}

	p.body = map[string]string{}

	switch p.contentType {
	case "application/x-www-form-urlencoded":

		body := []byte{}
		length, err := r.Body.Read(body)
		if err != nil {
			log.Println(err)
			return
		}

		bodyquery, err := url.ParseQuery(string(body[:length]))
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(bodyquery)

		for k, v := range bodyquery {
			if len(v) > 0 {
				p.body[k] = v[0]
			}

		}

	case "application/json":
		body := make([]byte, 2048)
		len, err := r.Body.Read(body)
		if err != nil {
			if err != io.EOF {
				return
			}
		}

		var jsonObj map[string]interface{}
		json.Unmarshal(body[:len], &jsonObj)

		for k, v := range jsonObj {
			switch reflect.TypeOf(v).String() {
			case "int":
				temp := strconv.Itoa(v.(int))
				p.body[k] = temp
			case "string":
				p.body[k] = v.(string)
			}
		}
	}
	return
}

func (p *patch) Database() error {
	dbContent := map[string]string{}

	for _, schema := range globals.Config.Data.DB.Schema {
		for k1, v1 := range p.body {
			for i := range schema.Columns {
				if schema.Columns[i].Name == k1 {
					switch schema.Columns[i].Content[0] {
					case "text":
						dbContent[k1] = v1
					case "sha256":
						sum := sha256.Sum256([]byte(fmt.Sprintf("%s%s", v1, schema.Columns[i].Content[1])))
						dbContent[k1] = fmt.Sprintf("%x", sum)
					case "md5":
						data := []byte(fmt.Sprintf("%s%s", v1, schema.Columns[i].Content[1]))
						dbContent[k1] = fmt.Sprintf("%x", md5.Sum(data))
					}

				}
			}
		}
	}

	return database.Patch(p.urlKey, p.urlValue, dbContent)

}
