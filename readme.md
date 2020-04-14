# Что это?

Здесь представлена краткая инструкция по настройке gRPC
и некоторой его экосистемы.

[Ссылка на официальную документацию по настройке](https://grpc.io/docs/quickstart/go/)

# Установка gRPC

Перед выполнением всех действий убедитесь, что у вас настроены переменные окружения:
- $GOPATH - Указывает на вашу домашнюю директорию Go
- $PATH - Должен содержать в себе путь до `$GOPATH/bin`

Сперва нужно установить gRPC
```
go get -u google.golang.org/grpc
```

Скачать **protoc** компилятор для вашей платформы.
https://github.com/google/protobuf/releases

Быстрые ссылки на архивы:
- [Protoc 3.11.4 Linux x64](https://github.com/protocolbuffers/protobuf/releases/download/v3.11.4/protoc-3.11.4-linux-x86_64.zip)
- [Protoc 3.11.4 Windows x64](https://github.com/protocolbuffers/protobuf/releases/download/v3.11.4/protoc-3.11.4-win64.zip)
- [Protoc 3.11.4 OSX x64](https://github.com/protocolbuffers/protobuf/releases/download/v3.11.4/protoc-3.11.4-osx-x86_64.zip)

В архиве будет файл `bin/protoc`, который следует закинуть
куда-нибудь, где у вас настроен **PATH**.
Например, вы можете положить его в `~/go/bin/`

Другую директорию `include` следует закинуть в `/usr/local/`.
Так же я видел варианты, когда эту директорию хранили прямо в проекте.

Потом нужно поставить плагин для protoc, который позволит
компилировать proto файл в код GoLang:
```
go get -u github.com/golang/protobuf/protoc-gen-go
```

# Проверка корректности установки

Прежде чем настраивать прочую инфраструктуру, предлагаю проверить,
что мы всё установили правильно и всё работает.

Зайдите в директорию
```
$GOPATH/src/google.golang.org/grpc/examples/helloworld
```
Она должна появиться вместе с тем, как вы скачали google.golang.org/grpc.
Если её нет, проверьте правильность настройки переменной **GOPATH**
и выполните
```
go get google.golang.org/grpc
```

Далее, находясь в этом проекте запустите два сервиса (при помощи
двух разных терминалов, или tmux/screen):
```
go run greeter_server/main.go
``` 
```
go run greeter_client/main.go
``` 

Если всё ок, то на клиенте вы увидите сообщение "Greeting: Hello world"

Теперь попробуем проверить корректность работы утилиты **protoc**,
для этого в той же директории с примером выполните:
```
protoc -I helloworld/ helloworld/helloworld.proto --go_out=plugins=grpc:helloworld
```

Если всё прошло без ошибок - gRPC полностью рабочий и настроенный!

# Настройка grpc-gateway и grpc-swagger

Это HTTP шлюз, который позволит взаимодействовать с вашим gRPC
сервисом другим сервисам, которые не могут в gRPC.
Вместо этого они могут использовать альтернативный HTTP шлюз.

[Документация по Gateway и Swagger](https://github.com/grpc-ecosystem/grpc-gateway)

Для установки этих инструментов выполните:
```
go install \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
    github.com/golang/protobuf/protoc-gen-go
```

Теперь в ваш proto файл можно добавить специальные аннотации.
Например, он может выглядеть как [файл в этом проекте](proto/communication.proto).

Такие аннотации помогают описать структуру HTTP запросов.

Выполните для получения дополнительных proto файлов расширения:
```
go get github.com/grpc-ecosystem/grpc-gateway
```

После этого должна появиться директория **$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway**  
Если этого не произошло - склонируйте её вручную туда.

Сам файл аннотаций и прочие расширения [можно взять здесь](https://github.com/googleapis/googleapis/blob/master/google/api/annotations.proto).

Далее нужно перегенерировать gRPC код с учётом плагина:
```
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=plugins=grpc:. \
  proto/communication.proto
 */
```

Теперь нужно сгенерировать код для HTTP Proxy:
```
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=logtostderr=true:. \
  proto/communication.proto
```

И генерация сваггера:
```
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --swagger_out=logtostderr=true,repeated_path_param_separator=ssv:. \
  proto/communication.proto
```

Всё описанное выше можно скастовать в одну команду:
```
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=plugins=grpc:. \
  --grpc-gateway_out=logtostderr=true:. \
  --swagger_out=logtostderr=true,repeated_path_param_separator=ssv:. \
  proto/communication.proto
```

Пример swagger аннотаций для proto файла [можно посмотреть здесь](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/examples/internal/proto/examplepb/a_bit_of_everything.proto).
Это просто большой и полноценный proto файл с кучей примеров.