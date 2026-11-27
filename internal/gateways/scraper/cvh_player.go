package scraper

import (
	"Hub/internal/entities"
	"Hub/internal/gateways"
	"context"
	"encoding/json"
	"io"

	"net/http"
)

type items struct{
	VkId string `json:"vkId"`
	VoiceStudio string`json:"voiceStudio"`
	VoiceType string `json:"voiceType"` 
	Season int `json:"season"`
	Episod int `json:"episod"`
}

type episodesJson struct{
	TitleName string `json:"titleName"`
	Items  []items `json:"items"`
}

type videoM3u8 struct{
	Status string `json:"status"`
	Source struct{
		HlsUrl string `json:"hlsUrl"`
	} `json:"sources"`
}

func CvhPlayerEpisodesParse(c *http.Client, url string, refer string, ctx context.Context)(*[]entities.Season ,error){
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil{
		return nil, err
	}

	r.Header.Set("accept","application/json")
	r.Header.Set("accept-language","ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	r.Header.Set("origin", refer)
	r.Header.Set("referer", refer)
	r.Header.Set("sec-ch-ua", `"Not-A.Brand";v="99", "Google Chrome";v="157", "Chromium";v="157"`)
	r.Header.Set("user-agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/157.0.0.0 Safari/537.36")

	res, err := c.Do(r)
	if err != nil{
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil{
		return nil, err
	}
	defer res.Body.Close()

	var episodesData episodesJson
	if err := json.Unmarshal(body, &episodesData); err != nil{
		return nil, err
	}
	
	var allSeason []entities.Season
	var actualSeason entities.Season
	countSeasons := -1
	for _, item := range episodesData.Items {
		if item.Season != countSeasons{
			countSeasons = item.Season 
			found := false
			for _, season := range allSeason{
				if season.Number == item.Season{
					actualSeason = season
					found = true 
					break
				}
			}
			if !found{
				actualSeason = entities.NewSeason(item.Season, make(map[string]entities.VoiceStudio))
				allSeason = append(allSeason, actualSeason)
			}
		
			
		}
		voice := item.VoiceStudio
		if voice == ""{
			voice = item.VoiceType
		}

		voiceStudio, ok := actualSeason.VoiceType[voice]
		if !ok {
			voiceStudio = entities.NewVoiceStudio(make(map[string]entities.Player))
			actualSeason.VoiceType[voice] = voiceStudio
		}

		player, ok := voiceStudio.Player["CVH"]
		if !ok{
			player = entities.NewPlayer(2, make(map[int]string))
			voiceStudio.Player["CVH"] = player
		}
		
		if _, ok = player.Episodes[item.Episod]; !ok{
			player.Episodes[item.Episod] = item.VkId
		}
	}
	
	
	return &allSeason, nil
}

func CvhPlayerVideoThread(refer string, idVideo string ,c * http.Client, ctx context.Context)(string, error){
	url := 	"https://plapi.cdnvideohub.com/api/v1/player/sv/video/" + idVideo
	
	r, err := http.NewRequestWithContext(ctx,http.MethodGet, url, nil)
	if err != nil{
		return "", err
	}

	r.Header.Set("accept", "application/json")
	r.Header.Set("accept-language","ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	r.Header.Set("origin", refer)
	r.Header.Set("sec-ch-ua",`"Not-A.Brand";v="99", "Google Chrome";v="157", "Chromium";v="157"`)
	r.Header.Set("user-agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/157.0.0.0 Safari/537.36")

	res, err := c.Do(r)
	if err != nil{
		return "", err 
	}
	
	body, err := io.ReadAll(res.Body)
	if err != nil{
		return "", err
	}

	var videoM3u8 videoM3u8

	if err := json.Unmarshal(body, &videoM3u8); err != nil{
		return "", err
	}
	
	if videoM3u8.Status != "ok"{
		return "", gateways.NotData
	}
	return videoM3u8.Source.HlsUrl, nil
}
