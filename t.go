package main

import (
	"fmt"
	"net/url"
)

func main() {
	url, _ := url.Parse("https://example.com?foo=value%2C1&bar=2")

	fmt.Printf("%#v\n", url.Query())

	for key, value := range url.Query() {
		fmt.Printf("%v: %v\n", key, value)
	}
}
