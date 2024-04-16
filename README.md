# Разработка веб-сервисов на Go

[Разработка веб-сервисов на Golang (Go), Василий Романов, stepik](https://stepik.org/187490)
Этот курс был создан в 2017 году на основе внедрения языка Go в Почту Mail.ru

3 части, 12 недель
- [part1](part1.md)
- [part2](part2.md)
- [part3](part3.md)

## silly-bus

https://stepik.org/course/187490/syllabus

### Неделя 1 - основы языка, Часть 1
- 1.1 Правила, чат поддержки, код к лекциям и домашние задания
- 1.2 Начало работы 
- 1.3 Основы языка 
- 1.4 Функции 
- 1.5 Структуры и методы 
- 1.6 Интерфейсы 
- 1.7 Практический пример - программа уникализации с тестами 
- 1.8 Задание 1 - программа вывода дерева файлов

Неделя 2 - асинхронная работа, Часть 1 
- 2.1 Методы обработки запросов 
- 2.2 Горутины и каналы 
- 2.3 Инструменты для многопроцессорного программирование 
- 2.4 Состояние гонки 
- 2.5 Задание 2 - асинхроннй пайплайн

Неделя 3 - json и бенчмарки, Часть 1 
- 3.1 JSON 
- 3.2 Работа с динамическими данными 
- 3.3 Бенчмарки и производительность 
- 3.4 Задание 3 - оптимизация кода

Неделя 4 - основы работы с HTTP, Часть 1 
- 4.1 Слушаем сетевое соединение 
- 4.2 Обработка HTTP-запросов 
- 4.3 Шаблонизация 
- 4.4 Профилирование веба 
- 4.5 Телеграм бот 
- 4.6 Задание 4 - тестовое покрытие для сервиса поиска по XML

### Неделя 5 - продолжаем работу с HTTP, Часть 2
- 5.1 Приветствие 
- 5.2 Middleware 
- 5.3 Роутинг http-запросов 
- 5.4 Валидация входящих данных 
- 5.5 Фреймворки 
- 5.6 Логирование 
- 5.7 Веб-сокеты 
- 5.8 Шаблонизация 
- 5.9 Задание 5 - веб-фреймворк на основе кодогенерации

Неделя 6 - базы данных, Часть 2 
- 6.1 SQL 
- 6.2 KV-хранилища 
- 6.3 Rabbit, Mongodb 
- 6.4 Задание 6 - универсальный сервис просмотра содержимого БД

Неделя 7 - основы микросервисов, Часть 2 
- 7.1 Что такое микросервис 
- 7.2 Делаем микросервис руками 
- 7.3 protobuf и gRPC 
- 7.4 Дополнительные темы 
- 7.5 Задание 7 - асинхронная система логирования

Неделя 8 - прочие темы, Часть 2 
- 8.1 Конфигурирование сервиса 
- 8.2 Мониторинг 
- 8.3 Низкоуровневое программирование 
- 8.4 Инструменты для статического анализа 
- 8.5 Задание 8 - заполнение полей структуры через рефлексию

### Неделя 9 - архитектура приложения, Часть 3
- 9.1 Структурируем приложение 
- 9.2 Тестируем комплексное приложение 
- 9.3 Авторизация и пароли 
- 9.4 CSRF-токены 
- 9.5 Сессии 
- 9.6 Задание 9 - архитектура типового приложения

Неделя 10 - oauth и рефакториг приложения, Часть 3 
- 10.1 OAuth 
- 10.2 Немного рефакторинга 
- 10.3 Проектирование API 
- 10.4 Задание 10 - телеграм бот

Неделя 11 - graphql, Часть 3 
- 11.1 Основы GraphQL 
- 11.2 GraphQL - интеграция в проект 
- 11.3 Организация пакетов в приложении 
- 11.4 Задание 11 - маркетплейс на основе GraphQL

Неделя 12 - сборка, s3 и трейсинг, Часть 3 
- 12.1 Сборка docker-контейнера 
- 12.2 Хранение файлов в проекте через S3 
- 12.3 Конфигурирование приложения 
- 12.4 Трейсинг запросов 
- 12.5 Обратная связь 
- 12.6 Задание 12 - многопользовательская MUD на основе асинхрона

## source code

Код ко всем лекциями можно скачать по следующей ссылке:
https://stepik.org/media/attachments/lesson/1177827/golang_web_services_2023-12-28.zip
- Данный код проверен на версии языка go 1.20
- В папке каждой недели есть домашка, она находится в подпапке 99_hw
- При работе с консолью пользуйтесь Bash

Я уже скачал: [handouts\golang_web_services_2023-12-28.zip](./handouts\golang_web_services_2023-12-28.zip)

У меня винда и WSL, поэтому я себе сделал такой раннер для кода:
```hs
pushd /mnt/c/Users/${LOGNAME}/data/github/golang-web_services-mrg_course/sandbox
alias gr='bash -vxe /mnt/c/Users/${LOGNAME}/data/github/golang-web_services-mrg_course/run.sh'
# всё запускаемое можно найти в скрипте
```
bash

## extra homework

### float numbers

> Numeric constants represent exact values of arbitrary precision and do not overflow.
Consequently, there are no constants denoting the IEEE-754 negative zero, infinity, and not-a-number values.

Напиши библиотечные функции:
- Предикат `IsInvalidFloat`, используемый в фильтрации коллекции чисел. Инвалид: inf, nan
- Сравнения двух чисел float, с учетом дельты, +/- zero, nan, +/- inf

### int numbers bits

> Two's complement is achieved by:
- Step 1: starting with the equivalent positive number.
- Step 2: inverting (or flipping) all bits – changing every 0 to 1, and every 1 to 0;
- Step 3: adding 1 to the entire inverted number, ignoring any overflow.

Напиши отображение 32-бит числа в int32, uint32. Битовое представление не меняется, меняется интерпретация набора бит.
Рассмотри два случая: int32 получен как результат хеширования
и 1) надо отобразить диапазон значений на uint32. 2) надо сохранить битовое представление.

