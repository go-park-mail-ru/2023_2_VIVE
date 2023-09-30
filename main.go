package main

import (
	"fmt"
	"models/app"
	"models/errors"
)

func main() {
	err := app.Run()

	if err != nil {
		fmt.Printf("err: %v\n in detail: %v\n", errors.SERVER_IS_NOT_RUNNUNG, err)
		return
	}
}
