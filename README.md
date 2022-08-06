# go-web_services-mrg_course

Golang course from MRG: [Разработка веб-сервисов на Golang](https://study.vk.team/learning/play/671)

Блоки программы. Часть 1:
- Введение в Golang
- Асинхронная работа
- Работа с динамическими данными и производительность
- Основы HTTP

## part 1, week 1

Введение в Go.
[Код, домашки, литература](week_01/materials.zip) https://cloud.mail.ru/public/pc7H/6Vx4txWWr

Зачем нужен еще один языка программирования?
Эффективность (работы стажёров в Гугл): компиляции, выполнения, разработки.
Утилизация многопроцессорных систем (легкие потоки, асинхронность);
Простой и понятный язык, читабельный, простая и быстрая сборка, с четким стилем.
Сборка в статический бинарник, решение проблемы dll-hell.

Для realtime, embedding не подходит (в наличии сборка мусора), но всё остальное OK.

### Первая программа

- https://play.golang.com/
- [run](run.sh)
- [hello_world](week_01/hello_world.go)

`package main`: основной (запускаемый) пакет программы.

`func main() {...}`: основной (запускаемый) функ. программы. Если такой нет, то запускатор будет ругаться
```s
go run hello_world.go
runtime.main_main·f: function main is undeclared in the main package
```

Как удобно в Scala, всё есть expression. Но в Golang не так, увы.

Используй `camelCase` для имён. Публичные обьекты называй с большой буквы `Println`.

### Переменные, базовые типы данных

- [vars_1](week_01/vars_1.go)
- [vars_2](week_01/vars_2.go)
- [strings](week_01/strings.go)
- [const](week_01/const.go)
- [types](week_01/types.go)
- [pointers](week_01/pointers.go)

`var name type` создает переменную со значением "по умолчанию";
считается достижением, что не бывает неопределенных значений, всегда память (выделенная под переменную) инициализируется.

`var name = 42` создает переменную с типом, выведенным автоматом из данного значения.

`name := 42` короткая форма создания новой переменной.

`name1, name2 = 42, 37` множественное присваивание (инициализация) работает.

[Почему нет префиксного инкремента?](https://go.dev/doc/faq#inc_dec)
> Why are ++ and -- statements and not expressions? And why postfix, not prefix?

Without pointer arithmetic, the convenience value of pre- and postfix increment operators drops.
By removing them from the expression hierarchy altogether, expression syntax is simplified and
the messy issues around order of evaluation of ++ and -- (consider f(i++) and p[i] = q[++i]) are eliminated as well.
The simplification is significant

`var i int` разрядность зависит от платформы.
Можно указать явно, `int8, int16 ... int64`.
Аналогично `uint`.

`float32, float64` нет просто `float`.

Есть `complex64, complex128` математики и физики радуются.

Строки в кавычках интерпретируются, `\n` и прочие будут транслированы.

Строки в бэктиках не интерпретируются.

Одинарные кавычки для задания byte (uint8) или rune (uint32).

Строки immutable.

Длина строки считается в байтах. Для подсчета в символах `utf8.RuneCountInString(someStr)`.
Соответственно, срезы тоже в байтах.
Строки можно легко конвертировать в байты и байты в строки.

Константы `const name = value`.
Блоки констант.
Опредение через `iota`, автоинкремент, всё сложно.
Нетипизированные константы.

Пользовательские типы данных, `type`. Полезно при моделировании, DSL.
Нет автоматического приведения типов.

Нет адресной арифметики, но есть ссылки, reference.
- `b := &a` получение ссылки.
- `*b = 42` запись значения по ссылке.
- `c := new(int)` создание ссылки на безымянную переменную.

### Переменные, составные типы данных: array, slice_1, slice_2, map
### Управляющие конструкции: control, loop
### Основы функций: functions
### Функция как объект первого класса, анонимные функции: firstclass
### Отложенное выполнение и обработка паники: defer, recover
### Основы работы со структурами: structs
### Методы структур: methods
### Пакеты и область видимости: dir.txt, visibility/
### Основы работы с интерфейсами: basic, many, cast
### Пустой интерфейс: empty_1, empty_2
### Композиция интерфейсов: embed
### Написание Программы Уникализации (ПУ): uniq, data_map.txt
### Написание тестов для ПУ: unique/unique, unique/unique_test

## part 1, week 2

Асинхронная работа.
[Код, домашки, литература](week_02/w2_materials.zip) https://cloud.mail.ru/public/YDEX/Dau2wVWuw/

- Методы обработки запросов и плюсы неблокирующего подхода:
  асинхронное выполнение, скорость процессор-кеш-память, время на переключение контекста, современные тенденции на многоядерность и параллельность,
  тяжелые процессы, потоки легче, асинхронные сопрограммы (green threads) еще легче. Невытесняющая многозадачность (eventloop, Windows 3.0) vs preemptive.
  Ввод-вывод и ожидание возврата из syscall. Время ожидания можно потратить на другие задачи, non-blocking IO. IO-bound vs CPU-bound.
  `Communicating Sequential Processes` by Tony Hoare. Горутины перемещаются между системными потоками.
- Горутины -- легковесные процессы: goroutines, `go` keyword.
- Каналы -- передаём данные между горутинами: chan_1, chan_2; `chan` keyword. Передача контроля над данными между потоками/горутинами.
  Чтение, запись -- операторы стрелочка `x <- myChannel; myChannel <- x`.
  Небуферизованные каналы vs буферизованные, размер буфера. Работа с каналом в цикле.
- Мультиплексирование каналов через (не блокирующий) оператор `select`: select_1, select_2, select_3. Выбор в цикле, канал "отмены".
- Таймеры и таймауты (как источник сигнала в каналах): timeout, tick, afterfunc. Тикер как источник регулярных/периодических сигналов в канале.
  AfterFunc как способ отложенного выполнения функции.

### TODO
- Пакет `context` и отмена выполнения
- Асинхронное получение данных
- Пул воркеров
- `sync.Waitgroup` -- ожидание завершения работы
- Ограничение по ресурсам
- Ситуация гонки на примере конкурентной записи в map
- `sync.Mutex` для синхронизации данных
- `sync.Atomic`

## week 3
Работа с динамическими данными и производительность

## week 4
Основы HTTP

## Info, links

- [If a map isn’t a reference variable, what is it?](https://dave.cheney.net/2017/04/30/if-a-map-isnt-a-reference-variable-what-is-it)
