package cmd

import (
	"github.com/craiggwilson/go-typeproviders/pkg/providers/mongodb"
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
		}

		p := mongodb.NewStructProvider(cfg)
		run(p)
	},
}
