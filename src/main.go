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

	api, err := globals.Config.Data.FindAPI(r)
	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "uri fine fail")
		return
	}
	if err := api.Check(r.Header.Get("Content-Type"), r.Method); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "api content type of method fail")
		return
	}

	resp, err := api.Handle(r)
	if err != nil {
		log.Printf("api handler fail, %+v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "api content type of method fail")
		return
	}

	database.Scan()

	w.Header().Set("Content-Type", api.Response.Type)
	// w.Header().Set("Content-Type", "application/json")
	w.Write(api.Resp(r.Method, resp))
}
