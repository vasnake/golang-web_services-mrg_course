package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	var defaultStr string               // empty by default
	var interprStr = "CR: \n, TAB: \t." // will be parsed
	var rawStr = `CR: \n, TAB: \t.`     // as is
	var utfStr = "ЯЫЧфйДжЮËъЭщП"        // all strings are utf8 (really?)
	var oneByte byte = '\x27'           // single quotes for symbols, uint8 in fact = 39
	var oneRune rune = 'Ы'              // int32 in fact = 1067
	fmt.Println("strings: ", defaultStr, interprStr, rawStr, utfStr, oneRune, oneByte)

	// concat
	fmt.Println("concat: ", rawStr+" wha?")

	// immutable
	//rawStr[0] = 39 // prohibited

	// string length
	byteLen := len(utfStr)
	numSymbols := utf8.RuneCountInString(utfStr)
	fmt.Println("string len: ", byteLen, numSymbols) // 26 13

	// slice, view
	fmt.Println("byte slice: ", utfStr[:3], utfStr[:4], rawStr[1]) // Я? ЯЫ 82
	bytes := []byte(utfStr)                                        //  [208 175 208 171 208 167 209 132 208 185 208 148 208 182 208 174 195 139 209 138 208 173 209 137 208 159]
	symbols := string(bytes[24:])                                  // П
	fmt.Println("bytes<->string: ", bytes, symbols)

}

/*
	// пустая строка по-умолчанию
	var str string

	// со спец символами
	var hello string = "Привет\n\t"

	// без спец символов
	var world string = `Мир\n\t`

	fmt.Println("str", str)
	fmt.Println("hello", hello)
	fmt.Println("world", world)

	// UTF-8 из коробки
	var helloWorld = "Привет, Мир!"
	hi := "你好，世界"

	fmt.Println("helloWorld", helloWorld)
	fmt.Println("hi", hi)

	// одинарные кавычки для байт (uint8)
	var rawBinary byte = '\x27'

	// rune (uint32) для UTF-8 символов
	var someChinese rune = '茶'

	fmt.Println(rawBinary, someChinese)

	helloWorld = "Привет Мир"
	// конкатенация строк
	andGoodMorning := helloWorld + " и доброе утро!"

	fmt.Println(helloWorld, andGoodMorning)

	// строки неизменяемы
	// cannot assign to helloWorld[0]
	// helloWorld[0] = 72

	// получение длины строки
	byteLen := len(helloWorld)                    // 19 байт
	symbols := utf8.RuneCountInString(helloWorld) // 10 рун

	fmt.Println(byteLen, symbols)

	// получение подстроки, в байтах, не символах!
	hello = helloWorld[:12] // Привет, 0-11 байты
	H := helloWorld[0]      // byte, 72, не "П"
	fmt.Println(H)

	// конвертация в слайс байт и обратно
	byteString := []byte(helloWorld)
	helloWorld = string(byteString)

	fmt.Println(byteString, helloWorld)
}

*/