Напиши конвертеры: (u)int => bitstring, bitstring => u(int).

### concurrency patterns

Напиши параллельную обработку по трем паттернам https://habr.com/ru/companies/timeweb/articles/770912/

### week 2, dynamic workers pool

[week 2, extra homework](./sandbox/week02_homework/wp_extra/wp_extra.md)

### week 4, search in XML

[week 4, extra homework](./part1.md#week4-homework)
Часть задания, в которой надо реализовать поиск "юзеров" в файле XML, согласно переданным параметрам.

## Info, links

- `Communicating Sequential Processes` by Tony Hoare
- https://cs.stanford.edu/people/eroberts/courses/soco/projects/2008-09/tony-hoare/csp.html
- http://www.usingcsp.com/cspbook.pdf
- CSP vs Actor model (channels vs actor/mailbox)
- [If a map isn’t a reference variable, what is it?](https://dave.cheney.net/2017/04/30/if-a-map-isnt-a-reference-variable-what-is-it)
- [Implementing a bignum calculator with Rob Pike](https://youtu.be/PXoG0WX0r_E)
- [Lexical Scanning in Go - Rob Pike](https://www.youtube.com/watch?v=HxaD_trXwRE)
- [Advanced Topics in Programming Languages: Concurrency/message passing Newsqueak](https://youtu.be/hB05UFqOtFA), prime sieve on channels
- [3.5 Years, 500k Lines of Go (Part 1)][https://npf.io/2017/03/3.5yrs-500k-lines-of-go/] `return fmt.Errorf("While doing foo: %v", err.Error())`
- The Go Playground https://go.dev/play/ or https://play.golang.com/
- Practical Go: Real world advice for writing maintainable Go programs / Dave Cheney https://www.google.com/search?q=Practical+Go%3A+Real+world+advice+for+writing+maintainable+Go+programs+%2F+Dave+Cheney
