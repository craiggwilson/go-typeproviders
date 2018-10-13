package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().StringP("pkg", "", "", "the name of the package to hold the structs")
	rootCmd.PersistentFlags().BoolP("embedStructs", "", false, "embed structs instead of giving them names")
}

// Execute starts the application using the provided arguments.
func Execute(args []string) {
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "typeprovider",
	Short: "Generate structs based on data.",
	Long:  `Generate structs based on data. `,
}
