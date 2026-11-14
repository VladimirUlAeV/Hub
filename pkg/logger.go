package pkg

import (
	"log/slog"
	"os"
)

func NewLogger(mode int, file *os.File) *slog.Logger{
	setingsHandler := slog.HandlerOptions{}
	switch mode{
	case 0:
		setingsHandler.Level = slog.LevelInfo
		setingsHandler.AddSource = false
	
	case 1:
		setingsHandler.Level = slog.LevelDebug
		setingsHandler.AddSource = true
	}

	handler := slog.NewTextHandler(file, &setingsHandler) //(файл для логов, структура с настройками хэндлера)
	return slog.New(handler)
}

func CreateNewLogFile() *os.File{
	fileRoute := "bin/app.log"
	file, err := os.OpenFile(fileRoute, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil{
		panic(err)
	}
	return file
}

/*
Строгая типизация (через конструкторы slog.Attr)

Работа с лог-записью:
slog.Debug("", slog.String("ключ", "текст"))


slog.String("ключ", "текст") — для строк
slog.Int("ключ", 42) / slog.Int64("ключ", 42) — для целых чисел
slog.Float64("ключ", 3.14) — для дробных чисел
slog.Bool("ключ", true) — для bool
slog.Duration("ключ", elapsed) — для time.Duration (удобно для замеров времени)
slog.Any("error", err) — для ошибок и всего остального

*/