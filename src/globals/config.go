package globals

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"tomlserver/src/database"
	"tomlserver/src/resp"
)

// Config ...
var Config config

type config struct {
	Port string `toml:"port"`
	API  []api  `toml:"api"`
	DB   struct {
		Schema []string `toml:"schema"`
		Type   []string `toml:"type"`
	} `toml:"db"`
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

func (a *api) Database(param map[string]string) {
	if !a.Db {
		return
	}
	store := map[string]string{}
	for k, v := range param {
		for i := range Config.DB.Schema {
			if k == Config.DB.Schema[i] {
				store[k] = v
			}
		}
	}
	database.Insert(store)

}

func (a *api) GetParam(r *http.Request) map[string]string {
	param := map[string]string{}
	switch a.Contenttype {
	case "application/x-www-form-urlencoded":
		for i := range a.Parameter {
			v := r.FormValue(a.Parameter[i])
			if v != "" {
				param[a.Parameter[i]] = v
			}
		}
	case "application/json":
		body := make([]byte, 2048)
		len, err := r.Body.Read(body)
		if err != nil {
			if err != io.EOF {
				return param
			}
		}

		var jsonObj map[string]interface{}
		json.Unmarshal(body[:len], &jsonObj)

		for k, v := range jsonObj {
			switch reflect.TypeOf(v).String() {
			case "int":
				temp := strconv.Itoa(v.(int))
				param[k] = temp
			case "string":
				param[k] = v.(string)
			}

		}
	}
	return param
}

func (a *api) Resp() []byte {

	switch a.Response.Type {
	case "application/json":
		var resp = resp.Response{
			ErrorCode:    a.Response.ErrorCode,
			ErrorMessage: a.Response.ErrorMessage,
			Data:         a.Response.Data.Content[0],
		}
		return resp.ToBytes()
	case "application/x-www-form-urlencoded":
		data := url.Values{}
		data.Set("errorCode", strconv.Itoa(a.Response.ErrorCode))
		data.Set("errorMessage", a.Response.ErrorMessage)
		data.Set("data", a.Response.Data.Content[0])
		return []byte(data.Encode())
	case "text/plain":

	}

	return nil

}

func (a *api) Check(contenttype, method string) (err error) {
	if a.Contenttype != contenttype {
		return errors.New("api Content type wrong")
	}
	method = strings.ToLower(method)
	for i := range a.Method {
		if strings.ToLower(a.Method[i]) == method {
			return nil
		}
	}
	return errors.New("api method wrong")
}

func (c *config) Find(uri string) (find api, err error) {
	for i := range c.API {
		if c.API[i].Router == uri {
			find = c.API[i]
			return
		}
	}
	err = errors.New("find uri fail")
	return
}

func (c *config) Check() (err error) { //讀取設定檔時的確認
	if len(c.DB.Schema) != len(c.DB.Type) {
		return errors.New("db schema should equal type")
	}

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
