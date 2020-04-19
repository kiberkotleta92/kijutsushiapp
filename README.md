# kijutsushiapp

kijutsushi (きじゅつし) - conjurer, magician, juggler, illusionist

## Задача

В данной задаче вам предстоит разработать программу, которая способна получать ссылку на видео файл и преобразовать речь из него в текст.

Программа должна работать следующим образом:

- На вход подается путь до директории с видео (не дольше 30 секунд).
- Получившиеся аудиофайлы отправляются на яндекс speachKit и в ответ возвращается текст.
- Из полученных фрагментов текста формируется текстовый файл.

### Итоговый результат

Итогом вашей работы должна стать программа, которая принимает путь до директории (через интерфейс или через консоль) и оцифровывает данные в текстовый файл.

## Решение

На виртуальной машине работает приложение-сервер, которое получает видео файл. Сервер, используя [ffmpeg](https://www.ffmpeg.org), вытаскивает из видео аудио и дробит на файлы по 30 секунд. Затем сервер отправляет эти файлы в [Yandex SpeechKit](https://cloud.yandex.ru/services/speechkit), обьединяет полученные фрагменты текста и возвращает отправителю видео.


Пользователь запускает приложение командной строки, которое принимает видео файл, отправляет приложению-серверу, и записывает ответ в стандартный вывод.


## Установка
### Установка сервера
```
go get github.com/kirilldenisov/kijutsushiapp/kijutsushiserver
scp trueenv.go login@adress:GOPATH/github.com/kirilldenisov/kijutsushiapp/kijutsushiserver/env.go
go install GOPATH/github.com/kirilldenisov/kijutsushiapp/kijutsushiserver

# запуск
kijutsushiserver
```

### Установка приложения коммандной строки
```
go get github.com/kirilldenisov/kijutsushiapp/kijutsushi
go install GOPATH/github.com/kirilldenisov/kijutsushiapp/kijutsushi

# запуск
kijutsushi <path_to_video>
```