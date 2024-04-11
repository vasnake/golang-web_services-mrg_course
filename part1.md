# Разработка веб-сервисов на Golang (Go). Часть 1

[Go course, MRG, Романов Василий](README.md)
[Разработка веб-сервисов на Golang (Go), Василий Романов, stepik](https://stepik.org/187490)

Неделя 1 - основы языка
- 1.1 Правила, чат поддержки, код к лекциям и домашние задания
- 1.2 Начало работы 
- 1.3 Основы языка 
- 1.4 Функции 
- 1.5 Структуры и методы 
- 1.6 Интерфейсы 
- 1.7 Практический пример - программа уникализации с тестами 
- 1.8 Задание 1 - программа вывода дерева файлов

Неделя 2 - асинхронная работа
- 2.1 Методы обработки запросов 
- 2.2 Горутины и каналы 
- 2.3 Инструменты для многопроцессорного программирование 
- 2.4 Состояние гонки 
- 2.5 Задание 2 - асинхроннй пайплайн

Неделя 3 - json и бенчмарки
- 3.1 JSON 
- 3.2 Работа с динамическими данными 
- 3.3 Бенчмарки и производительность 
- 3.4 Задание 3 - оптимизация кода

Неделя 4 - основы работы с HTTP
- 4.1 Слушаем сетевое соединение 
- 4.2 Обработка HTTP-запросов 
- 4.3 Шаблонизация 
- 4.4 Профилирование веба 
- 4.5 Телеграм бот 
- 4.6 Задание 4 - тестовое покрытие для сервиса поиска по XML

## part 1, week 1

Основы языка, введение в Go. [Код, домашки, литература](./handouts\golang_web_services_2023-12-28.zip)

Зачем нужен еще один язык программирования?

Go-team (backend development) хотела язык C/C++ но без их недостатков, плюс эффективная утилизация многопроцессорных систем.
Ключевые части проекта Go: зависимости, рантайм, garbage-collection, модель конкурентности (асинхронности).
Эффективность (работы программеров в Гугл): компиляции, выполнения, разработки.
Утилизация многопроцессорных систем через простой интерфейс: каналы и горутины (легкие потоки, асинхронность, CSP).
Простой и понятный язык, читабельный, с четким стилем.
Простая и быстрая сборка. Сборка в статический бинарник, решение проблем dll-hell.

Go не подходит для realtime, embedding (в наличии сборка мусора), но всё остальное OK.

Смысл всего упомянутого будет кристально ясен после просмотра серии YT-видео, объясняющих почему, откуда и зачем Go.
Ищи видео от Go-Team и конкретно Rob Pike.

Установка и наладка development environment:
- https://go.dev/doc/install
- https://go.googlesource.com/vscode-go/+/refs/heads/release.theia/README.md
- https://code.visualstudio.com/docs/languages/go
- https://github.com/golang/vscode-go/blob/master/docs/tools.md
- https://learn.microsoft.com/en-us/azure/developer/go/configure-visual-studio-code

Homework
```s
# run `go env`, explain output
GOENV - файл с определениями переменных окружения
GOMOD - путь к файлу `go.mod`, руками лучше не трогать
```

### Первая программа

Борьба с vscode. Выводы.

Чтобы vscode заставить работать с golps и модулями, надо понимать:
- Что такое и как работать с: `vscode workspace`, `workspace root dirs`, `Multi-root Workspaces`.
- Аппа Go состоит из модулей, модули содержат пакеты, пакеты содержат файлы `.go`, [details](https://go.dev/doc/modules/layout)
- Минимальный набор файлов для проекта с одной аппой из одного модуля: `prj/go.work`, `prj/app/go.mod`, `prj/app/main.go`.
- Команды для генерации такого бойлерплейта
```s
# https://github.com/golang/tools/blob/master/gopls/doc/workspace.md#multiple-modules
# https://code.visualstudio.com/docs/editor/multi-root-workspaces
export app_dir=~/go_sandbox_project/hello_world && mkdir -p $app_dir && pushd $app_dir
touch main.go
go mod init hello_world
pushd ..
go work init
go work use ./hello_world/
```

> Since the two modules are in the same workspace it’s easy to make a change in one module and use it in another.
https://go.dev/doc/tutorial/workspaces
```s
$ mkdir workspace
$ cd workspace

$ mkdir hello
$ cd hello

$ go mod init example.com/hello

$ cat > hello.go << EOT
package main
func main() { }
EOT

$ go run .

$ cd .. # workspace

$ go work init ./hello # The `use` directive tells Go that the module in the `hello` directory should be main modules when doing a build
$ go run ./hello

$ # create example/hello module ...
$ go work use ./example/hello

$ # use example/hello/... package in hello/main code
$ go run ./hello
```
workspaces.

- [run](run.sh) `$ GO_APP_SELECTOR=hello_world gr`
- Код из лекции [w01/hello_world](week_01/hello_world.go)
- Моя переработка кода [snb/hello_world](./sandbox/hello_world/main.go)

В общем, мои грабли (настройки vscode+gopls) оказались таковы: в vscode у меня воркспейс был из одной рут-директории `desktop`,
создание проекта и модуля в `desktop/sandbox/hello_world` привело к ругани gopls.
Пришлось перенести сэндбокс в дерево вне `desktop` и подключить этот сэндбокс к vscode воркспейс как отдельный корень.
После чего gopls всосал `go.work` в корне сэндбокс и оттуда прочухал модуль `hello_world`.

vscode command (palette) `Go: Install/Update Tools` требуется для обновления инструментария (golint, gopls) после установки новой версии go.
Если это само не обновится, есть шанс получить в редакторе сообщения от устаревшей версии языка.

Я пока не знаю, можно ли в коде давать hint для линтеров, поэтому:
раздражающие предупреждения отключаются так, пример:
```s
# go vet:
go vet -stringintconv=false ./spec

#vscode go tools: ~\AppData\Roaming\Code\User\settings.json
"gopls": {
    "ui.diagnostic.analyses": {
        "stringintconv": false
    }
}
```

Используй `camelCase` для имён. Публичные (экспортируемые) обьекты называй с большой буквы: `Println`.

Файл с кодом программы содержит текст:
- `package main`: декларация пакета, в котором расположен код файла. Название `main` означает, что это основной пакет программы,
в котором расположена основная функция программы.
- `func main() {...}`: декларация основной функции пакета. Функция будет выполнена при старте программы.
```s
go run hello_world.go

# если пакета main нет:
package hello_world is not a main package

# если функции main нет:
runtime.main_main·f: function main is undeclared in the main package
```

I:
> Как удобно в Scala, всё есть expression. Но в Golang не так, увы.

### https://go.dev/ref/spec

Почитай спецификацию языка, осознай особенности.

```s
pushd sandbox
mkdir -p spec && pushd ./spec
touch main.go
go mod init spec
go mod tidy
popd
go work use ./spec/
go vet -stringintconv=false spec
gofmt -w spec
go run spec
```
`GO_APP_SELECTOR=spec gr` [spec playground](./sandbox/spec/main.go)

Отдельные моменты из спеки, вызвавшие особый интерес:

> general-purpose language designed with systems programming in mind.
It is strongly typed and garbage-collected and has explicit support for concurrent programming.

Пакет: имя, нэймспейс для группировки, средство изоляции: констант, типов, переменных, функций.
Пакет состоит из: один или несколько файлов, в одной директории (по имени пакета).
Приватные и публичные элементы пакета: идентификатор может быть "экспортирован" из пакета, если имя начинается с Большой Буквы.

Модуль: коллекция пакетов, сопровождаемая файлом `go.mod`. Средство дистрибуции кода, см."зависимости".

> Programs are constructed from packages, whose properties allow efficient management of dependencies.

> Go programs are constructed by linking together packages.
A package in turn is constructed from one or more source files that together declare
constants, types, variables and functions belonging to the package and which are accessible in all files of the same package.
Those elements may be exported and used in another package. ...
An implementation may require that all source files for a package inhabit the same directory

> A module is a collection of Go packages ... with a `go.mod` file ...
The `go.mod` file defines the module’s ... path..., and its dependency requirements ...

> Source code is Unicode text encoded in UTF-8
The text is not canonicalized, ... use the unqualified term character to refer to a Unicode code point in the source text.
... a compiler may disallow the NUL character (U+0000) in the source text
... A byte order mark ((U+FEFF)) may be disallowed anywhere else in the source
Unicode character categories: newline, letter, digit, char (not newline). `_` is considered a lowercase letter.

Comments:
> Line comments start with the character sequence `//` and stop at the end of the line.
General comments start with the character sequence `/*` and stop with the first subsequent character sequence `*/`.
A general comment containing no newlines acts like a space. Any other comment acts like a newline.

Tokens:
> four classes: identifiers, keywords, operators and punctuation, and literals.
... the next token is the longest sequence of characters that form a valid token
The formal syntax uses semicolons ";" as terminators ... Go programs may omit most of these semicolons ...

> Identifiers name program entities ... An identifier is a sequence of one or more letters and digits. The first character in an identifier must be a letter.
Some identifiers are predeclared.

> keywords are reserved and may not be used as identifiers (`map` is a keyword)

> Operators combine operands into expressions

> If a variable has not yet been assigned a value, its value is the zero value for its type.
Variables of `interface` type also have a distinct dynamic type, which is the (non-interface) type of the value assigned to the variable at run time ...

Вот тут (variable of interface type) возникает очень интересная проблема
> [Why is my nil error value not equal to nil?](https://go.dev/doc/faq#nil_error)
Вкратце: `error` это интерфейс, поэтому при проверке на ошибки получают interface value и сравнивают его с predefined `nil`.
Если ошибка (которой нет, nil) была получена из кода, где, по дороге, на нее навесили тип отличный от `error`, то
interface value будет `(SomeType, nil)` и это `!= nil`. Ибо `nil: (nil, nil)`.
[details](https://go.dev/play/p/CRZ_caKYCBR).

Строки: иммутабельная последовательность байт, которая интерпретируется как:
последовательность символов (рун), Unicode code-points. Длина строки в рунах != длине строки в байтах.
Если надо работать с коллекцией байт или рун, то надо явно привести тип (`[]rune(str)` or `[]byte(str)`).
Такое приведение типа создает копию данных, где иммутабельность уже не соблюдается.
ЧСХ: длина в байтах, доступ по индексу к байту, итерирование `range str` дает руны (кушайте, не обляпайтесь).
> A string value is a sequence of bytes. The number of bytes is called the length of the string ...
Strings are immutable
It is illegal to take the address of a string's byte

> Numeric constants represent exact values of arbitrary precision and do not overflow.
Consequently, there are no constants denoting the IEEE-754 negative zero, infinity, and not-a-number values.
Constants may be typed or untyped.
An untyped constant has a default type

> Numeric types: uint8..64, int8..64, float32..64, complex64..128, byte (alias uint8), rune (alias uint32).
Predeclared integer types with implementation-specific sizes: uint, int, uintptr.

Array: коллекция элементов одного типа; массивы разных размеров это разные типы.
Определение типа массива не может быть рекурсивным.
> The length is part of the array's type
An array type T may not have an element of type T, or of a type containing T as a component

Slice: можно сказать, что слайс это view на array элементов.
У слайса есть переменные длина и капасити.
Слайс share "хранилище" (массив) с другими слайсами, они все смотрят на один и тот-же кусок памяти.

Struct: набор элементов (полей) структуры.
Может содержать embedded поля, promoted поля. Нельзя определять поля рекурсивно.
> A field declared with a type but no explicit field name is called an embedded field

Function type: сигнатура функции определяет ее тип.
Возможны `variadic` фунции, 0 и более аргументов на последней позиции списка (кортежа) аргументов.
Возвращаемые параметры могут быть именованы, возвращаемые параметры это, в общем случае, кортеж.

Interface type: любой тип, имплементирующий данный интерфейс, входит в type set этого интерфейса и
может быть адресован через переменную типа этого интерфейса.
Следствие: любой тип имплементирует пустой интерфейс, переменная типа пустой интерфейс может ссылаться на любой объект (не-интерфейс).
> For convenience, the predeclared type `any` is an alias for the empty interface

Embedding interface E in T: an interface T may use a (possibly qualified) interface type name E as an interface element

Интерфейс можно ограничить, указав набор типов, для которых интерфейс валиден. Это будет не-basic интерфейс.
> Interfaces that are not basic may only be used as type constraints, or as elements of other interfaces used as constraints. They cannot be the types of values or variables, or components of other, non-interface types

Нельзя определять интерфейс рекурсивно.

Map type: не-упорядоченная коллекция элементов, с индексами-ключами.
> The comparison operators `==` and `!=` must be fully defined for operands of the key type
... A nil map is equivalent to an empty map except that no elements may be added
> Note that the `zero value` for a `slice` or `map` type is not the same as an `initialized but empty` value of the same type

Channael type: средство коммуникации (между горутинами), обмен сообщениями (элементами) определенного типа. FIFO queue.
Канал может быть буферизован, это не отражается на его типе.
Небуферизованный канал работает только если приемщик и отправитель оба готовы.
Т.е. операция с каналом может быть блокирующей, если операция не может быть поддержана буфером!
Канал на отправку можно (нужно) закрывать вызовом `close`, что маркирует его как "отправок более не будет".
> A channel may be constrained only to send or only to receive by assignment or explicit conversion.
... The `<-` operator associates with the leftmost chan possible

В целом, операции с каналами обладают сложным протоколом:
опциональная буферизация, закрытие канала, (не)блокирующие операции, паника и zero-значения ...
Простой язык, говорили они, WTF!

Type parameters:
> comparing operands of type parameter type may panic at run-time

Composite literals:
это не константы, каждый раз создается новое значение.
> The LiteralType's core type T must be a struct, array, slice, or map type

> Function literals are closures

Method expressions:
> For a method with a value receiver, one can derive a function with an explicit pointer receiver
... the method does not overwrite the value whose address is passed in the function call
... a value-receiver function for a pointer-receiver method, is illegal

Method values:
ресивер вычисляется и сохраняется (создается копия) при получении method value,
в отличие от method expression, где функция создается с дополнительным параметром (ибо ресивер вычислить невозможно).

Type assertions: `var v, ok = x.(T)`

Variadic functions: список параметров (variadic) внутри функции виден как слайс.
Function `f := func(xs ...int){}` could be called as `var ys = []int{3, 7}; f(ys...)`

Generic functions (types): it's like templates, instantiation creates a new non-generic function/type.

> Arithmetic operators ... yield a result of the same type as the first operand
Но это не так для Constant Expressions!

Floating-point operators:
> An implementation may combine multiple floating-point operations into a single fused operation

Comparison operators:
> Two string values are compared lexically byte-wise
Two channel values are equal if they were created by the same call to `make`.
Slice, map, and function types are not comparable.
Boolean, numeric, string, pointer, and channel types are strictly comparable

Receive operator:
> Receiving from a nil channel blocks forever

Conversion:
> nil is not a constant ...
Converting a constant to a type (that is not a type parameter) yields a typed constant ...
conversions only change the type but not the representation of x (not applicable to `numeric <=> string`) ...
there is no indication of overflow (int) ...
the conversion succeeds (float) but the result value is implementation-dependent ...

Забудьте про адресную арифметику:
> There is no linguistic mechanism to convert between pointers and integers. 
The package `unsafe` implements this functionality under restricted circumstances

Исключение из правил (`string(rune(x))` vs `strconv.Itoa`):
> an integer value may be converted to a string type

> converting a slice to an array pointer yields a pointer to the underlying array of the slice

Constant expressions:
тип результата бинарной операции (над нетипизированными операндами) определяет тип самого правого операнда.
Что противоположно правилам Arithmetic Operators!

Order of evaluation:
порядок вычислений операндов определён не во всех случаях, всегда учитывай порядок вычисления операндов в выражениях, избегай магии!
> At package level, initialization dependencies determine the evaluation order ... but not for operands within each expression ...
all function calls, method calls, and communication operations are evaluated in lexical left-to-right order ...

`For` statements with `range` clause:
При итерировании строки имеем коллекцию рун, но ключ показывает индекс первого байта руны (не забывайте однако, строка это коллекция байт).
> If the range expression is a channel, at most one iteration variable is permitted, otherwise there may be up to two ...
The range expression is not evaluated: if at most one iteration variable is present and `len(x)` is constant

Select statements:
выбирает (произвольно) один из готовых комм.кейсов. Если нет готовых, то блочится или уходит в default.

Return statements:
Именованные возвращаемые параметры фунции/метода, их не требуется явно указывать в `return`.
> Any functions deferred by F are executed before F returns to its caller

Break statements:
Если есть label, то метка указывает блок из которого надо выйти. Это НЕ goto!

Continue statements:
Если есть label, то метка указывает блок где надо начать следующую итерацию. Это НЕ goto!

Fallthrough statements:
Противоположность `break` в традиционном `switch`. Чтобы провалиться надо явно написать `Fallthrough`.

Defer statements:
Типа деструктора (finally) для функции. Вызов блока перед возвращением из функции.
Несколько defer собираются в стек (FILO).

Built-in functions:
> The built-in functions do not have standard Go types, so they can only appear in call expressions;
they cannot be used as function values

`append`: добавляет элементы в слайс, элементы-слайса в слайс, байты-строки в слайс.

`copy`: копирует элементы слайса, байты строки.

`clear`: слайс - забивает zero-значениями по индексам 0..len; map - удаляет пары. (C - consistensy)

`close(channel)`:
> Closing a receive-only channel is an error.
Receive from closed ch - non-blocking zero value.
Other ops - panic. (C - consistency)

`make`: slice, map or channel. Returns `T` (not `*T`) and init memory.
`new`: allocate storage for a variable, returns `*T` (not `T`).

`panic(not_nil)`: propagate to top, terminate program. Deferred block still called.
`recover`: returns nil if no panic to recover from.

`print`, `println`: наличие этих функций не гарантировано.

Bootstrapping:
package `init` function.
> Multiple such functions may be defined per package, even within a single source file

unsafe: адресная арифметика, размеры, выравнивание структур, etc. Не надо, но если очень хочется, то можно.

Size and alignment:
> A struct or array type has size zero if it contains no fields (or elements, respectively) that have a size greater than zero.
Two distinct zero-size variables may have the same address in memory.
Отсюда вытекает трюк с коллекцией охулиарда пустых структур.

#### go.mod

https://go.dev/ref/mod#glossary

A `module` is a collection of packages that are released, versioned, and distributed together
Модуль: про дистрибуцию (пакет: про изоляцию).

A `module` is identified by a `module path`, which is declared in a `go.mod` file, together with information about the module’s dependencies

The `module root directory` is the directory that contains the `go.mod` file

The `main module` is the module containing the directory where the `go` command is invoked

A `module path` is the canonical name for a module, declared with the module directive in the module’s `go.mod` file.
A module’s path is the prefix for `package paths` within the module

Each `package` within a `module` is a collection of source files in the same directory that are compiled together.
A `package path` is the `module path` joined with the subdirectory containing the package

Typically, a `module path` consists of a repository root path, a directory within the repository (usually empty), and a major version suffix (only for `major version` 2 or higher)

If the module is released at `major version` 2 or higher, the `module path` must end with a `major version suffix` like `/v2`

A `version` identifies an immutable snapshot of a module, which may be either a `release` or a `pre-release`.
Each version starts with the letter `v`, followed by a `semantic version` https://semver.org/spec/v2.0.0.html

The `patch version` may be followed by an optional pre-release string starting with a `hyphen`.
The pre-release string or patch version may be followed by a `build metadata` string starting with a `plus`.
For example, `v0.0.0`, `v1.12.134`, `v8.0.5-pre`, and `v2.0.9+meta` are valid versions.

A version is considered `unstable` if its major version is `0` or it has a `pre-release suffix`
Unstable versions are not subject to compatibility requirements

A `pseudo-version` is a specially formatted `pre-release` version that encodes information about a specific `revision` in a version control repository

A module may be checked out at a specific `branch`, `tag`, or `revision` using a `version query`. `go get example.com/mod@master`

Module versions are distributed as `.zip files`.
There is rarely any need to interact directly with these files, since the `go` command creates, downloads, and extracts them automatically 

The `module cache` is the directory where the `go` command stores downloaded module files.
The module cache is distinct from the `build cache`, which contains compiled packages and other build artifacts
The default location of the module cache is `$GOPATH/pkg/mod`.
To use a different location, set the `GOMODCACHE` environment variable.

A `workspace` is a collection of modules on disk that are used as the `main modules` when running `minimal version selection` (`MVS`)
Go uses an algorithm called `Minimal version selection` (MVS) to select a set of module versions to use when building packages
`MVS` operates on a directed graph of modules, specified with `go.mod` files
`MVS` starts at the `main modules` (special vertices in the graph that have no version) and traverses the graph
Т.е. воркспейс содержит набор вершин графа зависимостей.

Most `go` commands may run in `Module-aware` mode or `GOPATH` mode
If `GO111MODULE`=off, the go command ignores `go.mod` files and runs in `GOPATH` mode

When using modules, the `go` command typically satisfies dependencies by
downloading modules from their sources into the `module cache`, then loading packages from those downloaded copies

`Vendoring` may be used to allow interoperation with older versions of Go, or
to ensure that all files used for a `build` are `stored` in a single file tree.

The `go mod vendor` command constructs a directory named `vendor` in the `main module’s root directory` containing copies of all packages needed

If the `vendor` directory is present in the main module’s root directory, it will be used automatically

### Переменные, базовые типы данных

- [vars_1](week_01/vars_1.go)
- [vars_2](week_01/vars_2.go)
- [strings](week_01/strings.go)
- [const](week_01/const.go)
- [types](week_01/types.go)
- [pointers](week_01/pointers.go)

```s
pushd sandbox
mkdir -p week01 && pushd ./week01
touch main.go
go mod init week01
go mod tidy
popd
go work use ./week01/
go vet week01
gofmt -w week01
go run week01
```
[week 1 playground](./sandbox/week01/main.go) `GO_APP_SELECTOR=week01 gr`

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

`float32, float64` но нет просто `float`. (C - consistensy)

Есть `complex64, complex128` математики и физики радуются.

Строки в кавычках интерпретируются, символы типа `\n` и прочие будут транслированы.

Строки в бэктиках не интерпретируются.

Одинарные кавычки для задания byte (uint8) или rune (uint32).

Строки immutable.

Длина строки считается в байтах, индексируется (доступ) по байтам.
Для подсчета в символах используй `utf8.RuneCountInString(someStr)`.
Итерация range по рунам.
Соответственно, срезы тоже в байтах.
Строки можно легко конвертировать в байты и байты в строки.

Константы `const name = value`.
Блоки констант в скобках `( ... )`.
Опредение через `iota`, автоинкремент, всё сложно.
Нетипизированные константы, тип присваивается при записи константы в переменную. Вроде макроса получается.

Пользовательские типы данных, `type`. Полезно при моделировании, DSL.
Нет автоматического приведения типов.

Нет адресной арифметики, но есть ссылки, reference. Полезно для передачи структур без копирования.
- `b := &a` получение ссылки.
- `*b = 42` запись значения в переменную по ссылке.
- `c := new(int)` создание именованной ссылки на безымянную переменную, инициализация нулевыми значениями.

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

- `var buf []int` создание пустого слайса без инициализации, nil
- `buf := []int{} // len:0, cap:0` создание пустого с инициализацией
- `buf := make([]int, 5, 10) // len:5, cap:10` срез это конструкция над массивом, у среза есть длина и емкость.
При длине более 0, элементы инициализируются нулевым значением типа элемента.

- `buf = append(buf, 9, 10) // len:2, cap:2` буфер может расти, при исчерпании емкости буфер пересоздается с удвоенной емкостью.
- `buf = append(buf, otherBuf...)` при добавлении другого буфера, его надо "распаковать".

Слайсы могут работать "по ссылке", оперируюя значениями в одном и том-же буфере.
То есть, если явно не выделять память под слайс (или неявно, через append), то работа идёт в одном и том-же буфере.
Т.е. срез это view на нижелижащий массив, который может быть пересоздан при операциях со слайсом.
Т.е. рассчитывать на то, что производные слайсы ссылаются на один буфер (в общем случае) - нельзя.

- `numCopied = copy(emptyBuf, buf)` копирование элементов в другой буфер, внутри проверка на выход за границы.
Копируется наименьшее из двух `len` колич.элементов.
- `copy(buf[1:3], []int{5, 6})` копирование под-диапазона.

Мапки `var user map[string]string`, можно, как и слайс, создать с нужной ёмкостью, через `make`.
- `mName, mNameExists := user["middleName"]` правильный способ получения значения из мапки, ибо по умолчанию, несуществующее значение = zero value.
- `delete(user, "lastName")` удаление записи из мапки.

Мапки реализованы через бакеты, бакеты могут пересоздаваться, поэтому любые референсы на значения невозможны.

### Управляющие конструкции

- [control](week_01/control.go)
- [loop](week_01/loop.go)

- `if boolVal { ... }` только тип bool допустим внутри иф.
- `if v, exists := myMap["name"]; exists { ... }` условие с блоком инициализации, здесь этот блок = `v, exists := myMap["name"];`
- `switch len(myMap) { 	case 0, 1: ... }` по умолчанию делает break при срабатывании условия
- `switch ... case k == "name" && v == "Bender": ...` сложные условия в switch
- `switch ... break` оператор выхода, можно выходить через несколько уровней, по метке;
Метка указывает блок кода из которого надо выйти, это не goto.

Циклы определяются ключевым словом `for`, есть разные формы циклов.

`for bytePosition, symb := range myStr { ... }` итерирование строки делается по символам (рунам), не байтам. (C - consistency)

### Основы функций

- [functions](week_01/functions.go)

`func sqrt(x int) int { ... }` и несколько более сложных вариантов объявления.
Например, именованный возвращаемый результат, иногда бывает удобно.

`func namedWithError(condition bool) (res int, err error) { ... }` осторожнее со значениями "по умолчанию".

`func sum(in ...int) int { ... }` списки параметров и кортежи возвращаемых значений -- это нормально.
Переменное количество входных параметров базируется на представлении параметров как слайса.

Синтаксис объявления функции работает (полноценно) только на уровне пакета.

### Функция как объект первого класса, анонимные функции

- [firstclass](week_01/firstclass.go)

Функция как значение переменной -- присваивать, передавать, возвращать.

`printer := func(msg string) { ... }` анонимная функция как значение переменной.

Замыкание, пример функции-фабрики, которая возвращает функцию печати-с-префиксом, где префикс берется из замыкания.
(Как люди мучаются без каррирования)

### Отложенное выполнение и обработка паники

- [defer](week_01/defer.go)
- [recover](week_01/recover.go)

- `defer doStuff("after work ...")` будет вызвана перед выходом из блока. Полезно как код финализации процедур.
- Несколько `defer` складываются в стек (FILO).
- Аргументы отложенных функций вычисляются НЕ отложенно а сразу. Чтобы этого избежать, такие аргументы заворачиваются в анонимную функцию.

`defer` полезен при восстановлении из паники ибо вызывается даже при возникновении паники.
Если внутри `defer` вызвать `recover()`, то программа не вывалится в панику а продолжит работать штатно.

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

Можно смотреть на методы как на обычные функции, где добавлен нулевой параметр, как "ресивер" - 
копия объекта или референс на существующий объект.

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

`GOPATH` определяет корневую директорию, в которой будут под-директории системы `bin, pkg, src`.

Имя пакета это имя директории.
Приватные поля (пакета) определяются именованием с маленькой буквы, публичные поля -- с большой.

Доступ к приватным полям возможен только в коде пакета, где определено приватное поле.

Крупные пакеты, с большим количеством файлов, предпочтительнее мелких пакетов, с малым количеством файлов.

Внешние зависимости складываются в директории `vendor`.

### Основы работы с интерфейсами

- [basic](week_01/basic.go)
- [many](week_01/many.go)
- [cast](week_01/cast.go)

- Интерфейс это тип `interface` с сигнатурами функций (и, возможно, констрейнтами на core тип) `type Payer interface { Pay(int) error }`
- Некий тип может "содержать" реализацию интерфейса `func (w *Wallet) Pay(amount int) error { ... }`
- При вызове метода нам похер на конкретный тип, указываем интерфейс `func Buy(p Payer) { err := p.Pay(10) }`

Можно (нужно?) держать переменную (поле структуры) с типом нужного интерфейса `var p Payer; p = &Card{Balance: 100}`

Можно кастовать тип от интерфейса обратно к конкретному типу, предварительно сматчив тип. И/Или сделать `type assertion`.

Так реализуется полиморфизм (версия ad-hoc polymorphism). Добавил реализацию интерфейса к произвольному типу и радуйся.

Т.е. где-то объявляется интерфейс (набор методов); в других местах для нужных типов данных этот интерфейс реализуется.
В другом месте пишется логика с использованием интерфейсов.
В другом месте этот код используется, с передачей в него значений (переменных) нужного типа (поддерживающий интерфейс).

Побочка: неясно, как глядя на код сразу сказать, кто реализует какой интерфейс и в каком объеме.

### Пустой интерфейс

- [empty_1](week_01/empty_1.go)
- [empty_2](week_01/empty_2.go)

Пустой интерфейс, выглядит внутри как пара `(type, value)` и не накладывает ограничений.
Любое значение реализует пустой интерфейс. Синоним `any`.
Использовать его можно только если ручками делать проверку/приведение типа:
`func Buy(in interface{}) { if p, ok = in.(Payer) ...}`

Рассказал на примере `fmt.Printf`, как аргумент типа "пустой интерфейс" может быть подерган за разные методы
(опции форматирования определяют вызываемый метод). `Stringer`.

Короче: пустой интерфейс это `any`. Штука настолько мощная, насколько вредная.

### Композиция интерфейсов

- [embed](week_01/embed.go)

Композиция интерфейсов.
Как и структуры, интерфесы могут быть составлены из других интерфейсов.
При этом, как и структуры, охватывающий интерфейс включает методы вложенных интерфейсов.

### Написание Программы Уникализации (ПУ)

- [uniq](week_01/uniq.go)
- [data_map.txt](week_01/data_map.txt)

Пример программы: получает на вход файл и выводит только уникальные строки из этого файла.
Более точная спека: print line with drop-duplicates.

Два варианта: 
- на мапке строка - уникальность; 
- на предположении, что файл сортирован ascending (panic: if prev > current; dup: if prev == current)

Демонстрация абстрагирования по входу и выходу, выделение кода в функцию, удобную для тестирования.

### Написание тестов для ПУ

- [unique/unique](week_01/unique/unique.go)
- [unique/unique_test](week_01/unique/unique_test.go) `GO_APP_SELECTOR=week01_test gr`

Чтобы поддерживать тестирование, зависимости передаются в функцию как аргументы
`func sortedInputUnique(input io.Reader, output io.Writer) error { ... }`

Модуль тестов должен быть файлом с именем с суффиксом `_test.go`.
Должен содержать функции-тесты с именами с префиксом `Test`:
`func TestSortedInput(t *testing.T) { ... }`

Тесты запускаются `go test -v ./unique`

Демонстрация написания примитивных юнит-тестов.

### week 1 homework

> Описание задания и локальные тесты вы можете найти в папке 99_hw в коде к соответствующей неделе.
Задание предназначено для самостоятельного решения, материалов лекций достаточно для его выполнения. 
Не надо ничего гуглить, уж тем более решения. Задача научиться приходить к решению, а не увидеть его.

Дана готовая структура пакета `main`, реализующего поведение утилиты `tree`.
Надо в пакет добавить реализацию функции,
выполняющей обход дерева каталогов и выдающей строки для вывода пользователю.
- [homework materials](week_01/materials.zip/week_1/99_hw/tree/)
- `handouts\golang_web_services_2023-12-28.zip\1\99_hw\tree\readme.md`
- [actual homework project](./sandbox/week01_homework/tree//hw1.md) `sandbox/week01_homework/tree/main.go`

```s
pushd week01_homework/tree
go mod init tree
go mod tidy
pushd ..
go work init
go work use ./tree/
go vet tree
gofmt -w tree
go test -v tree
go run tree . -f
cd tree && docker build -t mailgo_hw1 .
```
`GO_APP_SELECTOR=week01_tree_test gr` Реализовал два варианта: рекурсивный и нет (стек).

## part 1, week 2

Асинхронная работа.
[Код, домашки, литература](week_02/w2_materials.zip), updated: `handouts\golang_web_services_2023-12-28.zip\2\`

```s
pushd sandbox
mkdir -p week02 && pushd ./week02
touch main.go
go mod init week02
go mod tidy
popd # to sandbox
go work use ./week02/
go vet week02
gofmt -w week02
go run week02
```
[week 2 playground](./sandbox/week02/main.go) `GO_APP_SELECTOR=week02 gr`

### Плюсы неблокирующего IO (на примере web)

Рассказ про:
- Асинхронное выполнение (помните AJAX?);
- Скорость вычислений: поток данных процессор-кеш-память;
- Время на переключение контекста задач процессором (выгрузка-загрузка регистров);
- Современные тенденции на многоядерность и параллельность;
- Процессы тяжелые, потоки легче, асинхронные сопрограммы (green threads) еще легче.
- Утилизация дорогого железа -- сервер должен работать.
- Невытесняющая многозадачность (eventloop, Windows 3.0) vs preemptive. 

Ввод-вывод (IO) и ожидание возврата из syscall.
Время ожидания можно потратить на другие задачи, non-blocking IO.
Задачи IO-bound vs CPU-bound.

В целом: нужна асинхронность и параллельность, язык содержит средства получить это эффективно (потоки, горутины, каналы).

Node.js быстр но на одном ядре, Go еще быстрее
за счет утилизации всех ядер на комбинации подходов вытесняющей многозадачнисти и НЕ-вытесняющей.
Горутины в одном потоке исполняются на eventloop, но потоков есть по количеству ядер.
Все это в одном адресном пространстве одного процесса.

`Communicating Sequential Processes` by Tony Hoare.
> formal language for describing patterns of interaction in concurrent systems.
... member of... mathematical theories of concurrency... based on message passing via channels.

CSP это как акторы, но это только если вам пох на детали.
Каналы анонимны - акторы нет; 
сообщение явно передается в канал с неизвестным получателем, актор явно задается как получатель;
канал синхронизирует передачу сообщений (гарантия доставки), чтение сообщения актором не отслеживается в явном виде.
https://en.wikipedia.org/wiki/Communicating_sequential_processes#Comparison_with_the_actor_model

Горутины перемещаются между системными потоками,
код горутины может быть выполнен в любом потоке (как перемещается стек?).

Горутины очень легкие, их может быть оченьм много.
Коммуникация через "каналы".

### Горутины: легковесные "процессы"

Демонстрация асинхронной работы нескольких горутин
- [goroutines](week_02/goroutines.go)

`go doSomeWork(i)` функция не может вернуть значение обычным способом (предполагается общение через каналы).

`runtime.Gosched()` как способ уйти в планировщик, дав возможность запустить другие горутины (yield).

Есть шанс заблокировать шедулер, если молотить цикл без вызовов системных функций (вспоминаем eventloop Win 3.0).

### Каналы: передаём данные между горутинами

- [chan_1](week_02/chan_1.go)
- [chan_2](week_02/chan_2.go)

`chan` keyword, it is a type, like `int`.
Под капотом, при отправке/приеме данных через канал происходит передача контроля над данными между потоками/горутинами.
Каналы как средство синхронизации асинк. тасок.

Чтение, запись -- операторы стрелочка `x <- myChannel; myChannel <- x`.
В сигнатуре функции можно уточнить тип канала: только на чтение или только на запись.

- `ch1 := make(chan int, 0)` небуферизованный канал.
- `ch1 <- 42` запись в канал (блокирует), читатель уже должен ждать, ибо небуферизовано.
- `close(ch1)` закрытие канала, сигнал EOF для читателя.
- `v := <-in` чтение из канала.

Рабочий паттерн: создать канал, запустить горутину-читатель/писатель, запустить запись/чтение. Канал закрывает писатель.
- Чтение из канала в который никто не пишет: блочит горутину. Или fatal, deadlock.
- Запись в канал из которого никто не читает: fatal, deadlock. Или блок горутины.
- Запись в закрытый канал: паника.
- Чтение из закрытого канала: пустое значение и флаг (не)успеха.

Fatal, deadlock может произойти в "главной" горутине. Не важно, чтение или запись, важно главная или вторичная горутина.
Если читать/писать внутри вторичной горутины (в канал без второго конца), то эта вторичная горутина блочится на канале,
но fatal-deadlock не выпадает.

Небуферизованные каналы vs буферизованные, размер буфера. Работа с каналом в цикле.
Записанное в буферизованный канал может пропасть без следа, если никто не прочтёт. Поэтому буферизация это зло.

Если писать в небуферизованый канал (в главной горутине), из которого никто не читает, будет deadlock, очень плохо.
Наоборот -- тоже беда, нельзя читать из канала, в который никто не пишет.

Система отслеживает закрытие горутины, связанной с каналом и повисший канал вызывает ошибку.

Проблемы блокировки чтения/записи (и другие) позволяет решить оператор `select`, см.следующий урок.

### Мультиплексирование каналов: оператор `select`

- [select_1](week_02/select_1.go)
- [select_2](week_02/select_2.go)
- [select_3](week_02/select_3.go)

`select { case val := <-ch1: ...; case ch2 <- 1: ...; default: ... }`
Оператор `select` чекает каналы без блокирования.
Что-то делает только если канал готов. Если все каналы остановлены, то действие "по умолчанию".
В одном проходе выполняет только одну операцию, какую - рантайм определит, мы не можем повлиять.

Поместить `select` в цикл, действие по умолчанию если ни один канал не работоспособен -- можно выходить из цикла.

Использование канала "команд" вместе с каналом "данных": если получили команду "на выход" -- завершаем цикл.

### Таймеры и таймауты (как источник сигнала в каналах)

- [timeout](week_02/timeout.go)
- [tick](week_02/tick.go)
- [afterfunc](week_02/afterfunc.go)

- Таймер как источник сигнала (канал таймера) "через определенное время". Помогает в реализации логики тайм-аутов.
- Тикер как источник регулярных/периодических сигналов в канале. С поддержкой backpressure.
- AfterFunc как способ отложенного выполнения функции через заданное время.

Некоторые фунции пакета `time` могут быть удобны, но могут приводить к утечкам памяти, создавая вечные таймеры.
Пользуйтесь функциями, создающими таймеры, на которых можно вызвать метод `Stop`.

### Пакет `context` и отмена выполнения

- [cancel](week_02/context_cancel.go)
- [timeout](week_02/context_timeout.go)

Отмена асинхронных операций, вручную: с помощью примитивов пакета `context`.
Всем воркерам выдается общий контекст `ctx, cancelFunc := context.WithCancel(context.Background())`, в котором определена фунция завершения.
Когда получен нужный результат, дергается функция завершения контекста и воркеры понимают (слушая канал `ctx.Done()`), что пора выходить.

Отмена асинхронных операций, по таймауту: `ctx, cancelFunc := context.WithTimeout(context.Background(), workTime)`.
В этом варианте, в канал `ctx.Done()` приходит сообщение по таймеру, не по ручному вызову функции завершения.
Чтобы освободить ресурсы надо вызвать `cancelFunc`, так написано в докстринге.

Контекст -- основной способ отмены асинхронных операций, завершения воркеров, выполняющих циклы или долгие операции асинхронно.

### Асинхронное получение данных

- [async_work](week_02/async_work.go)

Практическое применение горутин и каналов: распараллеливание считывания статей и комментариев к ним на сайте.
Обработка страницы запускает (асинк) обработку комментариев, каменты получает в канале, созданном функцией обработки каментов.
В итоге страница и ее каменты обрабатываются параллельно.

Буферизованный канал как костыль для поддержки ситуации: писатель пишет в канал, читатель которого решил преждевременно сдохнуть.

### Пул воркеров

- [workerpool](week_02/workerpool.go)

Один канал как очередь сообщений: создается в главной таске, эта таска пишет туда задания (она же его и закрывает).
Несколько горутин как воркеры читающие из общей очереди. Рантайм сам распределяет кому-что.
По закрытию очереди, воркеры заканчивают. Это важный момент, надо следить, чтобы закрытие отработало.

### `sync.Waitgroup` -- ожидание завершения работы

- [waitgroup](week_02/waitgroup.go)

До сих пор мы использовали ввод с клавиатуры, чтобы программа не завершилась раньше своих воркеров. Теперь так делать не нужно.

Ресурс `wg := &sync.WaitGroup{}`, увеличивается `wg.Add(1)` при добавлении воркеров и уменьшается `defer wg.Done()` при удалении воркеров.
Можно использовать для ожидания `wg.Wait()` завершения всех воркеров.

Копировать созданную группу нельзя, работать с ней только по ссылке.

Добавлять в группу следует перед строчкой `go ...` ибо если унести добавление в код горутины, рантайм может добежать до
`Wait` раньше инициализации горутин.

Примитивы синхронизации, [sync package](https://pkg.go.dev/sync),
[sync/atomic package](https://pkg.go.dev/sync/atomic).

### Ограничение по ресурсам (семафор на каналах)

- [ratelim](week_02/ratelim.go)

Демонстрация ограничения количества одновременно работающих воркеров.

Канал (буферизованный) используется как семафор: запись сообщения в канал забирает ресурс, считывание сообщения освобождает ресурс.

Использование буферизованного канала (квота) как ограничителя.
Как буфер исчерпался, новая работа не начинается.

Перед стартом работы, воркер пытается записать сообщение в канал квоты.
Если в буфере место ещё есть, воркер запишет сообщение и сможет продолжить работу.
На выходе воркер считывает сообщение из канала квоты, освобождая место в буфере.

Почему `struct{}` считается хорошим вариантом сообщения в канале (когда важен только сам факт наличия сообщения)?
> Using a channel of empty structure will only increment a counter in the channel but not assign memory, copy elements ...
https://mariadesouza.com/2019/01/12/empty-struct/

> the empty struct has a width of zero. It occupies zero bytes of storage
https://dave.cheney.net/2014/03/25/the-empty-struct
https://pkg.go.dev/github.com/bradfitz/iter#example-N

### Ситуация гонки на примере конкурентной записи в map

- [race_1](week_02/race_1.go)

Демо: X воркеров конкурентно пишут в одну мапку.
Мапка по ходу операций перестраивается и разные воркеры начинают работать с разными копиями данных.
Кто из них победит? Программа падает с fatal error.

`go run -race ...` для диагностики гонки.

Что же делать? Ставить блокировки?

### `sync.Mutex` для синхронизации операций

- [race_2](week_02/race_2.go)

Берем мютекс (ссылка для предотвращения копирования) `mu := &sync.Mutex{}` и в нужном месте используем его
`mu.Lock(); counters[th*10+j]++; mu.Unlock()`

Вроде бы ОК. Но есть нюансы ...

### `sync.Atomic`

- [atomic_1](week_02/atomic_1.go) демо не-синхронизированных операций
- [atomic_2](week_02/atomic_2.go) демо с синхронизацией

Использовать Mutex слишком дорого для атомарного изменения одной переменной, счетчика например (регистра ЦП).
Есть альтернатива `atomic.AddInt32(&totalOperations, 1)`.

### week 2 homework

> асинхроннй пайплайн

Есть пакет `main` (file `signer.go`), он требует добавления кода. Есть тесты. Подробности в `*.md` файлах.
Вкратце: библиотека реализует "пайплайн" набора джобов, с распараллеливанием работы.
Каждая джоба читает из канала `in` и пишет в канал `out`. Если тесты выполняются, задача решена.
- [homework materials](week_02/w2_materials.zip/2/99_hw/signer/hw2.md), [extra](week_02/w2_materials.zip/2/99_hw/wp_extra/wp_extra.md)
- [actual homework project](./sandbox/week02_homework/signer/hw2.md)
- [actual homework project, extra](./sandbox/week02_homework/wp_extra/wp_extra.md)

Updated: `handouts\golang_web_services_2023-12-28.zip\2\99_hw\signer\readme.md`

```s
pushd sandbox # workspace
mkdir -p week02_homework/signer

pushd week02_homework/signer
go mod init signer

cat > signer.go << EOT
package main
func ExecutePipeline(jobs ...job) { panic("not yet") }
var SingleHash job = func(in, out chan interface{}) { panic("not yet") }
var MultiHash job = func(in, out chan interface{}) { panic("not yet") }
var CombineResults job = func(in, out chan interface{}) { panic("not yet") }
func main() { panic("not yet") }
EOT

go mod tidy

popd # workspace
# go work init
go work use ./week02_homework/signer

go vet signer
gofmt -w ./week02_homework/signer
go test -v -race signer
go run signer
```
signer app `GO_APP_SELECTOR=week02_signer_test gr`

Extra задание (динамический пул воркеров) без тестов и с очень свободным описанием. Отложил на "потом".

## part 1, week 3

Работа с динамическими данными и производительность. Рефлексия, кодогенерация, бенчмарки.
[Код, домашки, литература](week_03/part1_week3.zip), updated: `handouts\golang_web_services_2023-12-28.zip\3\`

```s
pushd sandbox
mkdir -p week03 && pushd ./week03

cat > main.go << EOT
package main
func main() { panic("not implemented yet") }
EOT

go mod init week03
go mod tidy
popd # to sandbox
go work use ./week03/
go vet week03
gofmt -w week03
go run week03
```
[week 3 playground](./sandbox/week03/main.go) `GO_APP_SELECTOR=week03 gr`

### JSON codec

- [json](week_03/json.go)
- [struct_tags](week_03/struct_tags.go)

- Пакет `encoding/json` дает нам кодек. `json.Unmarshal(bytes, emptyStructRef)`: Декодирование из слайса байт в структуру;
надо знать тип структуры и создать её, пустую, перед декодированием.
- Кодирование делается зеркально, `bytes, err := json.Marshal(someStruct)`.
- Приватные поля (с маленькой буквы которые) не обрабатываются, ибо `encoding/json` пакет не может получить доступ
к приватным полям нашего пакета.

Теги структуры: метаинформация структур.
При описании (определении) стурктуры, к полю добавляется строка определенного формата,
где записана метаинформация. В частности, как декодеру json обрабатывать поле, какое имя ему дать, какой тип использовать, etc.
```s
ID       int `json:"user_id,string"`
```

https://go.dev/blog/json

### JSON codec, структура неизвестна

- [dynamic](week_03/dynamic.go)

Как быть, если мы точно не знаем структуру, представленную json текстом?
Unmarshal в пустой интерфейс.

Также и Marshal, создав мапку `map[string]interface{}{ ... }`, можно ее кодировать в json.

### Пакет reflect, работаем с динамикой в рантайме

- [reflect_1](week_03/reflect_1.go)
- [reflect_2](week_03/reflect_2.go)

Демонстрация работы рефлексии для анализа (и процессинга) полей структур в рантайме.
Получить список полей и в цикле каждое обработать, посмотрев на имя, тип, тег.
Примитивный пример декодинга бинарных данных в произвольную структуру (при условии, что слайс байт отражает именно эту структуру).

- `reflect.ValueOf(x).Elem()` позволяет итерировать типы и значения структуры, переданной как "пустой интерфейс".
- Пример, как при помощи рефлексии распаковать слайс байт (бинарные данные на выходе perl `pack`) в структуру.
`reflect.ValueOf(x).Elem()` для получения списка полей структуры, которую надо восстановить из слайса байт.
Перебирая поля структуры, читаем данные из источника байт.

### Кодогенерация, программа пишет программу

`handouts\golang_web_services_2023-12-28.zip\3\codegen\Makefile`
- [unpack](week_03/unpack.go)
- [marshaller](week_03/marshaller.go)
- [codegen](week_03/codegen.go)

Демонстрация распаковки слайса байт в структуру (как в примере с рефлексией).
Только теперь жестко задан процесс распаковки, для каждого поля структуры записан (сгенерирован) шаблонный код.
Кодогенератор написан (ручками) с использованием анализатора AST го-шной программы.

Структура записана в файле х, кодеген читает этот файл и пишет шаблонный код (метод структуры) в файл у.
Сама программа собирается из файлов х, у, ...

К структуре добавлен метод `Unpack`, сгенерированный кодогенератором.
Метод реализует восстановление структуры из бинарных данных, созданных через perl `pack`.
Кодогенератор анализирует гошную структуру (пользуясь гошным AST), считанную из файла с исходником,
и, для всех её полей, создает код восстановления поля из слайса байт.

Кодогенерация полезна, когда некогда в рантайме тратить время на анализ и рефлексию.
Профилирование покажет насколько падает скорость (спойлер: в 2..4 раза).

### Система бенчмарков Go

тестируем: кодоген vs рефлексия; регулярки vs подстрока; slice.append prealloc vs no-prealloc

- [unpack_test](week_03/unpack_test.go)
- [json_test](week_03/json_test.go)
- [string_test](week_03/string_test.go)
- [prealloc_test](week_03/prealloc_test.go)

Бенчмарки в `testing.B`, `func BenchmarkFooBar( ... ) { ... }`, `go test -bench ...`

Бенчмарки это как тесты, только тестируется скорость.
Файл содержит суффикс `_test.go`, имя функции содержит префикс `func Benchmark`.

- `go test -bench . week03` показал, что распаковка структуры через кодогенерацию вдвое быстрее рефлексии.
- `go test -bench . -benchmem week03` замер расхода памяти (в дополнение к замеру скорости). Рефлексия жрет вдвое больше памяти.
- Пакет кодека json на кодогенерации (EasyJson) в 2..4 раза эффективнее стандартного (не в каждой метрике каждой операции, однако).
- Демонстрация медленности регулярных выражений на строках. Регулярка в 3..x раз медленнее простого `contains`.
- Демонстрация медленности добавления в слайс без преаллокации. В 8 раз медленнее.

### Профилирование через pprof

https://www.google.com/search?q=go+tool+pprof+illustrated&tbm=vid

Собрать профиль
`go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1 unpack_test.go`

Анализ профиля
`go tool pprof main.test.exe mem.out`

Команды: `top`, `list Unpack`, `web`, `alloc_space, top`, `alloc_objects, top`

Можно увидеть строки программы, где происходило "горячее": выделение памяти или жрал цпу.

Можно снимать дампы с работающей программы, не прерывая её работы.

### `sync.Pool`, mem allocation

- [pool_test](week_03/pool_test.go)

Демонстрация разницы в подходах: выделение буфера 64 байта для энкодера джейсон и выполнение кодирования, на каждой итерации теста.
Второй вариант: используя пул, память под буфер выделить один раз.

Хотим не выделять память каждый раз, хотим переиспользовать имеющийся буфер. Для скорости.

`dataPool = sync.Pool{ ... }; data := dataPool.Get().(*bytes.Buffer)`
Приведение к типу: `.(*bytes.Buffer)`, ибо геттер дает `any`.

sync.Pool это по сути синхронизированный синглтон, интерфейс: Get, Put.

Операций выделения памяти меньше. Нагрузка на GC, соответственно, падает.

### Покрытие кода тестами

Test coverage
- `go test -v -cover`,
- `go test -coverprofile=cover.out`,
- `go tool cover -html=cover.out -o cover.html`

Демо кода с тремя ветками выполнения, покрытие тестами всех веток.

### XML, complete file vs stream of tokens, performance

Паттерн борьбы за производительность: обработка потока (токенов) как замена обработке целого файла.
- [main](week_03/xml_main.go)
- [xml_test](week_03/xml_test.go)

Декодинг файла xml целиком и сразу, vs потоковая обработка токенов. Производительность ~1.5.
Но основная проблема не в производительности а в количестве ресурсов (памяти) потребных для обработки.
Для потока памяти нужна константа, для файла памяти надо по размеру файла.

Поточная обработка данных, на примере xml. `encoding/xml`.
Читаем по токенам, обрабатываем: `xml.NewDecoder(bytes.NewReader(data)).Token()`

### week 3 homework

Есть пакет `main`, он требует добавления кода в функцию `FastSearch`.
Есть тесты.
Подробности в `*.md` файлах.
`handouts\golang_web_services_2023-12-28.zip\3\99_hw\readme.md`

Вкратце: задание на оптимизацию (быстрее молотит, меньше памяти жрет, аллокаций меньшее количество) уже существующей (baseline) функции.
Научиться работать с `pprof`.
Есть функция `SlowSearch`, надо написать `FastSearch`, более производительную (во всех смыслах) версию медленной.
- [homework materials](week_03/part1_week3.zip/part1_week3/99_hw/hw3.md)
- [actual homework project](./sandbox/week03_homework/finder/hw3.md) `GO_APP_SELECTOR=week03_finder_test gr`

> Для выполнения задания необходимо чтобы
один из параметров ( ns/op, B/op, allocs/op ) был быстрее чем baseline ( fast < solution ) и
ещё один лучше baseline + 20%.

```s
pushd sandbox # workspace
mkdir -p week03_homework/finder

pushd week03_homework/finder
go mod init finder

cat > finder.go << EOT
package main
func main() { panic("not yet") }
EOT

go mod tidy

popd # workspace
# go work init
go work use ./week03_homework/finder

go vet finder
gofmt -w ./week03_homework/finder
go run finder
go test -v finder
go test -bench . -benchmem finder    
```
finder lib.

```s
# go test -bench . -benchmem $module
BenchmarkSlow-8               46          30084826 ns/op        20196022 B/op     189848 allocs/op
BenchmarkFast-8               78          16459328 ns/op          722669 B/op      12314 allocs/op
```
bench results.

https://betterstack.com/community/guides/scaling-go/json-in-go/

## part 1, week 4

Основы HTTP. Разработка web-сервера.
[Код, домашки, литература](week_04/part1_week4.zip), updated: `handouts\golang_web_services_2023-12-28.zip\4\`

```s
pushd sandbox
mkdir -p week04 && pushd ./week04

cat > main.go << EOT
package main
func main() { panic("not implemented yet") }
EOT

go mod init week04
go mod tidy
popd # to sandbox
go work use ./week04/
go vet week04
gofmt -w week04
go run week04
```
[week 4 playground](./sandbox/week04/main.go) `GO_APP_SELECTOR=week04 gr`

### Слушаем TCP-сокет с использованием пакета net

- [net_listen](week_04/net_listen.go)

Открытие TCP соединения, получение коннекта в цикле, общение с подключившимся клиентом в отдельной горутине.
`go build -o net_listen . && ./net_listen`
`netcat 127.0.0.1 8080`

Демонстрация многопоточности сервера, несколько клиентов могут общаться с сервером одновременно.

### Обслуживание HTTP-запросов

- [http](week_04/http.go) `net/http` пакет. Простейший вариант с одним хендлером одного роута.

- [pages](week_04/pages.go) Обработка трех роутов тремя хендлерами, два из них -- анонимные функции. 
Роут `/pages/` позволяет обрабатывать любое количество под-страниц. Роут `/page` задает только одну обрабатываемую страницу.

- [servehttp](week_04/servehttp.go) Хендлеры запросов как структуры, реализующие интерфейс.
Параметризуемые хендлеры, обладающие внутренним состоянием. Например, коннект к БД.
Все запросы на заданный урл обработаются хендлером привязанным к одному инстансу структуры.

- [mux](week_04/mux.go) Более низкоуровневое API к веб-серверу.
Отдельно создается и конфигурируется мультиплексор запросов (обработчики роутов); отдельно создается и конфигурируется сервер.

- [servers](week_04/servers.go) Несколько сервисов в разных горутинах, в рамках одной программы.
Сервис + мониторинг, к примеру можно так сделать.

### Работа с параметрами запросов

- [get](week_04/get.go) Как прочитать параметры из объекта `Request`. Параметры переданы через GET, PUT, POST -- всяко.
- [post](week_04/post.go) `if r.Method == http.MethodPost { loginFormValue := r.FormValue("login") }`, получение данных из POST формы.

- [cookies](week_04/cookies.go) Cookie, чтение, запись, устаревание.
Бизнес-логика логин/логаут с использованием cookie как идентификатора сессии пользователя.

- [headers](week_04/headers.go) Установка и чтение заголовков HTTP, header.

### Обслуживание статичных данных

- [static](week_04/static.go) `http.FileServer` реализует выдачу статичных файлов (без их интерпретации).
В проде так не надо, для локальной разработки (пет-проектов) сойдет.

### Upload файлов

- [file_upload](week_04/file_upload.go) Два варианта выгрузки файла:
    `Request.FormFile` даёт нам загруженный в форму (POST Multipart) файл.
    `io.ReadAll(Request.Body)` даёт нам байты тела запроса, загруженные через POST.

### HTTP-запросы во внешние сервисы

- [request](week_04/request.go) Три варианта (разные уровни API) запроса по HTTP, с точки зрения клиента:
`http.Get(uri); http.DefaultClient.Do(req); httpClient.Do(req)`.

### Тестирование обработчиков HTTP-запросов/ответов

- [request_test](week_04/request_test.go) `go test -v $module`. Тестирование функций сервера (хендлера): подстановка моков риквеста и респонса: `httptest.NewRequest`, `httptest.NewRecorder`.
- [server_test](week_04/server_test.go) Тестирование функций зависящих от внешнего сервиса: мокать внешний сервис `httptest.NewServer(http.HandlerFunc(ExtApiHandlerMock))`.

### Шаблоны inline, шаблоны из файлов

- [inline](week_04/inline.go) `text/template` пакет для парсинга и вывода шаблонов (формирование текста ответа сервисом).
- [file](week_04/file.go), [users.html](week_04/users.html) шаблон из файла, рендеринг с автоматической экранизацией выводимого html кода.

### Вызов методов и функций из шаблонов

- [method.go](week_04/method.go), [method.html](week_04/method.html) В шаблон передана структура, хотим там (в шаблоне) вызвать её метод (без параметров).

- [func.go](week_04/func.go), [func.html](week_04/func.html) `template.FuncMap` для регистрации фунций, если методов структур нам недостаточно. В шаблоне вызываем регистрированные фунции (параметр - текущий объект `.`).

### Профилирование через pprof

# I_AM_HERE

- [pprof_1.go](week_04/pprof_1.go), [pprof_1.sh](week_04/pprof_1.sh)

`net/http/pprof` для профилирования работающей программы под нагрузкой.
При импорте пакета регистрируются спец. обработчики профилировщика,
при работе программы дергаются эти урлы для снятия профиля cpu or mem.
В офлайн уже можно анализировать профиль.

- `ab -t 300 -n 1000000000 -c 10 http://127.0.0.1:8080/` имитация нагрузки.
- `curl http://127.0.0.1:8080/debug/pprof/heap -o mem_out.txt`, `curl http://127.0.0.1:8080/debug/pprof/profile?seconds=5 -o cpu_out.txt`
    снятие профиля.
- `go tool pprof -svg -inuse_space pprof_1.exe mem_out.txt > mem_is.svg` анализ профиля.

### Поиск утечки горутин

- [pprof_2](week_04/pprof_2.go)

- `pprof` позволяет снять стектрейс всех горутин, на примере утечки горутин.
`curl http://localhost:8080/debug/pprof/goroutine?debug=2 -o goroutines.txt` снимаем дамп прямо под нагрузкой.

### Трассировка поведения сервиса

[tracing](week_04/tracing.go): отслеживание хода выполнения программы за n секунд, трассировка вызовов. Прямо в проде под нагрузкой (хм...).

- `curl http://localhost:8080/debug/pprof/trace?seconds=10 -o trace.out` снятие трассы.
- `go tool trace -http "0.0.0.0:8081" tracing.exe trace.out` анализ трассы (в браузере).

### Пример с telegram-ботом

Создание телеграм-бота, пример.
Бот это веб-сервис с определенным API, он общается с телеграм сервером также по определенному API.

- [bot](week_04/bot.go) принимает сообщения от телеграм (в канал), берет сообщение из канала,
    выполняет некую работу (читает rss с хабра) и отправляет ответ на сообщение.

- `BotFather`: зарегить в телеграм.
- `ngrok.io` как прокси для вывода бота в интернет. Можете заюзать heroku или другие облачные решения для запуска сервисов.

## links, info

- [Package `fmt` implements formatted I/O with functions analogous to C's printf and scanf](https://pkg.go.dev/fmt)

Материалы для дополнительного чтения на английском:
* https://go.dev/ref/spec - спецификация по языку
* https://go.dev/ref/mem - модель памяти го. на начальном этапе не надо, но знать полезно
* https://go.dev/doc/code - How to Write Go Code
* https://pkg.go.dev/cmd/go - Пакеты, cmd/go
* https://go.dev/blog/strings - Strings, bytes, runes and characters in Go
* https://go.dev/blog/slices - Arrays, slices (and strings): The mechanics of 'append'
* https://go.dev/blog/slices-intro - Go Slices: usage and internals
* https://github.com/golang/go/wiki - вики го на гитхабе. очень много полезной информации
* https://go.dev/blog/maps - Go maps in action
* https://go.dev/blog/organizing-go-code - Organizing Go code
* https://go.dev/doc/effective_go - основной сборник тайного знания, сюда вы будуте обращатсья в первое время часто
* https://github.com/golang/go/wiki/CodeReviewComments - как ревьювить (и писать код). обязательно к прочтению
* https://divan.dev/posts/avoid_gotchas/ - материал аналогичный 50 оттенков го
* https://research.swtch.com/interfaces - Go Data Structures: Interfaces (Russ Cox)
* https://research.swtch.com/godata - Go Data Structures (Russ Cox)
* https://jordanorelli.com/post/42369331748/function-types-in-go-golang - Function Types in Go (golang)
* https://www.devdungeon.com/content/working-files-go - работа с файлами
* https://www.golangprograms.com - много how-to касательно базовых вещей в go
* https://yourbasic.org/golang/ - ещё большой набор how-to где можно получить углублённую информацию по всем базовым вещам. Очень полезны https://yourbasic.org/golang/blueprint/
* https://go101.org/article/101.html - e-book, похожий на предыдущий сайт с кучей информации по основам и основным местам
* https://github.com/Workiva/go-datastructures - A collection of useful, performant, and threadsafe Go datastructures
* https://github.com/enocom/gopher-reading-list - большая подборка статей по многим темам ( не только данной лекции )
* https://youtu.be/MzTcsI6tn-0 - как организовать код / Ashley McNamara + Brian Ketelsen. Go best practices
* https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1 - статья на предыдущую тему / Standard Package Layout
* https://dave.cheney.net/practical-go/presentations/qcon-china.html - Practical Go: Real world advice for writing maintainable Go programs / Dave Cheney dave@cheney.net Version 12c316-Dirty, 2019-04-24

На русском:
* https://habrahabr.ru/post/141853/ - как работают горутины
* https://habrahabr.ru/post/308070/ - как работают каналы
* https://habrahabr.ru/post/333654/ - как работает планировщик ( https://rakyll.org/scheduler/ )
* https://habrahabr.ru/post/271789/ - танцы с мютексами
* https://habr.com/ru/company/avito/blog/466495/ - как не ошибиться с конкурентностью в go

На английском:
* https://blog.golang.org/race-detector
* https://blog.golang.org/pipelines
* https://blog.golang.org/advanced-go-concurrency-patterns
* https://blog.golang.org/go-concurrency-patterns-timing-out-and
* https://talks.golang.org/2012/concurrency.slide#1
* https://www.goinggo.net/2017/10/the-behavior-of-channels.html
* http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/ - рассказ про оптимизацию воркер пула
* http://www.tapirgames.com/blog/golang-channel
* http://www.tapirgames.com/blog/golang-channel-closing
* https://github.com/golang/go/wiki/CommonMistakes

Видео:
* https://www.youtube.com/watch?v=5buaPyJ0XeQ - классное выступление Dave Cheney про функции первого класса и использование их с горутинами, очень рекомендую, оно небольшое
* https://www.youtube.com/watch?v=f6kdp27TYZs - Google I/O 2012 - Go Concurrency Patterns - очень рекомендую
* https://www.youtube.com/watch?v=rDRa23k70CU&list=PLDWZ5uzn69eyM81omhIZLzvRhTOXvpeX9&index=15 - ещё одно хорошее видео про паттерны конкуренции в го
* https://www.youtube.com/watch?v=KAWeC9evbGM - видео Андрея Смирнова с конференции Highload - в нём вы можете получить более детальную информацию по теме вводного видео (методы обработки запросов и плюсы неблокирующего подхода), о том, что там творится на системном уровне. На русском, не про go

linter tools
* https://pkg.go.dev/cmd/vet
* https://golangci-lint.run/usage/configuration
* https://github.com/golang/vscode-go/wiki/settings#uidiagnosticanalyses
* https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/stringintconv#hdr-Analyzer_stringintconv
