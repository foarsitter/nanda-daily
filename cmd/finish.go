package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"time"

	"github.com/spf13/cobra"
)

var finishCmd = &cobra.Command{
	Use:   "finish",
	Short: "Mark the current time as end of your day in Nanda",
	Long:  `Updates the finish_id timelog object or creates one`,
	Run:   FinishCommand,
}

func FinishCommand(cmd *cobra.Command, args []string) {
	DoFinish()
}

func DoFinish() {
	finishId := viper.GetString("finish_id")

	today := time.Now()

	if finishId == "" {
		res := PostTimelog(today, "ending the day")
		viper.Set("finish_id", res.Result.Id)
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		res := PatchTimelog(finishId, today, "Ending the day")
		print(res.Meta.Status)
	}

	fmt.Println("finish called")
}

func init() {
	rootCmd.AddCommand(finishCmd)
}
