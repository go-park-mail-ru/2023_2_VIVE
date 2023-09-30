package main

import (
	"fmt"
	"models/app"
	"models/serverErrors"
)

func main() {
	err := app.Run()

	if err != nil {
		fmt.Printf("err: %v\n in detail: %v\n", serverErrors.SERVER_IS_NOT_RUNNUNG, err)
		return
	}
}
