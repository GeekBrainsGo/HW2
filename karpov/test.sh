#!/bin/zsh
curl -X POST 127.0.0.1:3333/search -H 'Content-type: application/json' \
-d '{"search":"Чак Норрис", "sites":["https://ru.wikipedia.org/wiki/Норрис,_Чак","https://ru.wikipedia.org/wiki/Крутой_Уокер:_Правосудие_по-техасски","https://yandex.ru"]}'
