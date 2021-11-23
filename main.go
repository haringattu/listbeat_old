package main

import (
	"os"

	"github.com/haringattu/listbeat/cmd"

	_ "github.com/haringattu/listbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
