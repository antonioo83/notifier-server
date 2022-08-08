package commands

import (
	"github.com/antonioo83/notifier-server/config"
	"github.com/antonioo83/notifier-server/internal/services/client"
	"log"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Gets a status of the message",
	Long:  `Gets a status of the message.`,
	Run: func(cmd *cobra.Command, args []string) {
		configFromFile, err := config.LoadClientConfigFile("client_config.json")
		if err != nil {
			log.Fatalf("i can't load configuration file:" + err.Error())
		}
		cfg, err := config.GetClientConfigSettings(configFromFile)
		if err != nil {
			log.Fatalf("Can't read config: %s", err.Error())
		}

		messageID := cmd.Flag("s").Value.String()
		if messageID == "" {
			log.Println("Please set message ID using \"id\" flag!")
			return
		}

		ms := client.NewMessageService(cfg)
		resp, err := ms.GetStatus(messageID)
		if err != nil {
			log.Fatalf("Can't get a message status: %s", err.Error())
		}

		status := "sent!"
		if !resp.IsSent {
			status = "not sent"
		}

		log.Println("Message status:" + status)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.PersistentFlags().String("s", "1", "Message ID")
}
