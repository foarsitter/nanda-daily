package cmd

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var dbusCmd = &cobra.Command{
	Use:   "dbus",
	Short: "Monitor the dbus of changes in the lock screen",
	Long:  `The lock screen is part of the gnome ScreenSaver and can be used to detect if the user goes back to work`,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := dbus.ConnectSessionBus()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
			os.Exit(1)
		}
		defer conn.Close()

		if err = conn.AddMatchSignal(
			dbus.WithMatchInterface("org.gnome.ScreenSaver"),
			dbus.WithMatchMember("ActiveChanged"),
		); err != nil {
			panic(err)
		}

		c := make(chan *dbus.Signal, 10)
		conn.Signal(c)
		for v := range c {

			x := v.Body[0].(bool)

			dt := time.Now()

			if x {
				fmt.Println(dt.String(), "ActiveChange is set to True, as in the screensaver is turned on")
				DoFinish()
			} else {
				DoStart()
				fmt.Println(dt.String(), "ActiveChange is set to False, as in the screensaver is turned off")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(dbusCmd)
}
