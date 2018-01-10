package main

import (
	"fmt"

	"github.com/uphy/doopa/adapter/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println("Failed to execute the command: " + err.Error())
	}
}
