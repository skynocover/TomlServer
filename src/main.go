package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tomlserver/src/database"
	"tomlserver/src/globals"

	"github.com/BurntSushi/toml"
)

var tomlData string

func configInit() {
	loadConfig(true)
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGUSR2)
	go func() {
		for {
			<-s
			loadConfig(false)
			log.Println("Reloaded")
		}
	}()
}

func loadConfig(fail bool) {
	temp := new(globals.Tconfig)
	if _, err := toml.DecodeFile(tomlData, &temp); err != nil {
		log.Println("open config err: ", err)
		if fail {
			os.Exit(1)
		}
	}
	globals.Lock.Lock()
	globals.Config = temp
	globals.Lock.Unlock()
}

func main() {
	flag.StringVar(&tomlData, "config", "server.toml", "set the toml path")
	flag.Parse()

	configInit()

	log.Println(fmt.Sprintf("%+v", globals.Config))

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
	param := api.GetParam(r)
	log.Println(param)

	api.Database(param)

	for k, v := range database.Content {
		fmt.Println(k, v)
	}

	w.Header().Set("Content-Type", api.Response.Type)
	// w.Header().Set("Content-Type", "application/json")
	w.Write(api.Resp())
}
