package cmd

import (
	"fmt"
	"git.garena.com/frieda.hasanah/user_service/internal/populator"
	"github.com/spf13/cobra"
)

var populateCommand = &cobra.Command{
	Use:   "populate",
	Short: "Start Populating users in db",
	Run:   Populate,
}

func init() {
	rootCmd.AddCommand(populateCommand)
}
func Populate(cmd *cobra.Command, args []string) {
	fmt.Println("Populating users")
	populator.PopulateUsers(InitDBMySQL(), 9000000, 10000)
}
