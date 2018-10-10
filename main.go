package main

import (
	"os"

	"github.com/craiggwilson/go-typeproviders/cmd"
)

func main() {
	cmd.Execute(os.Args[1:])
}
