Укорачиватель ссылок
=============================
Предоставляет API по созданию сокращённых ссылок следующего формата:
- Ссылка уникальна и на один оригинальный URL должна ссылаться только одна сокращенная ссылка.
- Ссылка длинной 10 символов и состоит из символов латинского алфавита в нижнем и верхнем регистре, цифр и символа _ (подчеркивание)

Сервис принимает следующие запросы по gRPC:
1. Метод Generate - сохраняет оригинальный URL в базе и возвращает сокращённый

2. Метод Retrive- принимаtn сокращённый URL и возвращаtn оригинальный URL

Сервер gRPC запускается на порту 8020

Аргументы клиентской части, передаются при запуске приложения:
 -link="передаваемая ссылка"
 -g передает метод Generate
 -r передает метод Retrive

Например: go run cmd/client/main.go -link=https://www.google.com -g

Результат выполнения: Generated short link: http://localhost:8080/to/hpUPxmAgIL

Сервер HTTP запускается на порту 8080

Позволяет сохранять длинную ссылку и возвращать короткую. При переходе по короткой ссылке, происходит переадресация на сохраненную длинную ссылку (справедливо и для ссылок полученных через gRPC)
