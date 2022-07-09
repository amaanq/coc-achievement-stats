package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/amaanq/coc.go"
	"github.com/go-resty/resty/v2"
)

var (
	th_level = 3

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

func _main() {
	var err error
	client, err = coc.New(map[string]string{"dummy1@yopmail.com": "Password"})
	if err != nil {
		panic(err)
	}

	GetCosURLs()

	fmt.Println("Saving to file...")

	// save to file
	fp := fmt.Sprintf("players-th%d.json", th_level)
	file, err := os.Create(fp)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(players)
	if err != nil {
		panic(err)
	}

	fmt.Println("Done!")
}

func GetCosURLs() {
	var url_wg sync.WaitGroup
	url_wg.Add(len(all_urls) * len(pages))

	for _, url := range all_urls {
		for _, page := range pages {
			go func(url string, page int) {
				defer url_wg.Done()
				fmt.Printf("%s%d\n", url, page)
				resp, err := session.R().Get(fmt.Sprintf(url+fmt.Sprint(page), th_level))
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
}

func GetTags() {
	var api_wg sync.WaitGroup
	api_wg.Add(len(all_tags))
	for i, tag := range all_tags {
		time.Sleep(time.Millisecond * 10)
		go func(tag string, i int) {
			defer api_wg.Done()
			if i%100 == 0 {
				defer fmt.Printf("%d/%d\n", i, len(all_tags))
			}

			player, err := client.GetPlayer(tag)
			if err != nil && strings.Contains(err.Error(), "notFound") {
				fmt.Printf("%s not found (banned)\n", tag)
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
	api_wg.Wait()
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
