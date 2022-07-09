package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/amaanq/coc-achievement-stats/log"
)

func WgetTags() {
	all_tags = nil
	all_tags = make([]string, 0)

	var url_wg sync.WaitGroup
	url_wg.Add(len(all_urls) * len(pages))

	log.Log.Info("[+] Getting top players..")

	for _, url := range all_urls {
		for _, page := range pages {
			go func(url string, page int) {
				defer url_wg.Done()
				URL := fmt.Sprintf(url+fmt.Sprint(page), th_level)

				resp, err := session.R().Get(URL)
				if err != nil {
					panic(err)
				}
				if resp.StatusCode() != 200 {
					return
					panic(fmt.Sprint(resp.StatusCode()) + " " + fmt.Sprintf(url+fmt.Sprint(page), th_level))
				}

				// unmarshal into Response
				var response Response
				err = json.Unmarshal(resp.Body(), &response)
				if err != nil {
					panic(err)
				}
				for _, ranking := range response.Rankings {
					all_tags = append(all_tags, ranking.Tag)
				}
			}(url, page)
		}
	}
	url_wg.Wait()

	log.Log.Info("[+] Finished retrieving top players tags.")
}

func WgetAchievements() {
	var achievement_wg sync.WaitGroup
	achievement_wg.Add(len(all_tags))

	log.Log.Info("[+] Getting top players achievements..")

	start := time.Now()

	for i, tag := range all_tags {
		time.Sleep(time.Millisecond * 10)
		go func(tag string, i int) {
			defer achievement_wg.Done()
			if i%100 == 0 || i == len(all_tags)-1 {
				defer updateProgress(start, i, len(all_tags))
				if i == len(all_tags)-1 {
					time.Sleep(time.Millisecond * 500)
				}
			}

			player, err := client.GetPlayer(tag)
			if err != nil && strings.Contains(err.Error(), "notFound") {
				return
			}
			for err != nil && !strings.Contains(err.Error(), "notFound") {
				time.Sleep(time.Millisecond * 100)
				fmt.Println("tag " + tag + ": " + err.Error())
				player, err = client.GetPlayer(tag)
			}

			if player.TownHallLevel != th_level {
				return
			}

			players = append(players, *player)
		}(tag, i)
	}
	achievement_wg.Wait()

	log.Log.Info("[+] Finished retrieving top players achievements.")
}

func saveToFile(fp string) error {
	log.Log.Info("[+] Saving top players to file..")

	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(players)
	if err != nil {
		return err
	}
	return nil
}

func updateProgress(start time.Time, idx, total int) {
	percent := float64(idx) / float64(total) * 100
	t := time.Now()
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	date := fmt.Sprintf("%d/%02d/%02d %02d:%02d:%02d", year, month, day, hour, min, sec)
	if percent < 99.9 {
		fmt.Printf("\033[2K\r\033[0;32m[INFO] \033[0;34m %s \033[0m%.2f%% %.2fs", date, percent, time.Since(start).Seconds())
	} else {
		fmt.Printf("\033[2K\r\033[0;32m[INFO] \033[0;34m %s \033[0m100%% %.2fs\n", date, time.Since(start).Seconds())
	}
}

type Response struct {
	Size     int       `json:"size"`
	Rankings []Ranking `json:"rankings"`
}

type Ranking struct {
	Tag                 string  `json:"tag"`
	Value               int     `json:"value"`
	Rank                int     `json:"rank"`
	Name                string  `json:"name"`
	CharacterID         string  `json:"characterId"`
	TownHallLevel       int     `json:"townHallLevel"`
	TownHallWeaponLevel any     `json:"townHallWeaponLevel"`
	BuilderHallLevel    int     `json:"builderHallLevel"`
	Clan                *Clan   `json:"clan,omitempty"`
	ClanTag             *string `json:"clanTag"`
	IsVip               *bool   `json:"isVip,omitempty"`
}

type Clan struct {
	Name  string `json:"name"`
	Tag   string `json:"tag"`
	Badge string `json:"badge"`
}
