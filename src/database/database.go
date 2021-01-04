package database

import (
	"fmt"
)

// Content ...         table      key        values.......
var Content = map[string]map[string]map[string]string{}

// Insert ...
func Insert(table, key string, value map[string]string) error {
	row, ok := Content[table]
	if !ok {
		row = map[string]map[string]string{}
	}

	_, ok1 := row[key]
	if !ok1 {
		row[key] = value
		Content[table] = row
		return nil
	}
	return fmt.Errorf("Column already exist")

}

// Update ...
func Update(table, key string, value map[string]string) error {
	row, ok := Content[table]
	if !ok {
		return fmt.Errorf("table now exist")
	}

	_, ok = row[key]
	if !ok {
		return fmt.Errorf("row now exist")
	}
	Content[table][key] = value
	return nil
}

// Patch ...
func Patch(table, key string, value map[string]string) error {
	row, ok := Content[table]
	if !ok {
		return fmt.Errorf("table now exist")
	}

	columns, ok := row[key]
	if !ok {
		return fmt.Errorf("row now exist")
	}

	for k := range columns {
		for k1, v1 := range value {
			if k == k1 {
				Content[table][key][k1] = v1
			}
		}
	}
	return nil
}

// Delete ...
func Delete(table, key string) {
	delete(Content[table], key)
}

// ReadAll ...
func ReadAll(table string) map[string]map[string]string {
	return Content[table]
}

// Read ...
func Read(table, key string) map[string]string {
	return (Content[table])[key]
}

// Scan ...
func Scan() {

	for k, v := range Content {

		fmt.Println(k, v)

	}

}
