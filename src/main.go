package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
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
	log.Println(fmt.Sprintf("Run toml server on Port: %s", (globals.Config).Data.Port))
	log.Fatal(http.ListenAndServe(":"+globals.Config.Data.Port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {

	api, err := globals.Config.Data.Find(r.RequestURI)
	if err != nil {
		fmt.Fprintf(w, "uri fine fail")
		return
	}
	if err := api.Check(r.Header.Get("Content-Type"), r.Method); err != nil {
		fmt.Fprintf(w, "api content type of method fail")
		return
	}
	param := api.GetParam(r)
	log.Println(param)

	if api.Db {
		api.Database(param)
	}

	for k, v := range database.Content {
		fmt.Println(k, v)
	}

	w.Header().Set("Content-Type", api.Response.Type)
	// w.Header().Set("Content-Type", "application/json")
	w.Write(api.Resp())
}
