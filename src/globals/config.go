package globals

import (
	"log"
	"os"
	"sync"
	"time"
	"tomlserver/src/handler"

	"github.com/BurntSushi/toml"
)

// Config ...
var Config *tconfig

type tconfig struct {
	Filename       string
	LastModifyTime int64
	Lock           *sync.RWMutex
	Data           *handler.Data
}

// NewConfig ...
func NewConfig(filename string) {
	Config = &tconfig{
		Filename:       filename,
		Lock:           &sync.RWMutex{},
		Data:           &handler.Data{},
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
