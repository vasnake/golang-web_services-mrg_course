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
Можно указать явно, `int8, int16, int32, int64`.
Аналогично `uint`.

`float32, float64` нет просто `float`.

Есть `complex64, complex128` математики и физики радуются.

Строки в кавычках интерпретируются, символы типа `\n` и прочие будут транслированы.

Строки в бэктиках не интерпретируются.

Одинарные кавычки для задания byte (uint8) или rune (uint32).

Строки immutable.

Длина строки считается в байтах. Для подсчета в символах используй `utf8.RuneCountInString(someStr)`.
Соответственно, срезы тоже в байтах.
Строки можно легко конвертировать в байты и байты в строки.

Константы `const name = value`.
Блоки констант.
Опредение через `iota`, автоинкремент, всё сложно.
Нетипизированные константы, тип присваивается при записи константы в переменную. Вроде макроса получается.

Пользовательские типы данных, `type`. Полезно при моделировании, DSL.
Нет автоматического приведения типов.

Нет адресной арифметики, но есть ссылки, reference. Полезно для передачи структур без копирования.
- `b := &a` получение ссылки.
- `*b = 42` запись значения по ссылке.
- `c := new(int)` создание ссылки на безымянную переменную.

### Переменные, составные типы данных

- [array](week_01/array.go)
- [slice_1](week_01/slice_1.go)
- [slice_2](week_01/slice_2.go)
- [map](week_01/map.go)

`var arr3 [3]int` размерность массива входит в определение типа переменной, массивы разных размерностей не совместимы по типам.

Для определения размера массива можно использовать константы, но нельзя переменные.

`arr3 := [...]int{1, 2, 3}` Размер массива можно не задавать при явной инициализации.

При выходе за границы массива в рантайме получаем панику.

По причине несовместимости типоразмеров массивов, на практике работа с массивами вынесена на нижний уровень.
Сверху, то чем пользуются прикладники -- срезы, slice.

- `var buf []int` создание пустого слайса без инициализации
- `buf := []int{} // len:0, cap:0` создание пустого с инициализацией
- `buf := make([]int, 5, 10) // len:5, cap:10` срез это некий буфер, у него есть длина и емкость.

- `buf = append(buf, 9, 10) // len:2, cap:2` буфер может расти, при исчерпании емкости буфер пересоздается с удвоенной емкостью.
- `buf = append(buf, otherBuf...)` при добавлении другого буфера, его надо "распаковать".

Слайсы могут работать "по ссылке", оперируюя значениями в одном и том-же буфере.
То есть, если явно не выделять память под слайс (или неявно, через append), то работа идёт в одном и том-же буфере.

- `numCopied = copy(emptyBuf, buf)` копирование элементов в другой буфер, внутри проверка на выход за границы.
- `copy(buf[1:3], []int{5, 6})` копирование под-диапазона.

Мапки `var user map[string]string`, можно, как и слайс, создать с нужной ёмкостью, через `make`
- `mName, mNameExists := user["middleName"]` правильный способ получения значения из мапки, ибо по умолчанию, несуществующее значение = пустая строка.
- `delete(user, "lastName")` удаление ключа

### Управляющие конструкции

- [control](week_01/control.go)
- [loop](week_01/loop.go)

- `if boolVal { ... }` только тип bool.
- `if v, exists := myMap["name"]; exists { ... }` условие с блоком инициализации
- `switch len(myMap) { 	case 0, 1: ... }` по умолчанию делает break при срабатывании условия
- `switch ... case k == "name" && v == "Bender": ...` сложные условия в switch
- `switch ... break` оператор выхода, можно выходить через несколько уровней, по метке

Циклы определяются ключевым словом `for`, есть разные формы.

`for pos, symb := range myStr { ... }` итерирование строки делается по символам, не байтам.

### Основы функций

- [functions](week_01/functions.go)

`func sqrt(x int) int { ... }` и несколько более сложных вариантов объявления.
Например, именованный возвращаемый результат, иногда бывает удобно.

`func namedWithError(condition bool) (res int, err error) { ... }` осторожнее с значениями "по умолчанию".

`func sum(in ...int) int { ... }` кортежи параметров и кортежи возвращаемых значений -- это нормально.

Переменное количество входных параметров базируется на представлении параметров как слайса.

### Функция как объект первого класса, анонимные функции

- [firstclass](week_01/firstclass.go)

Функция как значение переменной -- присваивать, передавать, возвращать.

`printer := func(msg string) { ... }` анонимная функция как значение переменной.

### Отложенное выполнение и обработка паники

- [defer](week_01/defer.go)
- [recover](week_01/recover.go)

- `defer doStuff("after work ...")` будет вызвана перед выходом из области видимости. Полезно как код финализации процедур.
- Несколько defer складываются в стек (FILO).
- Аргументы отложенных функций вычисляются НЕ отложенно а сразу. Чтобы этого избежать, такие аргументы заворачиваются в анонимную функцию.

defer полезен при восстановлении из паники.
Если внутри defer вызвать `recover()`, то программа не вывалится в панику а продолжит работать штатно.

### Основы работы со структурами

- [structs](week_01/structs.go)

`type Person struct { Id int ... }` поля могут быть любых типов.

Полный формат инициализации структуры `var acc Account = Account{Id: 42, ... }`
с использованием имён полей. При пропуске поля -- его значение будет "по умолчанию".

Краткая форма инициализации `acc.Owner = Person{33, "Foo Bar", "Under The Bridge"}`
без имён полей, но пропускать поля уже нельзя.

Композиция структур: встраивание полей одной структуры в namespace другой структуры (подробнее см. код).
При конфликте имён, выигрывает поле вышележащее в иерархии вложенности.

### Методы структур

- [methods](week_01/methods.go)

Метод: функция, привязанная к типу данных.
Особенность языка: нет необходимости определять методы при определении типа. Можно привязать метод к типу когда угодно.

`func (p *Person) SetName(newName string) { p.Name = newName }` N.B. структура должна быть передана by-reference в методах-сеттерах.
Иначе метод будет применяться к копии объекта.
Язык не требует явного указания передачи by-reference при вызове метода, достаточно того, что by-reference указан в определении метода.

При композиции структур, внешняя структура получает все методы встроенных структур.

### Пакеты и область видимости

- [dir.txt](week_01/dir.txt)
- [visibility/main](week_01/visibility/main.go)
- [visibility/person/person](week_01/visibility/person/person.go)
- [visibility/person/func](week_01/visibility/person/func.go)

```s
# visibility
|---person
|   |---person.go
|   |---func.go
|---main.go

```

> As of Go v1.13, by default, go modules are used.
> Therefore, you need to tell explicitly if you don't want to do this. `GO111MODULE=off go run main.go`

GOPATH определяет корневую директорию, в которой будут под-директории `bin, pkg, src`.

Имя пакета это имя директории.
Приватные поля определяются именованием с маленькой буквы, публичные поля -- с большой.

Доступ к приватным полям возможен только в коде пакета, где определено приватное поле.

Крупные пакеты, с большим количеством файлов, предпочтительнее мелких пакетов, с малым количеством файлов.

Зависимости в директории vendor.

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
