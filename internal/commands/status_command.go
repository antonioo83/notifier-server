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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	statusCmd.PersistentFlags().String("s", "1", "A help for foo 1")
}
