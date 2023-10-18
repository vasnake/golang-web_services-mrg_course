# Разработка веб-сервисов на Go. Часть 1

[Go course, MRG, Романов Василий](README.md)
[Разработка веб-сервисов на Go, Часть 1](https://study.vk.team/learning/play/671)

- week 1, Введение в Go
- week 2, Асинхронная работа
- week 3, Работа с динамическими данными и производительность
- week 4, Основы HTTP

## part 1, week 1

Введение в Go. [Код, домашки, литература](week_01/materials.zip)

Зачем нужен еще один язык программирования?
Go-team (backend development) хотела язык C/C++ но без их недостатков, плюс эффективная утилизация многопроцессорных систем.
Эффективность (работы программеров в Гугл): компиляции, выполнения, разработки. Зависимости, рантайм, garbage-collection.
Утилизация многопроцессорных систем (легкие потоки, асинхронность, CSP);
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
- Аппа состоит из модулей, модули содержат пакеты, пакеты содержат файлы `.go`, [details](https://go.dev/doc/modules/layout)
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

[snb/hello_world](./sandbox/hello_world/main.go):
В общем, мои грабли оказались таковы: в vscode у меня воркспейс был из одной рут-директории `desktop`,
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

- https://play.golang.com/
- [run](run.sh)
- [w01/hello_world](week_01/hello_world.go)

Используй `camelCase` для имён. Публичные (экспортируемые) обьекты называй с большой буквы: `Println`.

Файл с кодом содержит:
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
[spec playground](./sandbox/spec/main.go)

Пакет: имя, нэймспейс для группировки: констант, типов, переменных, функций. Один или несколько файлов.
Идентификатор может быть "экспортирован" из пакета, если имя начинается с Большой Буквы.
Один пакет = одна директория (по имени пакета).

Модуль: коллекция пакетов, сопровождаемая файлом `go.mod`.

> general-purpose language designed with systems programming in mind. It is strongly typed and garbage-collected and has explicit support for concurrent programming.

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

Строки: последовательность символов (рун), Unicode code-points. Длина строки в рунах != длине строки в байтах.
Если надо работать с коллекцией байт или рун, то надо явно привести тип (`[]rune(str)` or `[]byte(str)`).
Такое приведение типа создает копию данных, где иммутабельность уже не соблюдается.
> A string value is a sequence of bytes. The number of bytes is called the length of the string ...
Strings are immutable
It is illegal to take the address of a string's byte

> Numeric constants represent exact values of arbitrary precision and do not overflow.
Consequently, there are no constants denoting the IEEE-754 negative zero, infinity, and not-a-number values.
Constants may be typed or untyped.
An untyped constant has a default type

> Numeric types: uint8..64, int8..64, float32..64, complex64..128, byte (alias uint8), rune (alias uint32).
Predeclared integer types with implementation-specific sizes: uint, int, uintptr.

Array: коллекция элементов одного типа, массивы разных размеров это разные типы.
Определение типа массива не может быть рекурсивным.
> The length is part of the array's type
An array type T may not have an element of type T, or of a type containing T as a component

Slice: можно сказать, что слайс это view на array с данными.
У слайса есть переменные длина и капасити.
Слайс разделяет "хранилище" (массив) с другими слайсами, они все смотрят на один и тот-же кусок памяти.

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
Канал может быть буферизован, это не отражается на его типе. Небуферизованный канал работает только если приемщик и отправитель оба готовы.
Т.е. операция с каналом может быть блокирующей, если операция не может быть поддержана буфером!
Канал на отправку можно (нужно) закрывать вызовом `close`, что маркирует его как "отправок более не будет".
> A channel may be constrained only to send or only to receive by assignment or explicit conversion.
... The `<-` operator associates with the leftmost chan possible

В целом, каналы обладают сложным протоколом: опциональная буферизация, закрытие канала, (не)блокирующие операции, паника и zero-значения ... WTF!

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
но это не так для Constant Expressions!

Floating-point operators:
> An implementation may combine multiple floating-point operations into a single fused operation

Comparison operators:
> Two string values are compared lexically byte-wise
Two channel values are equal if they were created by the same call to `make`
Slice, map, and function types are not comparable
Boolean, numeric, string, pointer, and channel types are strictly comparable

Receive operator:
> Receiving from a nil channel blocks forever

Conversion:
> nil is not a constant ...
Converting a constant to a type (that is not a type parameter) yields a typed constant ...
conversions only change the type but not the representation of x (not applicable to `numeric <=> string`) ...
there is no indication of overflow (int) ...
the conversion succeeds (float) but the result value is implementation-dependent ...

> There is no linguistic mechanism to convert between pointers and integers. 
The package `unsafe` implements this functionality under restricted circumstances

Исключение из правил (`string(rune(x))` vs `strconv.Itoa`):
> an integer value may be converted to a string type

> converting a slice to an array pointer yields a pointer to the underlying array of the slice

Constant expressions:
тип результата бинарной операции (над нетипизированными операндами) определяет тип самого правого операнда.
Что противоположно правилам Arithmetic Operators!

Order of evaluation:
порядок вычислений операндов определён не во всех случаях
> At package level, initialization dependencies determine the evaluation order ... but not for operands within each expression ...
all function calls, method calls, and communication operations are evaluated in lexical left-to-right order ...

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

Зависимости складываются в директории vendor.

### Основы работы с интерфейсами

- [basic](week_01/basic.go)
- [many](week_01/many.go)
- [cast](week_01/cast.go)

- Интерфейс это тип `type Payer interface { Pay(int) error }`
- Другой тип может содержать реализацию интерфейса `func (w *Wallet) Pay(amount int) error { ... }`
- При вызове похер на тип, указываем интерфейс `func Buy(p Payer) {	err := p.Pay(10) }`

Можно (нужно?) держать переменную (поле структуры) с типом нужного интерфейса `var p Payer; p = &Card{Balance: 100}`

Можно кастовать тип от интерфейса обратно к конкретному типу, предварительно сматчив тип.

Так реализуется полиморфизм. Добавил реализацию интерфейса в произвольный тип и радуйся.

Неясно, как, глядя на код, сразу сказать, кто реализует какой интерфейс и в каком объеме.

### Пустой интерфейс

- [empty_1](week_01/empty_1.go)
- [empty_2](week_01/empty_2.go)

Пустой интерфейс (type, value) не накладывает ограничений. Использовать его можно только если ручками делать проверку/приведение типа:
`func Buy(in interface{}) { if p, ok = in.(Payer) ...}`

Рассказал на примере `fmt.Printf`, как аргумент типа "пустой интерфейс" может быть подерган за разные методы (опции форматирования определяют вызываемый метод).

### Композиция интерфейсов

- [embed](week_01/embed.go)

Как и структуры, интерфесы могут быть составлены из других интерфейсов.
При этом, как и структуры, охватывающий интерфейс включает методы вложенных интерфейсов.

### Написание Программы Уникализации (ПУ)

- [uniq](week_01/uniq.go)
- [data_map.txt](week_01/data_map.txt)

Пример программы: получает на вход файл и выводит только уникальные строки из этого файла.

Два варианта: на мапке строка - уникальность; на предположении, что файл сортирован (prev == current).

### Написание тестов для ПУ

- [unique/unique](week_01/unique/unique.go)
- [unique/unique_test](week_01/unique/unique_test.go)

Чтобы поддерживать тестирование, зависимости передаются в функцию как аргументы
`func sortedInputUnique(input io.Reader, output io.Writer) error { ... }`

Модуль тестов должен быть файлом с именем с суффиксом `_test.go`.
Должен содержать функции-тесты с именами с префиксом `Test`:
`func TestSortedInput(t *testing.T) { ... }`

Тесты запускаются `go test -v ./unique`

## part 1, week 2

Асинхронная работа.
[Код, домашки, литература](week_02/w2_materials.zip) https://cloud.mail.ru/public/YDEX/Dau2wVWuw/

### Методы обработки запросов и плюсы неблокирующего подхода

Асинхронное выполнение (AJAX),
скорость потока данных процессор-кеш-память, время на переключение контекста (выгрузка-загрузка регистров),
современные тенденции на многоядерность и параллельность,
тяжелые процессы, потоки легче, асинхронные сопрограммы (green threads) еще легче.

Утилизация дорогого железа -- сервер должен работать.

Невытесняющая многозадачность (eventloop, Windows 3.0) vs preemptive.

Ввод-вывод и ожидание возврата из syscall. Время ожидания можно потратить на другие задачи, non-blocking IO.
IO-bound vs CPU-bound.

`Communicating Sequential Processes` by Tony Hoare. Горутины перемещаются между системными потоками,
код горутины может быть выполнен в любом потоке (как перемещается стек?).

### Горутины -- легковесные процессы

- [goroutines](week_02/goroutines.go)

`go doSomeWork(i)` функция не может вернуть значение обычным способом (см. каналы).

`runtime.Gosched()` уйти в планировщик, дав возможность запустить другие горутины, yield.

Есть шанс заблокировать шедулер, если молотить цикл без вызовов системных функций.

### Каналы -- передаём данные между горутинами

- [chan_1](week_02/chan_1.go)
- [chan_2](week_02/chan_2.go)

`chan` keyword. Передача контроля над данными между потоками/горутинами.

`ch1 := make(chan int)` небуферизованный канал.
`ch1 <- 42` запись в канал, читатель уже должен ждать, ибо небуферизовано.
`v := <-in` чтение из канала.

Чтение, запись -- операторы стрелочка `x <- myChannel; myChannel <- x`.

Небуферизованные каналы vs буферизованные, размер буфера. Работа с каналом в цикле.

Если писать в небуферизованый канал, из которого никто не читает, будет deadlock, очень плохо.
Наоборот тоже беда, нельзя читать из канала, в который никто не пишет.

Похоже, что система отслеживает закрытие горутины, связанной с каналом и повисший канал вызывает ошибку.

Записанное в буферизованный канал может пропасть без следа, если никто не прочтёт.

Эти проблемы позволяет решить оператор `select` ...

### Мультиплексирование каналов через (не блокирующий) оператор `select`

- [select_1](week_02/select_1.go)
- [select_2](week_02/select_2.go)
- [select_3](week_02/select_3.go)

`select { case val := <-ch1: ...; case ch2 <- 1: ...; default: ... }`
Что-то делает только если канал готов. Если все каналы остановлены, то действие "по умолчанию".

Выбор в цикле, действие по умолчанию если ни один канал не работоспособен -- можно выходить из цикла.

Использование канала "команд" вместе с каналом "данных".

### Таймеры и таймауты (как источник сигнала в каналах)

- [timeout](week_02/timeout.go)
- [tick](week_02/tick.go)
- [afterfunc](week_02/afterfunc.go)

- Таймер как источник сигнала "через определенное время". Реализация тайм-аут логики.
- Тикер как источник регулярных/периодических сигналов в канале.
- AfterFunc как способ отложенного выполнения функции.

Некоторые фунции пакета time могут быть удобны, но приводить к утечкам памяти, создавая вечные таймеры.

### Пакет `context` и отмена выполнения

- [cancel](week_02/context_cancel.go)
- [timeout](week_02/context_timeout.go)

Отмена асинхронных операций, вручную.
Всем воркерам выдается общий контекст `ctx, finish := context.WithCancel(context.Background())`, в котором определена фунция завершения.
Когда получен нужный результат, дергается функция завершения контекста и воркеры понимают (слушая канал `ctx.Done()`), что пора выходить.

Отмена по таймауту `ctx, _ := context.WithTimeout(context.Background(), workTime)`.
В канал `ctx.Done()` приходит сообщение по таймеру, не по ручному вызову функции завершения.

Контекст -- основной способ отмены асинхронных операций, завершения воркеров, выполняющих циклы или долгие операции асинхронно.

### Асинхронное получение данных

- [async_work](week_02/async_work.go)

Практическое применение горутин и каналов: распараллеливание считывания статей и комментариев к ним на сайте.

### Пул воркеров

- [workerpool](week_02/workerpool.go)

Один канал как очередь сообщений, несколько горутин как воркеры читающие из очереди. Главная программа сыпет задания в очередь.

### `sync.Waitgroup` -- ожидание завершения работы

- [waitgroup](week_02/waitgroup.go)

До сих пор использовали ввод с клавиатуры, чтобы программа не завершилась раньше своих воркеров. Теперь так делать не нужно.

Ресурс `wg := &sync.WaitGroup{}`, увеличивается `wg.Add(1)` при добавлении воркеров и уменьшается `defer wg.Done()` при удалении воркеров.
Можно использовать для ожидания `wg.Wait()` завершения всех воркеров.

### Ограничение по ресурсам

- [ratelim](week_02/ratelim.go)

Использование буферизованного канала как ограничителя. Как буфер исчерпался, новая работа не поступает.
Перед стартом работы, воркер пытается записать сообщение в канал квоты. Если в буфере место ещё есть, воркер запишет
сообщение и сможет продолжить работу.
На выходе воркер считывает сообщение из канала квоты, освобождая место в буфере.

### Ситуация гонки на примере конкурентной записи в map

- [race_1](week_02/race_1.go)

Пять воркеров конкурентно пишут в одну мапку. Мапка по ходу операций перестраивается и разные воркеры начинают работать с разными копиями данных.
Кто из них победит? Программа падает с fatal error.

`go run -race ...` для диагностики.

Что же делать?  Ставить блокировки, ...

### `sync.Mutex` для синхронизации данных

- [race_2](week_02/race_2.go)

Берем мютекс `mu := &sync.Mutex{}` и в нужном месте используем его
`mu.Lock(); counters[th*10+j]++; mu.Unlock()`

Но есть нюансы ...

### `sync.Atomic`

- [atomic_1](week_02/atomic_1.go)
- [atomic_2](week_02/atomic_2.go)

Для атомарного изменения одной переменной (регистра ЦП) использовать Mutex слишком дорого.
Есть альтернатива `atomic.AddInt32(&totalOperations, 1)`.

## part 1, week 3

Работа с динамическими данными и производительность
[Код, домашки, литература](week_03/part1_week3.zip) https://cloud.mail.ru/public/2iXh/RC437wn11

### Распаковываем JSON

- [json](week_03/json.go)
- [struct_tags](week_03/struct_tags.go)

- Пакет `encoding/json` дает нам кодек.
Декодирование из слайса байт в структуру (надо знать тип структуры и создать её, пустую, перед декодированием) `json.Unmarshal(bytes, emptyStructRef)`
- Кодирование зеркально, `bytes, err := json.Marshal(someStruct)`.
- Приватные поля (с маленькой буквы которые) не обрабатываются, ибо `encoding/json` пакет не может получить доступ к приватным полям нашего пакета `main`.

Теги структуры -- метаинформация структур. При описании (определении) стурктуры, к полю добавляется строка определенного формата,
где записана метаинформация. В частности, как декодеру json обрабатывать поле, какое имя ему дать, какой тип использовать, etc.

```s
ID       int `json:"user_id,string"`
```

### Нюансы работы с JSON

- [dynamic](week_03/dynamic.go)

Как быть, если мы точно не знаем структуру, представленную json текстом?
Unmarshal в пустой интерфейс.

Также и Marshal, создав мапку `map[string]interface{}{ ... }`, можно ее кодировать в json.

### Пакет reflect -- работаем с динамикой в рантайме

- [reflect_1](week_03/reflect_1.go)
- [reflect_2](week_03/reflect_2.go)

- `reflect.ValueOf(x).Elem()` позволяет итерировать типы и значения из которых состоит "пустой интерфейс".
- Пример, как при помощи рефлексии распаковать слайс байт (бинарные данные на выходе perl `pack`) в структуру.
`reflect.ValueOf(x).Elem()` для получения списка полей структуры, которую надо восстановить из слайса байт.
Перебирая поля структуры, читаем данные из источника байт.

### Кодогенерация -- программа пишет программу

- [unpack](week_03/unpack.go)
- [marshaller](week_03/marshaller.go)
- [codegen](week_03/codegen.go)

К структуре добавлен метод `Unpack`, сгенерированный кодогенератором.
Метод реализует восстановление структуры из бинарных данных, созданных через perl `pack`.
Кодогенератор анализирует гошную структуру (пользуясь гошным AST), считанную из файла с исходником,
и, для всех её полей, создает код восстановления поля из слайса байт.

Кодогенерация полезна, когда некогда в рантайме тратить время на анализ и рефлексию. Профилирование покажет насколько падает скорость ...

### Система бенчмарков Go

- [unpack_test](week_03/unpack_test.go)
- [json_test](week_03/json_test.go)
- [string_test](week_03/string_test.go)
- [prealloc_test](week_03/prealloc_test.go)

Бенчмарки в `testing.B`, `func BenchmarkFooBar( ... ) { ... }`, `go test -bench ...`

Распаковка структуры через кодогенерацию вдвое быстрее рефлексии.

`go test -bench . -benchmem unpack_test.go` замер расхода памяти. Рефлексия жрет вдвое больше памяти.

Пакет кодека json на кодогенерации, EasyJson, в 4 раза эффективнее стандартного.

Демонстрация медленности регулярных выражений на строках. На порядок.

Демонстрация медленности добавления в слайс без преаллокации. В 20 раз медленнее.

### Профилирование через pprof

Собрать профиль
`go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1 unpack_test.go`

Анализ профиля
`go tool pprof main.test.exe mem.out`

Команды: `top`, `list Unpack`, `web`, `alloc_space, top`, `alloc_objects, top`

Можно снимать дампы с работающей программы, не прерывая её работы.

### sync.Pool

- [pool_test](week_03/pool_test.go)

Хотим не выделять память каждый раз, хотим использовать пул преаллоцированной памяти. Для скорости.

`dataPool = sync.Pool{ ... }; data := dataPool.Get().(*bytes.Buffer)`

Процентов 10 скорости выиграли. Операций выделения памяти на порядок меньше. Нагрузка на GC, соответственно, падает драматически.

### Покрытие кода тестами

`go test -v -cover`,
`go test -coverprofile=cover.out`,
`go tool cover -html=cover.out -o cover.html`

### XML

- [main](week_03/xml_main.go)
- [xml_test](week_03/xml_test.go)

Поточная обработка данных, на примере xml. `encoding/xml`.
`xml.NewDecoder(bytes.NewReader(data)).Token()` читаем по токенам, обрабатываем.

## part 1, week 4

Основы HTTP
[Код, домашки, литература](week_04/part1_week4.zip) https://cloud.mail.ru/public/NTFa/barMiVYZd

### Слушаем TCP-сокет с использованием пакета net

- [net_listen](week_04/net_listen.go)

Открытие TCP соединения, получение коннекта в цикле и общение с подключившимся клиентом в отдельной горутине.
`go build -o net_listen.exe . && ./net_listen.exe`
`telnet 127.0.0.1`

Демонстрация многопоточности сервера, несколько клиентов могут общаться с сервером одновременно.

### Обслуживание HTTP-запросов

- [http](week_04/http.go)
- [pages](week_04/pages.go)
- [servehttp](week_04/servehttp.go)
- [mux](week_04/mux.go)
- [servers](week_04/servers.go)

- `net/http` пакет. Простейший вариант с одним хендлером одного роута.
- Обработка трех роутов тремя хендлерами, два из них -- анонимные функции. N.B. набор урлов `/pages/` и один урл `/page`.
- Хендлеры как структуры, реализующие интерфейс. Параметризуемые хендлеры, обладающие внутренним состоянием. Например, коннект к БД.
- Более низкоуровневое API к веб-серверу. Отдельно конфигурируется мультиплексор запросов (обработчики роутов) и отдельно конфигурируется сервер.
    Или несколько серверов.
- Несколько сервисов в рамках одной программы. Сервис + мониторинг, к примеру.

### Работа с параметрами запросов

- [get](week_04/get.go)
- [post](week_04/post.go)
- [cookies](week_04/cookies.go)
- [headers](week_04/headers.go)

- Как считать параметры из объекта риквест, из урла, GET
- `http.MethodPost`, получение данных из формы, POST
- Логин/логаут с использованием cookie (кука как идентификатор сессии пользователя)
- Установка и чтение заголовков, header

### Обслуживание статичных данных

- [static](week_04/static.go)

- `http.FileServer` реализует выдачу статичных файлов (без их интерпретации). В проде не надо, для локальной разработки (пет-проектов) может пригодиться.

### Загрузка файлов формы

- [file_upload](week_04/file_upload.go)

- Два варианта:
    `Request.FormFile` даёт нам загруженный в форму файл.
    `Request.Body` даёт нам байты, загруженные через POST.

### HTTP-запросы во внешние сервисы

- [request](week_04/request.go)

- Три варианта (уровня API) запрос-ответ по HTTP, с точки зрения клиента.

### Тестирование HTTP-запросов и ответов

- [request_test](week_04/request_test.go)
- [server_test](week_04/server_test.go)

- Тестирование функций хендлера с подстановкой моков риквеста и респонса: `httptest.NewRequest`, `httptest.NewRecorder`, `go test ...`.
- Тестирование функции зависящей от внешнего сервиса, мокать внешний сервис.

### Inline-шаблоны и шаблоны из файлов

- [inline](week_04/inline.go)
- [file](week_04/file.go), [users.html](week_04/users.html)

- `text/template` пакет для парсинга и вывода шаблонов (формирование текста ответа сервисом).
- `html/template` шаблонизатор с автоматической экранизацией выводимого html кода.

### Вызов методов и функций из шаблонов

- [method.go](week_04/method.go), [method.html](week_04/method.html)
- [func.go](week_04/func.go), [func.html](week_04/func.html)

- В шаблон передана структура, хотим там (в шаблоне) вызвать её метод (без параметров).
- `template.FuncMap` регистрация фунций, если методов структур нам недостаточно. В шаблоне вызываем регистрированные фунции.

### Профилирование через pprof

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

linter tools
* https://pkg.go.dev/cmd/vet
* https://golangci-lint.run/usage/configuration
* https://github.com/golang/vscode-go/wiki/settings#uidiagnosticanalyses
* https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/stringintconv#hdr-Analyzer_stringintconv
