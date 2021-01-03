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

// Update ...
func Update(key int, value map[string]string) {
	Content[key] = value
}

// Patch ...
func Patch(key int, value map[string]string) {
	for k := range Content[key] {
		for k1, v1 := range value {
			if k == k1 {
				Content[key][k] = v1
			}
		}
	}
}

// Read ...
func Read(key int) map[string]string {
	return Content[key]
}

// Scan ...
func Scan() {

	for k, v := range Content {

		fmt.Println(k, v)

	}

}
