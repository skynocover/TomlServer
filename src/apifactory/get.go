package apifactory

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"tomlserver/src/database"
	"tomlserver/src/resp"
)

type get struct {
	//input
	responseType string
	errorCode    int
	errorMessage string
	content      []string
	dataType     string
	//generate
	urlKey   string
	urlValue string
	// response
	singleMap map[string]string
	allMap    map[string]map[string]string
}

func (g *get) GetParam(r *http.Request) {

	params := strings.Split(r.RequestURI, "/")
	log.Println(params)

	switch len(params) {
	case 2:
		g.urlKey = params[1]
	case 3:
		g.urlKey = params[1]
		g.urlValue = params[2]
	default:
		return
	}

	return
}

func (g *get) Database() error {
	if g.urlKey != "" && g.urlValue != "" {
		g.singleMap = database.Read(g.urlKey, g.urlValue)
	} else if g.urlKey != "" {
		g.allMap = database.ReadAll(g.urlKey)
	}
	return nil
}

func (g *get) Response() []byte {
	switch g.dataType {
	case "text":
	case "hash":
		data := ""
		switch {
		case g.singleMap != nil:
			for _, v := range g.singleMap {
				data = fmt.Sprintf("%s%s", data, v)
			}
		case g.allMap != nil:
			for _, v1 := range g.allMap {
				for _, v2 := range v1 {
					data = fmt.Sprintf("%s%s", data, v2)
				}
			}
		}

		switch g.content[0] {
		case "sha256":
			sum := []byte(fmt.Sprintf("%s%s", data, g.content[1]))
			data = fmt.Sprintf("%x", sha256.Sum256(sum))
		case "md5":
			sum := []byte(fmt.Sprintf("%s%s", data, g.content[1]))
			data = fmt.Sprintf("%x", md5.Sum(sum))
		}
		g.content[0] = data
	}

	switch g.responseType {
	case "application/json":

		var resp = resp.Response{
			ErrorCode:    g.errorCode,
			ErrorMessage: g.errorMessage,
		}

		switch g.dataType {
		case "db":
			return resp.ToBytesWithObject(g.singleMap)
		case "text":
			return []byte(g.content[0])
		case "hash":
			return []byte(g.content[0])
		}

		return resp.ToBytes()
	case "application/x-www-form-urlencoded":

		value := url.Values{}
		value.Set("errorCode", strconv.Itoa(g.errorCode))
		value.Set("errorMessage", g.errorMessage)

		switch g.dataType {
		case "db":
			for k, v := range g.singleMap {
				value.Set(k, v)
			}
		default:
			value.Set("data", g.content[0])
		}
		return []byte(value.Encode())
	case "text/plain":
		switch g.dataType {
		case "db":
			jsondata, _ := json.Marshal(g.singleMap)
			return jsondata
		default:
			return []byte(g.content[0])
		}
	}

	return nil
}
