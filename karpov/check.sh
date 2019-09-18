#!/bin/zsh
# success
curl -X POST 127.0.0.1:3333/v1/search -H 'Content-type: application/json' \
-d '{"search":"Чак Норрис", "sites":["https://ru.wikipedia.org/wiki/Норрис,_Чак","https://ru.wikipedia.org/wiki/Крутой_Уокер:_Правосудие_по-техасски","https://yandex.ru"]}'

# not found
curl -X POST 127.0.0.1:3333/v1/search -H 'Content-type: application/json' \
-d '{"search":"Чак Норрис2", "sites":["https://ru.wikipedia.org/wiki/Норрис,_Чак","https://ru.wikipedia.org/wiki/Крутой_Уокер:_Правосудие_по-техасски","https://yandex.ru"]}'

# wrong header
curl -X POST 127.0.0.1:3333/v1/search -H 'Content-type: application/xml' \
-d '{"search":"Чак Норрис", "sites":["https://ru.wikipedia.org/wiki/Норрис,_Чак","https://ru.wikipedia.org/wiki/Крутой_Уокер:_Правосудие_по-техасски","https://yandex.ru"]}'

# wrong payload
curl -X POST 127.0.0.1:3333/v1/search -H 'Content-type: application/json' \
-d '"Чак Норрис", "https://ru.wikipedia.org/wiki/Норрис,_Чак","https://ru.wikipedia.org/wiki/Крутой_Уокер:_Правосудие_по-техасски","https://yandex.ru"'

## COOKIES
# set COOKIES
curl -v -c ./cookie-jar.txt 'http://127.0.0.1:3333/v1/cookies?username=vasya'
# get COOKIES
curl -v -b ./cookie-jar.txt http://127.0.0.1:3333/v1/cookies/data
