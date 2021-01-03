package globals

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"tomlserver/src/database"
	"tomlserver/src/resp"

	"github.com/BurntSushi/toml"
)

// Config ...
var Config *tconfig

type tconfig struct {
	Filename       string
	LastModifyTime int64
	Lock           *sync.RWMutex
	Data           *configData
}

// NewConfig ...
func NewConfig(filename string) {
	Config = &tconfig{
		Filename:       filename,
		Lock:           &sync.RWMutex{},
		Data:           &configData{},
		LastModifyTime: 0,
	}

	Config.parse(true)

	go func() {
		for {
			time.Sleep(time.Second * 5)
			Config.parse(false)
		}
	}()
}

func (c *tconfig) parse(init bool) {
	fileInfo, _ := os.Stat(c.Filename)
	currModifyTime := fileInfo.ModTime().Unix()
	if currModifyTime > c.LastModifyTime {
		c.LastModifyTime = currModifyTime
		Config.Lock.Lock()
		if _, err := toml.DecodeFile(c.Filename, c.Data); err != nil {
			log.Println("open config err: ", err)
			if err != nil && init {
				os.Exit(1)
			}
		}
		Config.Lock.Unlock()
		log.Printf("Config = %+v\n", c.Data)
	}
}

// Tconfig ...
type configData struct {
	Port string `toml:"port"`
	API  []api  `toml:"api"`
	DB   struct {
		Schema []schema `toml:"schema"`
	} `toml:"db"`
}

type schema struct {
	Key  string `toml:"key"`
	Type string `toml:"type"`
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
	store := map[string]string{}
	for k, v := range param {
		for _, schema := range Config.Data.DB.Schema {
			if k == schema.Key {
				store[k] = v
			}
		}
	}
	database.Insert(store)
}

func (a *api) Method(r *http.Request) {
	switch method {
	case "get":

	case "post":
	case "put":
	case "patch":
	case "delete":

	}
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

	data := ""
	switch a.Response.Data.Type {
	case "text":
		data = a.Response.Data.Content[0]
	case "db":
		for _, content := range a.Response.Data.Content {
			for _, db := range database.Content {
				v, ok := db[content]
				if ok {
					data = fmt.Sprintf("%s: %s, ", content, v)
				}
			}
		}
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
		value.Set("data", data)
		return []byte(value.Encode())
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

func (c *configData) Find(uri string) (find api, err error) {
	for i := range c.API {
		if c.API[i].Router == uri {
			find = c.API[i]
			return
		}
	}
	err = errors.New("find uri fail")
	return
}

func (c *configData) Check() (err error) { //讀取設定檔時的確認

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
