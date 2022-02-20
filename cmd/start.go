package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"time"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Create a entry at start",
	Long:  ``,
	Run:   StartCommand,
}

type PostTimelogResponse struct {
	Meta struct {
		Time          float64       `json:"time"`
		Log           []interface{} `json:"log"`
		TransactionId string        `json:"transaction_id"`
		Status        int           `json:"status"`
		Dt            float64       `json:"dt"`
	} `json:"meta"`
	Result struct {
		Id           string        `json:"id"`
		ReadGroup    string        `json:"read_group"`
		Account      string        `json:"account"`
		Project      interface{}   `json:"project"`
		Organization interface{}   `json:"organization"`
		RangeFrom    time.Time     `json:"range_from"`
		RangeUntil   time.Time     `json:"range_until"`
		DateOnly     int           `json:"date_only"`
		Duration     int           `json:"duration"`
		Date         string        `json:"date"`
		Time         string        `json:"time"`
		Year         int           `json:"year"`
		Month        int           `json:"month"`
		Day          int           `json:"day"`
		Week         int           `json:"week"`
		Weekday      int           `json:"weekday"`
		Hour         int           `json:"hour"`
		Description  string        `json:"description"`
		LastModified time.Time     `json:"last_modified"`
		Labels       []interface{} `json:"labels"`
	} `json:"result"`
}

func StartCommand(cmd *cobra.Command, args []string) {

	DoStart()
}

func DoStart() {
	lastSyncDate := viper.GetString("last_sync_date")

	today := time.Now()

	if lastSyncDate == today.Format("2006-01-02") {
		fmt.Println("The day is already synced")
		return
	}

	res := PostTimelog(today, "Starting the day")

	if res.Meta.Status == 201 {
		viper.Set("start_id", res.Result.Id)
		viper.Set("finish_id", "")

		viper.Set("last_sync_date", today.Format("2006-01-02"))
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("start called")
}

func init() {
	rootCmd.AddCommand(startCmd)
}
