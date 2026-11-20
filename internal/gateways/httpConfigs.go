package gateways

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"time"

	"golang.org/x/net/publicsuffix"
)

type Client struct{
	Client *http.Client
}

func NewClient() (*Client, error){ //Общий патерн клиента для всех запросов
	// Проверка в конфиг файле на настройки прокси
	CustomTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment, //значение по дефолту
		TLSClientConfig: &tls.Config{
			//InsecureSkipVerify: false, //проверка сертификатов сайта
			MinVersion: tls.VersionTLS12, //обязательное использование современного шифрования
		},
		TLSHandshakeTimeout: 5 * time.Second, // ограничение на TLS рукопажатие
		/*
		DisableKeepAlive: false, //Переиспользование сокетов, по умолчанию включено
		MaxIdleConns:        100,              // держать до 100 сокетов в пуле суммарно
		MaxIdleConnsPerHost: 20,               // держать до 20 сокетов в пуле для одного хоста
		MaxConnsPerHost:     50,               // не более 50 одновременных сокетов на хост
		IdleConnTimeout:     60 * time.Second, // закрыть сокет, если простаивает 60 секунд
		*/
		ForceAttemptHTTP2: true, // Включаем поддержку HTTP/2
	}
	CustomJar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List, //защита от утечек между субдоменами, также для иметаций браузера
	}) 
	if err != nil{
		return &Client{}, err
		//логирование или что-то еще
	}

	return &Client{
			Client: &http.Client{
			Transport: CustomTransport,
			//CheckRedirect: ,
			Jar: CustomJar,
			Timeout: 10 * time.Second,
		},
	}, nil 
}

func NewBasicRequst(ctx context.Context, url string)(*http.Request, error){  //Общий патерн запроса для палучения html-страницы
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil{
		return nil, err
	}

	r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Gecko/20100101 Firefox/132.0")
	r.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	r.Header.Set("Accept","text/html,application/xhtml+xml,*/*;q=0.8")


	return r, nil
}
