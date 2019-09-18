package main

// ErrorModel - Ошибка отвечаемая сервером
type ErrorModel struct {
	Code     int         `json:"code"`
	Err      string      `json:"error"`
	Desc     string      `json:"desc"`
	Internal interface{} `json:"internal"`
}

// SitesSearchReqModel - модель входящих данных для поиска по сайтам
type SitesSearchReqModel struct {
	Query string   `json:"search"`
	Sites []string `json:"sites"`
}

// SitesSearchRespModel - модель исходящих данных для поиска по сайтам
type SitesSearchRespModel struct {
	Sites []string `json:"sites"`
}
