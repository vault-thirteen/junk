// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package request

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonF4c71a47DecodeTestFPRCPkgModelsHttpRequest(in *jlexer.Lexer, out *UserLogInRequest) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "user":
			easyjsonF4c71a47DecodeTestFPRCPkgModelsHttpRequest1(in, &out.User)
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
func easyjsonF4c71a47EncodeTestFPRCPkgModelsHttpRequest(out *jwriter.Writer, in UserLogInRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"user\":"
		out.RawString(prefix[1:])
		easyjsonF4c71a47EncodeTestFPRCPkgModelsHttpRequest1(out, in.User)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserLogInRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF4c71a47EncodeTestFPRCPkgModelsHttpRequest(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserLogInRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF4c71a47EncodeTestFPRCPkgModelsHttpRequest(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserLogInRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF4c71a47DecodeTestFPRCPkgModelsHttpRequest(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserLogInRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF4c71a47DecodeTestFPRCPkgModelsHttpRequest(l, v)
}
func easyjsonF4c71a47DecodeTestFPRCPkgModelsHttpRequest1(in *jlexer.Lexer, out *UserLogInRequestUser) {
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
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "internal_name":
			out.InternalName = string(in.String())
		case "password":
			out.Password = string(in.String())
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
func easyjsonF4c71a47EncodeTestFPRCPkgModelsHttpRequest1(out *jwriter.Writer, in UserLogInRequestUser) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"internal_name\":"
		out.RawString(prefix[1:])
		out.String(string(in.InternalName))
	}
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	out.RawByte('}')
}