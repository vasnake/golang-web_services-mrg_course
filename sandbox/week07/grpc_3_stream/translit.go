package grpc_3_stream

import (
	"fmt"
	// "gws/7/microservices/grpc_stream/translit"
	"io"

	tr "github.com/gen1us2k/go-translit"
)

type TrServer struct {
}

func NewTr() *TrServer {
	return &TrServer{}
}

func (srv *TrServer) EnRu(inStream Transliteration_EnRuServer) error {
	for {
		inWord, err := inStream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		out := &Word{
			Word: tr.Translit(inWord.Word),
		}
		fmt.Println(inWord.Word, "->", out.Word)

		inStream.Send(out)
	}
}

// mustEmbedUnimplementedTransliterationServer implements TransliterationServer.
func (srv *TrServer) mustEmbedUnimplementedTransliterationServer() {
	panic("unimplemented")
}
