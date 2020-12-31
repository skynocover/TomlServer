package database

import (
	"fmt"
)

var id = 0

// Content ...
var Content = make(map[int]map[string]string, 0)

// Insert ...
func Insert(param map[string]string) {
	Content[id] = param
	id++
}

// Scan ...
func Scan() {

	for k, v := range Content {

		fmt.Println(k, v)

	}

}
