# ДЗ
1. Используя функцию для поиска из прошлого практического задания, постройте сервер, который будет принимать JSON с поисковым запросом в POST-запросе и возвращать ответ в виде массива строк в JSON.
```JSON
{
  "search":"фраза для поиска",
  "sites": [
      "первый сайт",
      "второй сайт"
  ]
}
```
2. Напишите два роута: один будет записывать информацию в Cookie (например, имя), а второй — получать ее и выводить в ответе на запрос.

### Немного саморазвития и ответ на вопрос "Что почитать?"
	  
#### Советую ознакомиться с этими библиотеками для Go:
1) https://github.com/go-chi/chi (Продвинутые роуты)
2) https://github.com/Sirupsen/logrus (Продвинутый логгер)

#### Советую почитать про архитектурный подход REST и RESTfull:
1) https://habr.com/ru/post/38730/ (Кратко о REST)
2) https://habr.com/ru/post/351890/ (Best practices)
3) https://habr.com/ru/post/265845/ (Немного холивара о RESTfull)

#### Для совсем хардкорщиков, гляньте в сторону gRPC (библиотека от гугла):
1) https://grpc.io/ (Сайт)
2) https://github.com/grpc/grpc-go (Библиотека для Go)


<p>

[postserver](packages/postserver/postserver.go), [postserver_test](packages/postserver/postserver_test.go) - Using function to search from the past practice, build a server that will accept JSON with search request in the POST request and return the response as an array of strings in JSON.
</p>
<p>

[cookieserver](packages/cookieserver/cookieserver.go), [cookieserver_test](packages/cookieserver/cookieserver_test.go) - Write two routes: one will record information in a cookie (for example, a name), and the second will receive it and display it in response to a request.
</p>