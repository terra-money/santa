package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/terra-project/santa/utils"
	yaml "gopkg.in/yaml.v2"
)

// versionCmd represents the version command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Sets a default config file",
	Run: func(cmd *cobra.Command, args []string) {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println("Error finding homedir:", err)
			return
		}

		g := utils.SantaApp{
			KeyDir: fmt.Sprintf("%s/.santa", home),
			Node:   "http://localhost:26657",

			TriggerInterval: "5",
			FeeAmount:       "1000000uluna",

			FailWebHookURL:     "",
			FailWebHookDataKey: "text",

			SuccessWebHookURL:     "",
			SuccessWebHookDataKey: "text",
		}

		if _, err := os.Stat(g.KeyDir); os.IsNotExist(err) {
			err := os.MkdirAll(g.KeyDir, 0777)
			if err != nil {
				fmt.Println("Error creating directory:", err)
				return
			}
		}

		conf := fmt.Sprintf("%s/config.yaml", g.KeyDir)
		if _, err := os.Stat(conf); os.IsNotExist(err) {
			out, err := yaml.Marshal(g)
			if err != nil {
				fmt.Println("Error marshaling config:", err)
				return
			}
			file, err := os.Create(conf)
			if err != nil {
				fmt.Println("Error creating config file:", err)
				return
			}
			defer file.Close()
			fmt.Fprintf(file, string(out))
		} else {
			fmt.Println("Config file already exists, skipping...")
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
