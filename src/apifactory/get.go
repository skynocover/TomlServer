package apifactory

import (
	"log"
	"net/http"
	"strings"
	"tomlserver/src/database"
)

type get struct {
	//input
	responseType string
	errorCode    int
	errorMessage string
	content      []string
	//generate
	urlKey    string
	urlValue  string
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

func (g *get) Database() {
	if g.urlKey != "" && g.urlValue != "" {
		g.singleMap = database.Read(g.urlKey, g.urlValue)
	} else if g.urlKey != "" {
		g.allMap = database.ReadAll(g.urlKey)
	}
}

func (g *get) Response() []byte {
	return nil
}
