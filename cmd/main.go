package main

import (
	"fmt"
	"os"

	"github.com/leandrowiemesfilho/go-todo-cli/internal/cli"
)

func main() {
	app := cli.NewCLI()

	if err := app.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
