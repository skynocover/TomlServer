package handler

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"tomlserver/src/resp"
)

// Data ...
type Data struct {
	Port string `toml:"port"`
	API  []api  `toml:"api"`
	DB   struct {
		Schema []schema `toml:"schema"`
	} `toml:"db"`
}

type schema struct {
	Table   string    `toml:"table"`
	Key     string    `toml:"key"`
	Columns []columns `toml:"columns"`
}

type columns struct {
	Name    string   `toml:"name"`
	Content []string `toml:"content"`
}

type api struct {
	Router      string   `toml:"router"`
	Method      []string `toml:"method,omitempty"`
	Parameter   []string `toml:"parameter,omitempty"`
	Contenttype string   `toml:"contenttype,omitempty"`
	Db          bool     `toml:"db,omitempty"`
	Response    struct {
		Type         string `toml:"type"`
		ErrorCode    int    `toml:"errorCode"`
		ErrorMessage string `toml:"errorMessage"`
		Data         struct {
			Type    string   `toml:"type"`
			Content []string `toml:"content"`
		} `toml:"data"`
	} `toml:"response"`
}

// func (a *api) Handle(r *http.Request) []byte {

/*
	if !a.Db {
		_, _, v := a.getParam(r) //account , _ , password&code
		log.Println(v)
		res := ""
		for i := range a.Parameter {
			res = fmt.Sprintf("%s%s", res, v[a.Parameter[i]])
		}
		return res, nil
	}

	switch strings.ToLower(r.Method) {
	case "get":
		k, v := a.getURIParam(r)
		if k != "" && v != "" {
			res := database.Read(k, v)
			return res, nil
		} else if k != "" {
			res := database.ReadAll(k)
			return res, nil
		} else {
			return "", errors.New("fail uri param")
		}

	case "post":
		k, _, v := a.getParam(r) //account , _ , password&code
		dbContent := map[string]string{}

		for _, schema := range Config.Data.DB.Schema {
			for k1, v1 := range v {
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

		for _, schema := range Config.Data.DB.Schema {
			for k1, _ := range v {
				if schema.Key == k1 {
					database.Insert(k, v[schema.Key], dbContent)
				}
			}
		}
		return "success", nil

	case "put":
		k, r, v := a.getParam(r)
		dbContent := map[string]string{}
		for _, schema := range Config.Data.DB.Schema {
			for k1, v1 := range v {
				for i := range schema.Columns {
					if schema.Columns[i].Name == k1 {
						dbContent[k1] = v1
					}
				}
			}
		}

		err := database.Update(k, r, dbContent)

		return "success", err

	case "patch":
		k, r, v := a.getParam(r)

		dbContent := map[string]string{}
		for _, schema := range Config.Data.DB.Schema {
			for k1, v1 := range v {
				for i := range schema.Columns {
					if schema.Columns[i].Name == k1 {
						dbContent[k1] = v1
					}
				}
			}
		}

		err := database.Patch(k, r, dbContent)

		return "success", err

	case "delete":
		k, v := a.getURIParam(r)
		database.Delete(k, v)
		return "", nil
	default:
		return "", fmt.Errorf("no math method")
	}
*/
// }

func (a *api) getParam(r *http.Request) (table, row string, value map[string]string) {
	value = map[string]string{}

	params := strings.Split(r.RequestURI, "/")
	log.Println(params)

	switch len(params) {
	case 2:
		table = params[1]
	case 3:
		table = params[1]
		row = params[2]
	default:
		return
	}

	switch a.Contenttype {
	case "application/x-www-form-urlencoded":
		for i := range a.Parameter {
			v := r.FormValue(a.Parameter[i])
			if v != "" {
				value[a.Parameter[i]] = v
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
				value[k] = temp
			case "string":
				value[k] = v.(string)
			}
		}
	}
	return

}

func (a *api) getURIParam(r *http.Request) (key, value string) {
	params := strings.Split(r.RequestURI, "/")
	log.Println(params)

	switch len(params) {
	case 0, 1:
	case 2:
		key = params[1]
	case 3:
		key = params[1]
		value = params[2]
	}
	return
}

func (a *api) Resp(method string, response interface{}) []byte {

	data := ""
	switch a.Response.Data.Type {
	case "text":
		data = a.Response.Data.Content[0]
	case "db":
	case "hash":
		switch a.Response.Data.Content[0] {
		case "sha256":
			result := response.(string)
			sum := sha256.Sum256([]byte(fmt.Sprintf("%s%s", result, a.Response.Data.Content[1])))
			data = fmt.Sprintf("%x", sum)
		case "md5":
			result := response.(string)
			md5data := []byte(fmt.Sprintf("%s%s", result, a.Response.Data.Content[1]))
			data = fmt.Sprintf("%x", md5.Sum(md5data))
		}
		// data = response
	}

	switch a.Response.Type {
	case "application/json":
		var resp = resp.Response{
			ErrorCode:    a.Response.ErrorCode,
			ErrorMessage: a.Response.ErrorMessage,
			Data:         data,
		}
		return resp.ToBytes()
	case "application/x-www-form-urlencoded":

		value := url.Values{}
		value.Set("errorCode", strconv.Itoa(a.Response.ErrorCode))
		value.Set("errorMessage", a.Response.ErrorMessage)

		if a.Response.Data.Type == "db" && strings.ToLower(method) == "get" {
			resp := response.(map[string]string)
			for k, v := range resp {
				value.Set(k, v)
			}
		} else {
			value.Set("data", data)
		}

		return []byte(value.Encode())
	case "text/plain":
		return []byte(response.(string))
	}

	return nil
}

func (a *api) Check(contenttype, method string) (err error) {
	method = strings.ToLower(method)
	for i := range a.Method {
		if strings.ToLower(a.Method[i]) == method {
			if method == "post" || method == "put" || method == "patch" {
				if a.Contenttype != contenttype {
					return errors.New("api Content type wrong")
				}
			}
			return nil
		}
	}
	return errors.New("api method wrong")
}

func (c *Data) FindAPI(r *http.Request) (find api, err error) {
	uri := ""
	realRouters := strings.Split(r.RequestURI, "/")
	switch len(realRouters) {
	case 0, 1:
		err = errors.New("find api fail")
		return
	default:
		uri = realRouters[1]
	}

	switch strings.ToLower(r.Method) {
	case "get", "delete", "put", "patch":
		for i := range c.API {
			if strings.Contains(c.API[i].Router, uri) {
				find = c.API[i]
				return
			}
		}
	case "post":
		for i := range c.API {
			if strings.Contains(c.API[i].Router, uri) {
				find = c.API[i]
				return
			}
		}
	}

	err = errors.New("find uri fail")
	return
}

func (c *Data) Check() (err error) { //讀取設定檔時的確認

	for i := range c.API {
		switch c.API[i].Contenttype {
		case "application/x-www-form-urlencoded":
		case "application/json":
		default:
			return errors.New("content type shoule be form or json")
		}
	}

	return nil
}
