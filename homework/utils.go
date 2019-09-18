package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

// SiteSearch - Делает запрос к массиву адресов и возвращает те, в теле которых найдена подстрока
func SiteSearch(needle string, urls []string) ([]string, error) {
	// Структура для безопасной групповой работы через несколько горутин
	group := struct {
		errgroup.Group          //< Запуск горутин с отловом возвращаемых ошибок
		sync.Mutex              //< Синхронизация горутин при доступе к объектам структуры
		urls           []string //< Срез подходящих ссылок
	}{
		urls: make([]string, 0, len(urls)),
	}

	// Для каждой ссылки запускаем горутину
	for _, u := range urls {
		url := u //< Сохраням значение переменной, т.к. u будет меняться
		group.Go(func() error {
			// Делаем GET запрос на адрес из ссылки
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			// Читаем тело ответа на запрос
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			// Если тело содержит искомую строку - добавляем в массив ответа
			if strings.Contains(string(body), needle) {
				group.Lock()
				group.urls = append(group.urls, url)
				group.Unlock()
			}

			return nil
		})
	}

	// Ожидаем завершения всех горутин
	err := group.Wait()
	return group.urls, err
}
