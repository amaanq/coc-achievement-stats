/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/amaanq/coc-achievement-stats/log"
	"github.com/amaanq/coc.go"
	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// var (
// 	allCmds = []*cobra.Command{
// 		downloadTHCmd,
// 		listFilesCmd,
// 	}
// )

// rootCmd represents the base command when called without any subcommands

var rootCmd = &cobra.Command{
	Use:   "coc-achievements",
	Short: "A tool to download and parse Clash of Clans achievements record holders for those that aren't visible on Clash of Stats",
	Long:  `To get started run one of the following commands below:`,
	// RunE: func(cmd *cobra.Command, args []string) error { }
}

func displayCommandsToRun() (string, error) {
	allCmds := []string{"download", "list"}

	templates := &promptui.SelectTemplates{
		Label: "		{{ .Name }}?",
		Active: "		     ↳ {{ .Name | cyan }}",
		Inactive: "			{{ .Name | cyan }}",
		Selected: "Selected: {{ .Name | red }}",
		Details: `
			Selected:
			{{ .Name }}
			`,
	}
	prompt := promptui.Select{
		Label:     "Which command do you want to run",
		Items:     allCmds,
		Templates: templates,
		Size:      10,
	}
	index, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return allCmds[index], nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fix.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	var err error

	// load dotenv
	err = godotenv.Load()
	if err != nil {
		log.Log.Errorf("Error loading .env file: %v", err)
	}

	if os.Getenv("email") == "" || os.Getenv("password") == "" {
		log.Log.Error("You must set your email and password in the .env file")
		os.Exit(1)
	}

	client, err = coc.New(map[string]string{os.Getenv("email"): os.Getenv("password")})
	if err != nil {
		panic(err)
	}

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
