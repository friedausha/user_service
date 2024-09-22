package cmd

import (
	"fmt"
	"git.garena.com/frieda.hasanah/user_service/internal/data"
	"git.garena.com/frieda.hasanah/user_service/internal/data/cache"
	"git.garena.com/frieda.hasanah/user_service/internal/model"
	"git.garena.com/frieda.hasanah/user_service/internal/service"
	"git.garena.com/frieda.hasanah/user_service/utils/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "User Service",
	Short: "A service to handle everything related to user data",
	Long:  `A service to handle everything related to user data, including registration, authentication and authorizationF`,
}

var (
	userRepository model.IUserRepository
	userService    model.IUserService
	authService    model.IAuthService
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func init() {
	cobra.OnInitialize(initConfigReader, initApp)
}

func initConfigReader() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Print(".env file not found")
		}
	}
}

func initApp() {
	dbConn := InitDBMySQL()
	userRepository = data.NewUserRepository(dbConn)

	userCache := cache.NewUserCache()

	authService = service.NewAuthService(userRepository, userCache)
	userService = service.NewService(userRepository, userCache)
}
