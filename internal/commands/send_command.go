package commands

import (
	"github.com/antonioo83/notifier-server/config"
	"github.com/antonioo83/notifier-server/internal/services/client"
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Sends messages to the notifier server",
	Long:  `Sends messages to the notifier server. File have to json format`,
	Run: func(cmd *cobra.Command, args []string) {
		configFromFile, err := config.LoadClientConfigFile("client_config.json")
		if err != nil {
			log.Fatalf("i can't load configuration file:" + err.Error())
		}
		cfg, err := config.GetClientConfigSettings(configFromFile)
		if err != nil {
			log.Fatalf("Can't read config: %s", err.Error())
		}

		filepath := cmd.Flag("f").Value.String()
		if filepath == "" {
			filepath = "messages.json"
		}

		ms := client.NewMessageService(cfg)
		status, err := ms.SendMessages(filepath)
		if err != nil {
			log.Println("I can't send message: " + err.Error())
			return
		} else {
			log.Println("I sent request and get next HTTP status:" + strconv.Itoa(status))
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	sendCmd.PersistentFlags().String("f", "messages.json", "File path")
}
