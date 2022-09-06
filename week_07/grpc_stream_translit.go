package main

import (
	"coursera/microservices/grpc_stream/translit"
	"fmt"
	"io"

	tr "github.com/gen1us2k/go-translit"
)

// service class

type TrServer struct {
}

func NewTr() *TrServer {
	return &TrServer{}
}

// service implementation, NB one stream for in and out
func (srv *TrServer) EnRu(inStream translit.Transliteration_EnRuServer) error {
	for {
		inWord, err := inStream.Recv()
		if err == io.EOF { // closed stream
			return nil
		}
		if err != nil {
			return err
		}

		out := &translit.Word{
			Word: tr.Translit(inWord.Word),
		}

		fmt.Println(inWord.Word, "->", out.Word)
		inStream.Send(out)
	}
	return nil
}
