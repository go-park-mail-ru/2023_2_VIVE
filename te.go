package main

import (
	"HnH/pkg/nullTypes"
	"fmt"
)

func main() {
	fmt.Println(nullTypes.NewNullInt(10, true).Value())
}
