package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"github.com/amaanq/coc-achievement-stats/log"
	"github.com/amaanq/coc.go"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var mainStatsCmd = &cobra.Command{
	Use:   "main",
	Short: "main area to compare data across all achievements",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
	TH:
		th_level, err = renderViewOfThs("Which th level do you want to check out?")
		if err != nil {
			return err
		}

		if th_level == -1 {
			log.Log.Error("Exiting gracefully")
		}

		if !achievementsAlreadyExist() {
			log.Log.Info("Th" + strconv.Itoa(th_level) + " data doesn't exist, prompting user to download")
			if !askToDownload() {
				log.Log.Info("Aborting")
				return nil
			}

			err = downloadTHToFile()
			if err != nil {
				log.Log.Errorf("Error downloading th data: %v", err)
				return err
			}
		}

	ACHIEVEMENT:
		achievement, err := renderAchievementSelection()
		if err != nil {
			return err
		}

		if achievement == "Go back to TH level selection" {
			goto TH
		} else if achievement == "Exit" {
			log.Log.Info("Exiting gracefully")
			return nil
		} else {
			err = renderCompareAchievement(achievement)
			if err != nil {
				log.Log.Errorf("Error comparing achievement: %v", err)
				return err
			}
			goto ACHIEVEMENT
		}
	},
}

func renderAchievementSelection() (string, error) {
	achievements := []string{
		"Anti-Artillery",
		"Bust This!",
		"Clan War Wealth",
		"Firefighter",
		"Not So Easy This Time",
		"Shattered and Scattered",
		"X-Bow Exterminator",
		"Go back to TH level selection",
		"Exit",
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F4A5 {{ . | green }}",
		Inactive: "  {{ . | cyan }}",
		Selected: "{{ . | red }}",
		Details: `
--------- {{ . | cyan }} ----------
`,
	}
	prompt := promptui.Select{
		Label:     "Which achievement do you want to check out the rankings for?",
		Items:     achievements,
		Templates: templates,
		Size:      10,
	}
	index, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return achievements[index], nil
}

func renderCompareAchievement(achievement string) error {
	log.Log.Info("Comparing achievement " + achievement + " for th level " + fmt.Sprint(th_level))

	data, err := ioutil.ReadFile("players-th" + strconv.Itoa(th_level) + ".json")
	if err != nil {
		log.Log.Errorf("Error reading players-th%v.json: %v", th_level, err)
		return err
	}
	err = json.Unmarshal(data, &players)
	if err != nil {
		log.Log.Errorf("Error unmarshalling players-th%v.json: %v", th_level, err)
		return err
	}

	inCollapsed := make(map[coc.PlayerTag]bool)
	collapsedPlayers := make([]CollapsedPlayerStruct, len(players))
	for i, player := range players {
		if _, ok := inCollapsed[player.Tag]; !ok {
			_achievement := getAchievementByName(player, achievement)
			if _achievement == nil {
				continue
			}
			c := CollapsedPlayerStruct{
				Name: player.Name,
				Tag:  string(player.Tag),
				TH:                        player.TownHallLevel,
				AchievementName:           _achievement.Name,
				AchievementCompletionInfo: _achievement.CompletionInfo,
				AchievementValue:          _achievement.Value,
			}
			collapsedPlayers[i] = c
			inCollapsed[player.Tag] = true
		}
	}

	sort.SliceStable(collapsedPlayers, func(i, j int) bool {
		// sort by achievement name value descending
		return collapsedPlayers[i].AchievementValue > collapsedPlayers[j].AchievementValue
	})
	for i, cPlayer := range collapsedPlayers {
		collapsedPlayers[i].RenderedName = fmt.Sprintf("#%d. %s (%s) %d", i+1, cPlayer.Name, cPlayer.Tag, cPlayer.AchievementValue)
	}

	templates := &promptui.SelectTemplates{
		Label:    fmt.Sprintf("%s rankings for TH%d", achievement, th_level),
		Active:   "\U0001F4A5 {{ .RenderedName | green }}",
		Inactive: "  {{ .RenderedName | cyan }}",
		Selected: "{{ .RenderedName | red }}",
		Details: `
--------- {{ .RenderedName | cyan }} ----------
`,
	}

	prompt := promptui.Select{
		Label:     "Hit enter to exit", //fmt.Sprintf("%s rankings for TH%d", achievement, th_level),
		Items:     collapsedPlayers,
		Templates: templates,
		Size:      10,
	}
	_, _, err = prompt.Run()
	return nil
}

func getAchievementByName(player coc.Player, achievement string) *coc.Achievement {
	for _, a := range player.Achievements {
		if strings.EqualFold(a.Name, achievement) {
			return &a
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(mainStatsCmd)
}

type CollapsedPlayerStruct struct {
	Name                      string
	Tag                       string
	RenderedName              string
	TH                        int
	AchievementName           string
	AchievementCompletionInfo string
	AchievementValue          int
}
