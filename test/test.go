package main

import "log"

func main() {
	params := make(map[string]int)

	params["s"] = 1
	params["y"] = 2

	for key, value := range params {
		log.Println(key)
		log.Println(value)
	}
}
