package database

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
)

// Content ...
var Content = &sync.Map{}

// Insert ...
func Insert(value string) {
	id := uuid.New().String()
	Content.Store(id, value)
}

// Scan ...
func Scan() {
	Content.Range(func(key interface{}, value interface{}) bool { //遍歷需要使用func
		k := key.(string)
		v := value.(string)
		log.Println(fmt.Sprintf("key: %s, value: %s", k, v))
		return true //回傳true會繼續下一輪
	})
}
