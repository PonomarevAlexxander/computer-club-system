# Тестовое задание Go
______
## Формулировка задания
Требуется написать прототип системы, которая следит за работой компьютерного клуба, обрабатывает события и подсчитывает выручку за день и время занятости каждого стола.
Входные данные представляют собой текстовый файл. Файл указывается первым аргументом при запуске программы.
Программа должна запускатьcя в Linux или Windows с использованием docker container-a (требуется написание Dockerfile). Требуется использование стандартной библиотеки (https://pkg.go.dev/std).

## Описание решения
Решение представляет из себя программу ня `Go` с возможностью сборки в `Dockerfile`.

В репозитории представлены тестовые примеры для запуска. Они могут быть найдены в директории [examples](/examples/).

Основные шаги для сборки и запуска программы локально можно найти в [Makefile'e](Makefile):
- Сборка приложения осуществляется с помощью
```shell
make build
```
- Сборка `Docker-image` может быть запущена с помощью
```shell
make docker-build
```
- Запуск `Docker-image` с использованием тестовых сценариев из `/examples`
```shell
make docker-run file=<название файла из examples>
```
Стоит отметить, что Вы можете запустить контейнер сами, если примонтируете volume (`-v`),
в котором будет находиться текстовый файл для запуска и передадите путь в параметры `docker run`.

Например, `docker run -v "/путь/к/file.txt:/путь/к/file.txt" computer-club-system "/путь/к/file.txt"`
- Запуск тестовых сценариев из `/examples` может быть запущен вручную или через
```shell
make test
```
- Запуск unit тестов
```shell
make unit-test
```
