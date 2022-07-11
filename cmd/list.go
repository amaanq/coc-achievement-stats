package cmd

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var listFilesCmd = &cobra.Command{
	Use:   "list",
	Short: "List town hall levels with available data",
	Long:  `If you have downloaded some town hall levels' data but aren't sure which ones, you can use this command to list the town hall levels with available data.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ths := GetAvailableThsData()

		fmt.Println("Town hall levels with available data:")
		for _, th := range ths {
			fmt.Println("    ", th)
		}

		return rootCmd.Execute()
	},
}

func GetAvailableThsData() []string {
	ths := make([]string, 0)

	filepath.WalkDir("./", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".json" {
			// print path without extension
			n := strings.Index(path, "-")
			if n == -1 {
				return nil
			}
			//fmt.Println(strings.ToUpper(path[n+1 : len(path)-5]))
			ths = append(ths, strings.ToUpper(path[n+1:len(path)-5]))
		}
		return nil
	})
	sort.SliceStable(ths, func(i, j int) bool {
		// parse last digit or two digits
		i1, _ := strconv.Atoi(ths[i][2:])
		j1, _ := strconv.Atoi(ths[j][2:])
		return i1 < j1
	})
	return ths
}

func init() {
	rootCmd.AddCommand(listFilesCmd)
}
