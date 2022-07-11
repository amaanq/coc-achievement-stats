/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/amaanq/coc.go"
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
	client, err = coc.New(map[string]string{"dummy1@yopmail.com": "Password"})
	if err != nil {
		panic(err)
	}

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
