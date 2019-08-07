package cmd

import (
	"bufio"
	"log"
	"os"

	input "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var startCmd = &cobra.Command{
	Use:   "start [name]",
	Short: "Runs auto withdraw program, with given key info",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		reader := bufio.NewReader(os.Stdin)
		password, err := input.GetPassword("Enter the passphrase:", reader)
		if err != nil {
			log.Fatalf("failed reading password: %s", err.Error())
			return
		}

		generator.KeyName = args[0]
		generator.KeyPassword = password

		kb, err := keys.NewKeyBaseFromDir(generator.KeyDir)
		if err != nil {
			log.Fatalf("failed to open keybase: %s", err.Error())
			return
		}

		_, err = kb.Get(generator.KeyName)
		if err != nil {
			log.Fatalf("failed to get account: %s", err.Error())
			return
		}

		generator.ListenNewBLock()
		// err := generator.SendTx(10000, "vodka")
		// fmt.Println(err)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
