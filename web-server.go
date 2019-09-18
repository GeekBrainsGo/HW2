package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"sync"

	"github.com/Sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type SearchStruct struct {
	Pattern string   `json:"search"`
	URLs    []string `json:"sites,omitempty"`
}

func main() {
	stopchan := make(chan os.Signal)

	// logrus.SetReportCaller(true)

	router := http.NewServeMux()

	router.HandleFunc("/", searchHandler)

	go func() {
		logrus.Info("Сервер запущен")
		err := http.ListenAndServe(":8080", router)
		logrus.Fatal(err)
	}()

	signal.Notify(stopchan, os.Kill, os.Interrupt)

	runPost() // запустим для проверки

	<-stopchan

	logrus.Info("Сервер остановлен!")
}

// runPost функция проверки сервера
//
func runPost() {

	data := SearchStruct{}
	data.Pattern = "Вакансии"
	data.URLs = []string{"http://www.mail.ru", "http://www.ibm.ru", "http://www.yandex.ru", "http://www.hp.ru"}
	url := "http://127.0.0.1:8080"

	obj, err := json.Marshal(data)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(obj))
	if err != nil {
		logrus.Info("Произошла ошибка передачи запроса: " + err.Error())
		return
	}
	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close() // важный пункт!

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Info("Произошла ошибка чтения ответа: " + err.Error())
			return
		}
	
		answer := SearchStruct{}
		err = json.Unmarshal(respBody, &answer)
	
		// выводим найденную информацию
		fmt.Println("Ищем:", answer.Pattern)
		for _, item := range answer.URLs {
			fmt.Println("Найден сайт:", item)
		}
		} else {
			fmt.Println("Ответ сервера:", resp.Status)
		}
}

func searchHandler(wr http.ResponseWriter, req *http.Request) {
	var resultSearch []string

	defer req.Body.Close() // важный пункт!
	byteJSON, _ := ioutil.ReadAll(req.Body)
	search := SearchStruct{}

	if req.Method != "POST" {
		http.Error(wr, "Неверный метод запроса", http.StatusBadRequest)
		logrus.Info("Неверный метод запроса: " + req.Method)
		return
	}

	if err := json.Unmarshal(byteJSON, &search); err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		logrus.Info("Неверная структура запроса: " + err.Error())
		return
	}

	if search.Pattern == "" {
		http.Error(wr, "Пустой шаблон поиска", http.StatusBadRequest)
		return
	}
	if search.URLs == nil {
		http.Error(wr, "Пустой список сайтов поиска", http.StatusBadRequest)
		return
	}

	resultSearch, err := multiSearch(search.URLs, search.Pattern)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		logrus.Info("Ошибка в поиске шаблона: " + err.Error())
		return
	}

	search.URLs = resultSearch

	answer, err := json.Marshal(search)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		logrus.Info("Ошибка в формировании ответа")
		return
	}
	fmt.Fprintf(wr, "%s", []byte(answer))

}

// multiSearch - функция поиска заданой строки на ресурсах представленных в массиве URL
// Изменено с учетом дополнительной информацией полученной на уроке
func multiSearch(arrLink []string, pattern string) ([]string, error) {

	// Структура для безопасной групповой работы через несколько горутин
	group := struct {
		errgroup.Group          //< Запуск горутин с отловом возвращаемых ошибок
		sync.Mutex              //< Синхронизация горутин при доступе к объектам структуры
		urls           []string //< Срез подходящих ссылок
	}{
		urls: make([]string, 0, len(arrLink)),
	}

	for _, item := range arrLink {
		// u := item
		url := item //< Сохраням значение переменной, т.к. u будет меняться
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

			// Если тело содержит строку удовлетворяющую шаблону - добавляем в массив ответа
			matched, err := regexp.MatchString(pattern, string(body))
			if err != nil {
				return err
			}
			if matched {
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
