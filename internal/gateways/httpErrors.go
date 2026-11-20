package gateways

import "errors"

var (
	BadRequst = errors.New("Плохой запрос")
	NotFullId = errors.New("Не полное айди")
	NotFound = errors.New("Не найдено")
	NotData = errors.New("Нет данных")
)