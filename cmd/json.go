package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/craiggwilson/go-typeproviders/pkg/providers/json"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(jsonCmd)

	jsonCmd.Flags().StringP("name", "n", "AutoGenerated", "The name of the struct.")
}

var jsonCmd = &cobra.Command{
	Use:   "json [filename]",
	Short: "Generate structs based on a json file.",
	Long:  "Generate structs based on a json file.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var r io.Reader
		var structName string
		if len(args) == 1 {
			f, err := os.Open(args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer func() {
				_ = f.Close()
			}()
			r = f

			base := filepath.Base(args[0])
			ext := filepath.Ext(base)
			structName = base[0 : len(base)-len(ext)]
		} else {
			r = os.Stdin
		}

		if structName == "" || cmd.Flags().Lookup("name").Changed {
			structName = cmd.Flags().Lookup("name").Value.String()
		}

		cfg := json.Config{
			Input:      r,
			StructName: structName,
		}

		p := json.NewStructProvider(cfg)
		run(p)
	},
}
