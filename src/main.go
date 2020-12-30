package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"tomlserver/src/globals"

	"github.com/BurntSushi/toml"
)

func main() {
	var tomlData string
	flag.StringVar(&tomlData, "config", "server.toml", "set the toml path")
	flag.Parse()

	if _, err := toml.DecodeFile(tomlData, &globals.Config); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handler)
	log.Println(fmt.Sprintf("Run toml server on Port: %s", globals.Config.Port))
	log.Fatal(http.ListenAndServe(":"+globals.Config.Port, nil))

}

func handler(w http.ResponseWriter, r *http.Request) {

	api, err := globals.Config.Find(r.RequestURI)
	if err != nil {
		fmt.Fprintf(w, "uri fine fail")
		return
	}
	if err := api.Check(r.Header.Get("Content-Type"), r.Method); err != nil {
		fmt.Fprintf(w, "api content type of method fail")
		return
	}
	api.GetParam(r)

	w.Header().Set("Content-Type", "application/json")
	// w.Write(api.Resp())
	fmt.Fprintf(w, api.Resp())
}
