package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/craiggwilson/typeproviders/internal/generators/mongodb"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(mongodbCmd)

	mongodbCmd.Flags().StringP("uri", "u", "localhost:27017", "The mongodb URI to read from.")
	mongodbCmd.Flags().StringP("database", "d", "", "The mongodb database to use.")
	mongodbCmd.Flags().StringP("collection", "c", "", "The mongodb collection to use.")

	mongodbCmd.MarkFlagRequired("database")
}

var mongodbCmd = &cobra.Command{
	Use:   "mongodb",
	Short: "Generate structs based on a mongodb collection.",
	Long:  "Generate structs based on a mongodb collection.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := mongodb.Config{
			URI:            cmd.Flags().Lookup("uri").Value.String(),
			DatabaseName:   cmd.Flags().Lookup("database").Value.String(),
			CollectionName: cmd.Flags().Lookup("collection").Value.String(),
			Package:        rootCmd.PersistentFlags().Lookup("pkg").Value.String(),
		}

		gen := mongodb.NewGenerator(cfg)

		var buf bytes.Buffer
		if err := gen.Write(&buf); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(buf.String())
	},
}
