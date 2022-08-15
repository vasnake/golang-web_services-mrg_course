package main

import (
	"encoding/json"
	"fmt"
)

var jsonStr = `[
	{"id": 17, "username": "iivan", "phone": 0},
	{"id": "17", "address": "none", "company": "Mail.ru"}
]` // list of map, N.B. id data type

func main() {
	data := []byte(jsonStr)

	var user1 interface{}
	json.Unmarshal(data, &user1)
	fmt.Printf("unpacked in empty interface:\n%#v\n\n", user1) // list of map string->interface{}

	user2 := map[string]interface{}{
		"id":       42,
		"username": "rvasily",
	}
	var user2i interface{} = user2 // cast to empty interface, just for fun
	result, _ := json.Marshal(user2i)
	fmt.Printf("json string from map:\n %s\n", string(result)) // valid json
}
