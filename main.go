package main

import (
	"os"

	"github.com/craiggwilson/typeproviders/cmd"
)

func main() {
	cmd.Execute(os.Args[1:])
}
