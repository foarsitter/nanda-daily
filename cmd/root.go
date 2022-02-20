package cmd

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "nanda-daily",
	Short: "Tool for logging the start en finish times of your work day to nanda",
	Long:  `Based on your lockscreen this tool determines your first login and last login of each day. It create two entries in Nanda marking your workday so you can fill the gap later.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nanda-daily.toml)")
}

var cfgFile = "$HOME/.nanda-daily.toml)"

func initConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".nanda-daily")
		viper.SetConfigType("toml")

		//err = viper.WriteConfigAs(home + "/.nanda-daily.toml")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
