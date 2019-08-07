package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/terra-project/feegiver/utils"
	yaml "gopkg.in/yaml.v2"
)

var (
	// Path to config
	cfgFile string

	// The actual app config
	generator utils.Generator

	// Version for the application. Set via ldflags
	Version = "undefined"

	// Commit (git) for the application. Set via ldflags
	Commit = "undefined"

	// Branch (git) for the application. Set via ldflags
	Branch = "undefined"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "feegiver",
	Short: "An fee giver server for terra",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.feegiver/config.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	generator = utils.Generator{
		Version: Version,
		Commit:  Commit,
		Branch:  Branch,
	}

	var bz []byte

	if cfgFile != "" {
		var err error
		bz, err = ioutil.ReadFile(cfgFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cfgFile = fmt.Sprintf("%s/.feegiver/config.yaml", home)
		bz, err = ioutil.ReadFile(cfgFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	yaml.Unmarshal(bz, &generator)
}
