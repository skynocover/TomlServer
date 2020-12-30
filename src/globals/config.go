package globals

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

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
		ErrorCode    int    `toml:"errorCode"`
		ErrorMessage string `toml:"errorMessage"`
		Data         struct {
			Type    string   `toml:"type"`
			Content []string `toml:"content"`
		} `toml:"data"`
	} `toml:"response"`
}

func (a *api) GetParam(r *http.Request) {
	switch a.Contenttype {
	case "application/x-www-form-urlencoded":
		for i := range a.Parameter {
			v := r.FormValue(a.Parameter[i])
			log.Println(v)
			// a.Parameter[i]
		}
	}

}

func (a *api) Resp() string {
	var resp = resp.Response{
		ErrorCode:    a.Response.ErrorCode,
		ErrorMessage: a.Response.ErrorMessage,
		Data:         a.Response.Data.Content[0],
	}

	return fmt.Sprintf("errorCode=%d&errorMessage=%s&data=%s", resp.ErrorCode, resp.ErrorMessage, resp.Data)

	// return resp.ToBytes()

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

func (c *config) Check() (err error) {
	if len(c.DB.Schema) != len(c.DB.Type) {
		return errors.New("db schema should equal type")
	}

	for i := range c.API {
		switch c.API[i].Contenttype {
		case "x-www-form-urlencoded":
		case "json":
		default:
			return errors.New("content type shoule be form or json")
		}
	}

	return nil
}
