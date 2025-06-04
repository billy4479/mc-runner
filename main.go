package main

import (
	"fmt"
	"os"

	"github.com/billy4479/mc-runner/cmd"
)

func main() {
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
