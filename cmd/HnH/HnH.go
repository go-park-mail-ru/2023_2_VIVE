package main

import (
	"HnH/app"
	"HnH/pkg/serverErrors"
	"fmt"
)

func main() {
	err := app.Run()

	if err != nil {
		fmt.Printf("err: %v\n in detail: %v\n", serverErrors.SERVER_IS_NOT_RUNNUNG, err)
		return
	}
}
