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
