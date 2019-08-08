package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	bip39 "github.com/bartekn/go-bip39"
	input "github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagRecover   = "recover"
	flagOldHdPath = "old-hd-path"
	flagBech      = "bech"
)

// versionCmd represents the version command
var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Runs keys calls",
}

// get keys list
var keysList = &cobra.Command{
	Use:   "list",
	Short: "Fetch all keys managed by the server",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := app.GetKeys()
		if err != nil {
			log.Fatalf("Failed: %s", err.Error())
			return
		}
		fmt.Println(string(out))
	},
}

// keys add
var keysAdd = &cobra.Command{
	Use:   "add [name]",
	Args:  cobra.ExactArgs(1),
	Short: "Add a new key to the keyserver, optionally pass a mnemonic to restore the key",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		password, err := input.GetCheckPassword(
			"Enter a passphrase to encrypt your key to disk:",
			"Repeat the passphrase:", reader)

		if err != nil {
			log.Fatalf("failed reading password: %s", err.Error())
			return
		}

		var mnemonic string
		oldHdPath := false
		if viper.GetBool(flagRecover) {
			bip39Message := "Enter your bip39 mnemonic"
			mnemonic, err = input.GetString(bip39Message, reader)
			if err != nil {
				log.Fatalf("failed reading mnemonic: %s", err.Error())
				return
			}

			if !bip39.IsMnemonicValid(mnemonic) {
				log.Fatal("invalid mnemonic")
				return
			}

			oldHdPath = viper.GetBool(flagOldHdPath)
		} else if viper.GetBool(flagOldHdPath) {
			log.Fatal("--old-hd-path can not be used without --recover flag")
			return
		}

		out, err := app.AddNewKey(args[0], password, mnemonic, oldHdPath)

		if err != nil {
			log.Fatalf("Failed: %s", err.Error())
			return
		}
		fmt.Println(string(out))
	},
}

// key show
var keyShow = &cobra.Command{
	Use:   "show [name]",
	Args:  cobra.ExactArgs(1),
	Short: "Fetch details for one key",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := app.GetKey(args[0], viper.GetString(flagBech))

		if err != nil {
			log.Fatalf("Failed: %s", err.Error())
			return
		}
		fmt.Println(string(out))
	},
}

// /keys/{name} PUT
var keyPut = &cobra.Command{
	Use:   "put [name]",
	Args:  cobra.ExactArgs(1),
	Short: "Update the password on a key",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		password, err := input.GetPassword("Enter the current passphrase:", reader)
		if err != nil {
			log.Fatalf("failed reading password: %s", err.Error())
			return
		}

		newPassword, err := input.GetCheckPassword(
			"Enter the new passphrase:",
			"Repeat the new passphrase:", reader)

		if err != nil {
			log.Fatalf("failed reading new password: %s", err.Error())
			return
		}
		err = app.UpdateKey(args[0], password, newPassword)
		if err != nil {
			log.Fatalf("Failed: %s", err.Error())
			return
		}
		fmt.Println("ok")
	},
}

// /keys/{name} DELETE
var keyDelete = &cobra.Command{
	Use:   "delete [name]",
	Args:  cobra.ExactArgs(1),
	Short: "Delete a key",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		password, err := input.GetPassword("Enter the passphrase:", reader)
		if err != nil {
			log.Fatalf("failed reading password: %s", err.Error())
			return
		}

		err = app.DeleteKey(args[0], password)
		if err != nil {
			log.Fatalf("Failed: %s", err.Error())
			return
		}
		fmt.Println("ok")
	},
}

func init() {
	keysCmd.AddCommand(keysList)

	keysAdd.Flags().Bool(flagRecover, false, "Recovering key option; mnemonic is required")
	keysAdd.Flags().Bool(flagOldHdPath, false, "Recover key with old hd path")
	viper.BindPFlag(flagRecover, keysAdd.Flags().Lookup(flagRecover))
	viper.BindPFlag(flagOldHdPath, keysAdd.Flags().Lookup(flagOldHdPath))
	keysCmd.AddCommand(keysAdd)

	keyShow.Flags().String(flagBech, "", "bech32 prefix; acc or val or cons")
	viper.BindPFlag(flagBech, keyShow.Flags().Lookup(flagBech))
	keysCmd.AddCommand(keyShow)
	keysCmd.AddCommand(keyPut)
	keysCmd.AddCommand(keyDelete)
	rootCmd.AddCommand(keysCmd)

}
