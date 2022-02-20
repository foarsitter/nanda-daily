package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
	"net/http"
	"net/url"
	"os"
	"strings"
	"syscall"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

func credentials() (string, string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	username := viper.GetString("username")
	domain := viper.GetString("domain")

	fmt.Printf("Enter your nanda domain [%s]: ", domain)
	domain1, err := reader.ReadString('\n')
	if err != nil {
		return "", "", "", err
	}

	if domain1 != "\n" {
		domain = domain1
	}

	fmt.Printf("Enter Username [%s]: ", username)
	username1, err := reader.ReadString('\n')

	if username1 != "\n" {
		username = username1
	}

	if err != nil {
		return "", "", "", err
	}

	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", "", err
	}

	password := string(bytePassword)
	return strings.TrimSpace(domain), strings.TrimSpace(username), strings.TrimSpace(password), nil
}

var loginCmd = &cobra.Command{
	Use:   "LoginCommand",
	Short: "Login with your domain, e-mail and password",
	Run:   LoginCommand,
}

type LoginResponse struct {
	Meta struct {
		Time   float64       `json:"time"`
		Log    []interface{} `json:"log"`
		Status int           `json:"status"`
		Dt     float64       `json:"dt"`
	} `json:"meta"`
	Result struct {
		Account struct {
			ID                string   `json:"id"`
			Name              string   `json:"name"`
			Email             string   `json:"email"`
			Activated         int      `json:"activated"`
			Archived          int      `json:"archived"`
			TotalBookedCount  int      `json:"total_booked_count"`
			TotalBookedSum    int      `json:"total_booked_sum"`
			RecentBookedCount int      `json:"recent_booked_count"`
			Rank              int      `json:"rank"`
			Groups            []string `json:"groups"`
			Permissions       []string `json:"permissions"`
		} `json:"account"`
		AccessToken string `json:"access_token"`
		UIConfig    struct {
		} `json:"ui_config"`
	} `json:"result"`
	Error []struct {
		Code      string        `json:"code"`
		Message   string        `json:"message"`
		Stack     []interface{} `json:"stack"`
		Reference interface{}   `json:"reference"`
	} `json:"error"`
}

func LoginCommand(cmd *cobra.Command, args []string) {

	domain, username, password, _ := credentials()

	response := DoLogin(username, domain, password)

	viper.Set("domain", domain)
	viper.Set("access_token", response.Result.AccessToken)
	viper.Set("username", username)
	err := viper.WriteConfig()
	if err != nil {
		fmt.Println(err)
	}
}

func DoLogin(username string, domain string, password string) LoginResponse {
	data := url.Values{"email": {username}, "domain": {domain}, "password": {password}}

	resp, err := http.PostForm(requestUrl("auth/LoginCommand"), data)

	if err != nil {
		fmt.Println(err)
	}

	var res LoginResponse

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("AccessToken:", res.Result.AccessToken)
	fmt.Println("Done")

	return res
}

func requestUrl(path string) string {
	return fmt.Sprintf("%s/%s", "https://api.nanda.io/api/v1/", path)
}
