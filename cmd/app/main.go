package main

import (
	"Hub/internal/gateways"
	"Hub/internal/gateways/scraper"
	"Hub/pkg"
	"context"
	"log/slog"
	"time"
)


func main() {
	ctx,cancel := context.WithTimeout(context.Background(), 1 * time.Second)
	defer cancel()

	httpClient, err := gateways.NewClient()
	if err != nil{
		return
	}

	LogFile := pkg.CreateNewLogFile()
	slog.SetDefault(pkg.NewLogger(0, LogFile)) //логгер(режим, файл для записи логов)

	anime, err := scraper.NewScraperAnimeGo(httpClient)
	if err != nil{
		return
	}
	_, err = anime.AnimeParser(ctx,"https://animego.me/anime/saga-o-vinlande-1089")
	if err !=nil{
		return 
	}
	
}
