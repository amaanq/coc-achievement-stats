package cmd

import (
	"encoding/json"
	"fmt"
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
	time.Sleep(time.Millisecond * 5000)
	for _, tag := range all_tags {
		fmt.Println(tag)
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
