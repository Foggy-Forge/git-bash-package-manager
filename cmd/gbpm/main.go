package main

import (
	"log"

	"github.com/Foggy-Forge/git-bash-package-manager/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
