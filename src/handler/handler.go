package handler

import (
	"errors"
	"net/http"
	"strings"
)

// Data ...
type Data struct {
	Port string `toml:"port"`
	API  []API  `toml:"api"`
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

// API config
type API struct {
	Router      string   `toml:"router"`
	Method      []string `toml:"method,omitempty"`
	Parameter   []string `toml:"parameter,omitempty"`
	ContentType string   `toml:"contenttype,omitempty"`
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

// FindAPI ...
func (d *Data) FindAPI(r *http.Request) (find API, err error) {
	uri := ""
	realRouters := strings.Split(r.RequestURI, "/")
	switch len(realRouters) {
	case 0, 1:
		err = errors.New("find api fail")
		return
	default:
		uri = realRouters[1]
	}

	method := strings.ToLower(r.Method)

	switch method {
	case "get", "delete":
		for i := range d.API {
			if strings.Contains(d.API[i].Router, uri) {
				for i := range d.API[i].Method {
					if strings.ToLower(d.API[i].Method[i]) == method {
						return d.API[i], nil
					}

					err = errors.New("API method wrong")
					return
				}
			}
		}

	case "post", "put", "patch":
		for i := range d.API {
			if strings.Contains(d.API[i].Router, uri) {
				for i := range d.API[i].Method {
					if strings.ToLower(d.API[i].Method[i]) == method {
						if d.API[i].ContentType != r.Header.Get("Content-Type") {
							err = errors.New("API Content-Type wrong")
							return
						}
						return d.API[i], nil
					}
				}
				err = errors.New("API method wrong")
				return
			}
		}
	}

	err = errors.New("find uri fail")
	return
}

// Check ...
func (d *Data) Check() (err error) { //讀取設定檔時的確認

	for i := range d.API {
		switch d.API[i].ContentType {
		case "application/x-www-form-urlencoded":
		case "application/json":
		default:
			return errors.New("content type shoule be form or json")
		}
	}

	return nil
}
