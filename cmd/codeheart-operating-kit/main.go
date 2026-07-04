package main

import (
	"os"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
