package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/amaanq/coc-achievement-stats/log"
	"github.com/amaanq/coc.go"
	"github.com/go-resty/resty/v2"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	th_level = -1

	pb_url  = "https://api.clashofstats.com/rankings/players/best-trophies?location=global&level=%d&page="
	ws_url  = "https://api.clashofstats.com/rankings/players/war-stars?location=global&level=%d&page="
	cwl_url = "https://api.clashofstats.com/rankings/players/war-league-legend?location=global&level=%d&page="
	atw_url = "https://api.clashofstats.com/rankings/players/attack-wins?location=global&level=%d&page="
	dfw_url = "https://api.clashofstats.com/rankings/players/defense-wins?location=global&level=%d&page="
	hh_url  = "https://api.clashofstats.com/rankings/players/heroic-heist?location=global&level=%d&page="
	cnq_url = "https://api.clashofstats.com/rankings/players/conqueror?location=global&level=%d&page="
	unb_url = "https://api.clashofstats.com/rankings/players/unbreakable?location=global&level=%d&page="
	hum_url = "https://api.clashofstats.com/rankings/players/humiliator?location=global&level=%d&page="
	gch_url = "https://api.clashofstats.com/rankings/players/games-champion?location=global&level=%d&page="
	don_url = "https://api.clashofstats.com/rankings/players/donations?location=global&level=%d&page="
	rcv_url = "https://api.clashofstats.com/rankings/players/donations-received?location=global&level=%d&page="
	fin_url = "https://api.clashofstats.com/rankings/players/friends-in-need?location=global&level=%d&page="
	exp_url = "https://api.clashofstats.com/rankings/players/exp-level?location=global&level=%d&page="
	wsn_url = "https://api.clashofstats.com/rankings/players/well-seasoned?location=global&level=%d&page="
	gob_url = "https://api.clashofstats.com/rankings/players/get-those-goblins?location=global&level=%d&page="
	nnt_url = "https://api.clashofstats.com/rankings/players/nice-and-tidy?location=global&level=%d&page="

	all_urls = []string{pb_url, ws_url, cwl_url, atw_url, dfw_url, hh_url, cnq_url, unb_url, hum_url, gch_url, don_url, rcv_url, fin_url, exp_url, wsn_url, gob_url, nnt_url}
	pages    = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	client  *coc.Client
	session = resty.New()

	all_tags = make([]string, 0)
	players  = make([]coc.Player, 0)
)

var downloadTHCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a town hall level's data",
	Long:  `Red is not downloaded already, green is downloaded already. You can redownload a town hall level by selecting it, and it will be downloaded again overwriting the old data.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		th_level, err = renderViewOfThs("Which Town Hall level do you want to download (Red are not downloaded already)")
		if err != nil {
			log.Log.Errorf("Error rendering view of ths: %v", err)
			return err
		}
		if th_level == -1 {
			log.Log.Error("Exiting gracefully")
			return nil
		}

		// check if players-th{th_level}.json exists
		if achievementsAlreadyExist() {
			log.Log.Info("Th" + strconv.Itoa(th_level) + " data already exists, prompting user to overwrite")
			if !askToOverwrite() {
				log.Log.Info("Aborting")
				return nil
			}
		}
		
		return downloadTHToFile()
	},
}

func downloadTHToFile() error {
	var err error
	log.Log.Info("Downloading th" + strconv.Itoa(th_level) + " data")

	WgetTags()
	WgetAchievements()

	fp := fmt.Sprintf("players-th%d.json", th_level)
	err = saveToFile(fp)
	if err != nil {
		log.Log.Errorf("Error saving to file %s: %v", fp, err)
		return err
	}
	log.Log.Info("Done downloading th" + strconv.Itoa(th_level) + " data")
	return nil
}

func achievementsAlreadyExist() bool {
	if _, err := os.Stat(filepath.Join(".", "players-th"+strconv.Itoa(th_level)+".json")); err == nil {
		return true
	}
	return false
}

func renderViewOfThs(_prompt string) (int, error) {
	ths := GetAvailableThsData()

	all_ths := make([]string, 0)
	_range := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, -1}

	// missing ths should have whats in range and not in ths
	for _, th := range _range {
		if th == -1 {
			all_ths = append(all_ths, "Exit")
			continue
		}

		if !contains(ths, th) {
			//missing_ths = append(missing_ths, fmt.Sprintf("\u001b[31mTH%d", th))
			all_ths = append(all_ths, fmt.Sprintf("\u001b[31mTH%d", th))
		} else {
			//notmissing_ths = append(notmissing_ths, fmt.Sprintf("\u001b[32mTH%d", th))
			all_ths = append(all_ths, fmt.Sprintf("\u001b[32mTH%d", th))
		}
	}
	all_ths = append(all_ths, "\u001b[31mExit")

	templates := &promptui.SelectTemplates{
		Label: "		{{ . }}?",
		Active: "		     ↳ {{ . | cyan }}",
		Inactive: "			{{ . | cyan }}",
		Selected: "Selected: {{ . | red }}",
		Details: `			
			Selected:
			{{ . }}
			`,
	}
	prompt := promptui.Select{
		Label:     _prompt,
		Items:     all_ths,
		Templates: templates,
		Size:      10,
	}
	index, _, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	return _range[index], nil
}

func askToDownload() bool {
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Do you want to download th%d data?", th_level),
		IsConfirm: true,
		Validate: func(input string) error {
			if input != "y" && input != "Y" && input != "n" && input != "N" {
				return errors.New("please enter y or n")
			}
			return nil
		},
	}
	result, err := prompt.Run()
	if err != nil {
		log.Log.Errorf("Error running prompt: %v", err)
		return false
	}
	return result == "y" || result == "Y"
}

func askToOverwrite() bool {
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Do you want to overwrite and redownload th%d data?", th_level),
		IsConfirm: true,
		Validate: func(input string) error {
			if input != "y" && input != "Y" && input != "n" && input != "N" {
				return errors.New("please enter y or n")
			}
			return nil
		},
	}
	result, err := prompt.Run()
	if err != nil {
		log.Log.Errorf("Error running prompt: %v", err)
		return false
	}
	return result == "y" || result == "Y"
}

func contains(s []string, e int) bool {
	for _, a := range s {
		th_level, _ := strconv.Atoi(a[2:])
		if th_level == e {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(downloadTHCmd)
}
