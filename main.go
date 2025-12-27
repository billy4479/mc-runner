package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/billy4479/mc-runner/cmd"
)

//go:embed frontend/dist
var frontend embed.FS

func main() {
	fmt.Println(frontend.ReadDir("frontend/dist"))
	err := cmd.Run(frontend)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
