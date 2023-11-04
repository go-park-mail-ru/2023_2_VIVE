package main

import "fmt"

func main1() {
	IDsNum := 5
	endOfQuery := ""
	for i := 0; i < IDsNum; i++ {
		if i == 0 {
			endOfQuery += fmt.Sprintf("$%d", i + 1)
		} else {
			endOfQuery += fmt.Sprintf(", $%d", i + 1)
		}
	}

	fmt.Printf("(%s)\n", endOfQuery)
}
