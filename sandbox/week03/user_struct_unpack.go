package main

import (
	"bytes"
	"encoding/binary"
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// Code generated by our little toy example

func (in *User) Unpack(data []byte) error {
	r := bytes.NewReader(data)

	// ID
	var IDRaw uint32
	binary.Read(r, binary.LittleEndian, &IDRaw)
	in.ID = int(IDRaw)

	// Login
	var LoginLenRaw uint32
	binary.Read(r, binary.LittleEndian, &LoginLenRaw)
	LoginRaw := make([]byte, LoginLenRaw)
	binary.Read(r, binary.LittleEndian, &LoginRaw)
	in.Login = string(LoginRaw)

	// Flags
	var FlagsRaw uint32
	binary.Read(r, binary.LittleEndian, &FlagsRaw)
	in.Flags = int(FlagsRaw)
	return nil
}

// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson9f2eff5fDecodeSt(in *jlexer.Lexer, out *UserV2) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Id":
			out.Id = int(in.Int())
		case "RealName":
			out.RealName = string(in.String())
		case "Login":
			out.Login = string(in.String())
		case "Flags":
			out.Flags = int(in.Int())
		case "Status":
			out.Status = int(in.Int())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func easyjson9f2eff5fEncodeSt(out *jwriter.Writer, in UserV2) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"Id\":")
	out.Int(int(in.Id))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"RealName\":")
	out.String(string(in.RealName))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"Login\":")
	out.String(string(in.Login))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"Flags\":")
	out.Int(int(in.Flags))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"Status\":")
	out.Int(int(in.Status))
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserV2) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9f2eff5fEncodeSt(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserV2) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9f2eff5fEncodeSt(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserV2) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9f2eff5fDecodeSt(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserV2) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9f2eff5fDecodeSt(l, v)
}
