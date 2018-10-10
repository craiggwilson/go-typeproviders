package cmd

import (
	"github.com/craiggwilson/go-typeproviders/pkg/generate"
)

func run(p generate.StructProvider) {
	pkg := rootCmd.PersistentFlags().Lookup("pkg").Value.String()
	generate.Generate(p, "", pkg)
}
