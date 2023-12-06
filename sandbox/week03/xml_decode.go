package main

import (
	"bytes"
	"encoding/xml"
	"io"
)

type UserV3 struct {
	ID      int    `xml:"id,attr"`
	Login   string `xml:"login"`
	Name    string `xml:"name"`
	Browser string `xml:"browser"`
}

type Users struct {
	Version string `xml:"version,attr"`
	List    []User `xml:"user"`
}

func findAllLoginsInDoc() []string {
	logins := make([]string, 0, 64)
	docRef := new(Users)

	err := xml.Unmarshal(xmlDocBytes, &docRef) // all at once, memory pressing
	if err != nil {
		show("findAllLoginsInDoc, error: ", err)

	} else {
		for _, u := range docRef.List {
			logins = append(logins, u.Login)
		}
	}

	return logins
}

func findAllLoginsInStream() []string {
	logins := make([]string, 64)
	bytesReaderRef := bytes.NewReader(xmlDocBytes)
	decoderRef := xml.NewDecoder(bytesReaderRef)
	var login string

	for {
		token, err := decoderRef.Token()
		if err != nil {
			if err != io.EOF {
				show("findAllLoginsInStream, shit happend: ", err)
			}
			break
		}

		if token == nil {
			show("findAllLoginsInStream, token is nil, bad XML?")
			break
		}

		switch tok := token.(type) {
		case xml.StartElement:
			if tok.Name.Local == "login" {
				if err := decoderRef.DecodeElement(&login, &tok); err != nil {
					show("findAllLoginsInStream, shit happend", err)
				}
				logins = append(logins, login)
			}
		}
	}

	return logins
}

/*
	go test -bench . -benchmem xml_test.go
*/

var xmlDocBytes = []byte(`<?xml version="1.0" encoding="utf-8"?>
	<users>
		<user id="1">
			<login>user1</login>
			<name>Василий Романов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36</browser>
		</user>
		<user id="2">
			<login>user2</login>
			<name>Иван Иванов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36</browser>
		</user>
		<user id="2">
			<login>user3</login>
			<name>Иван Петров</name>
			<browser>Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; Trident/5.0)</browser>
		</user>
		<user id="1">
			<login>user1</login>
			<name>Василий Романов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36</browser>
		</user>
		<user id="2">
			<login>user2</login>
			<name>Иван Иванов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36</browser>
		</user>
		<user id="2">
			<login>user3</login>
			<name>Иван Петров</name>
			<browser>Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; Trident/5.0)</browser>
		</user>
		<user id="2">
			<login>user3</login>
			<name>Иван Петров</name>
			<browser>Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; Trident/5.0)</browser>
		</user>
		<user id="1">
			<login>user1</login>
			<name>Василий Романов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36</browser>
		</user>
		<user id="2">
			<login>user2</login>
			<name>Иван Иванов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36</browser>
		</user>
		<user id="2">
			<login>user3</login>
			<name>Иван Петров</name>
			<browser>Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; Trident/5.0)</browser>
		</user>
		<user id="2">
			<login>user3</login>
			<name>Иван Петров</name>
			<browser>Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; Trident/5.0)</browser>
		</user>
		<user id="1">
			<login>user1</login>
			<name>Василий Романов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36</browser>
		</user>
		<user id="2">
			<login>user2</login>
			<name>Иван Иванов</name>
			<browser>Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36</browser>
		</user>
		<user id="2">
			<login>user3</login>
			<name>Иван Петров</name>
			<browser>Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; Trident/5.0)</browser>
		</user>
	</users>`)
