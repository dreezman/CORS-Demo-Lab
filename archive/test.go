package main

import "fmt"

func main() {
	var m map[string]int
	m = make(map[string]int)
	m["route"] = 66
	j := m["root"]
	fmt.Println("value",j)
}