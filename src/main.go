package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"tomlserver/src/apifactory"
	"tomlserver/src/database"
	"tomlserver/src/globals"
)

var tomlData string

func main() {
	flag.StringVar(&tomlData, "config", "server.toml", "set the toml path")
	flag.Parse()

	globals.NewConfig(tomlData)

	// configInit()

	log.Println(fmt.Sprintf("%+v", globals.Config))

	http.HandleFunc("/", handler)
	log.Println(fmt.Sprintf("Run toml server on Port: %s", globals.Config.Data.Port))
	log.Fatal(http.ListenAndServe(":"+globals.Config.Data.Port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {

	api, err := globals.Config.Data.FindAPI(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "API fine fail, err: %v", err)
		return
	}

	newapi := apifactory.NewAPI(api.ContentType, api.Response.Type, api.Response.Data.Type, r.Method, api.Response.ErrorCode, api.Response.ErrorMessage, api.Response.Data.Content)

	newapi.GetParam(r)
	if api.Db {
		if err := newapi.Database(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, err.Error())
			return
		}
	}
	result := newapi.Response()

	database.Scan()

	w.Header().Set("Content-Type", api.Response.Type)
	w.Write(result)
}
