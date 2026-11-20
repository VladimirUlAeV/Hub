package scraper // Скрапер - это программа или скрипт для автоматического сбора (парсинга) данных с веб-сайтов.

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	//	"log/slog"
	"net/http"
	"strings"

	"Hub/internal/entities"
	"Hub/internal/gateways"

	"github.com/PuerkitoBio/goquery"
)

type ScraperAnimeGo struct {
	client *http.Client
}
type animeLDJSON struct{
	Name string `json:"name"`
	AlterName string `json:"alternateName"`
	Image string `json:"image"`
}
type JsonRes struct {
	Status string `json:"status"`
	Data   struct {
		Content string `json:"content"`
	} `json:"data"`
}

func NewScraperAnimeGo(c * gateways.Client) (*ScraperAnimeGo, error){
	return &ScraperAnimeGo{
		client: c.Client,
	}, nil
}

/*func validUrl(Url string) (string, error) { //проверка ссылки
	URL := strings.ReplaceAll(Url, " ", "")
	if URL == "" || len(URL) == 0 {
		return "", gateways.EmptyUrl
	}

	return URL, nil
}*/

func animeContentUrl (url string, c *http.Client, ctx context.Context)(*animeLDJSON, error){
	r, err := gateways.NewBasicRequst(ctx, url)
	if err != nil {
		//		slog.Error("Ошибка создания запроса", err)
		return nil, err
	}

	animePage, err := c.Do(r)
	if err != nil {
		//		slog.Error("Ошибка запроса", err)
		return nil, err
	}

	body := animePage.Body
	defer animePage.Body.Close()
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil{
		return nil, err
	}

	rawLDJson := strings.TrimSpace(doc.Find(`script[type="application/ld+json"]`).First().Text())
	if rawLDJson == ""{
		return nil, gateways.NotData
	}
	var meta animeLDJSON
	if err := json.Unmarshal([]byte(rawLDJson), &meta); err != nil{
		return nil, err
	}

	return &meta, nil
}

func playerContent(refer string, url string, c *http.Client, ctx context.Context) ([]string, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Accept", "application/json")
	r.Header.Set("accept-language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	r.Header.Set("referer", refer)
	r.Header.Set("sec-ch-ua", `"Not-A.Brand";v="99", "Google Chrome";v="157", "Chromium";v="157"`)
	r.Header.Set("sec-fetch-site", "same-origin")
	r.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/157.0.0.0 Safari/537.36")
	r.Header.Set("x-requested-with", "XMLHttpRequest")
	r.Header.Set("sec-fetch-mode", "cors")
	r.Header.Set("sec-fetch-dest", "empty")
	r.Header.Set("sec-ch-ua-mobile", "?0")
	r.Header.Set("priority", "u=1, i")
	r.Header.Set("sec-ch-ua-platform", "Windows")
	
	resDataContent, err := c.Do(r)
	if err != nil {
		return nil, err
	}
	defer resDataContent.Body.Close()
	if resDataContent.StatusCode == http.StatusNotFound {
		return nil, gateways.NotFound
	}
	body, err := io.ReadAll(resDataContent.Body)
	if err != nil {
		return nil, err
	} 
	var JsonRes JsonRes

	if err := json.Unmarshal(body, &JsonRes); err != nil {
		return nil, err
	} 

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(JsonRes.Data.Content))
	if err != nil {
		return nil, err
	}

	const prefixLen = 15
	var links []string

	outer: // Метка для внешнего цикла
	for _, btn := range doc.Find("button[data-player]").EachIter() {
		if len(links) >= 4 { 
			break
		}

		link, exists := btn.Attr("data-player")
		if !exists || len(link) < prefixLen {
			continue
		}

		prefixLink := link[:prefixLen]
		for _, saved := range links {
			if prefixLink == saved[:prefixLen] {
				continue outer // Переходим сразу к следующей кнопке во внешнем цикле
		}
		}

		links = append(links, link)
	}
	
	return links, nil
}

//func apiPlayer() {}

//func mediaRequst() {}

func (s *ScraperAnimeGo) AnimeParser(ctx context.Context, url string) (entities.Anime, error) {
	ValidUrl := url
/*	if err != nil {
		return entities.Anime{}, err
	} */
	IdAnimeSlice := []rune(url)
	if len(IdAnimeSlice) <= 4 {
		return entities.Anime{}, gateways.NotFullId
	}
	IdAnime := string(IdAnimeSlice[len(IdAnimeSlice)-4:]) // получение последних 4 символов из ссылки
	playerUrl := "https://animego.me/player/" + IdAnime

	client := s.client
	
	
	nameAndImage, err := animeContentUrl(ValidUrl,client,ctx); 
	if err != nil {
		return entities.Anime{}, err
	}

	players, err := playerContent(url, playerUrl, client, ctx) 
	if err != nil {
		return entities.Anime{}, err
		//slog.Error("Ошибка парсинга данных", err)

	}

//	slog.Info("Парсинг", slog.String("страницы", url), slog.String("статус", "успешно"))
	return entities.Anime{}, nil
}
