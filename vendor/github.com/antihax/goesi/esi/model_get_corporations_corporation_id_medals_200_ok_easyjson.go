// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package esi

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

func easyjsonB6f9bc1aDecodeGithubComAntihaxGoesiEsi(in *jlexer.Lexer, out *GetCorporationsCorporationIdMedals200OkList) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(GetCorporationsCorporationIdMedals200OkList, 0, 1)
			} else {
				*out = GetCorporationsCorporationIdMedals200OkList{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 GetCorporationsCorporationIdMedals200Ok
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonB6f9bc1aEncodeGithubComAntihaxGoesiEsi(out *jwriter.Writer, in GetCorporationsCorporationIdMedals200OkList) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v GetCorporationsCorporationIdMedals200OkList) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonB6f9bc1aEncodeGithubComAntihaxGoesiEsi(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetCorporationsCorporationIdMedals200OkList) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonB6f9bc1aEncodeGithubComAntihaxGoesiEsi(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetCorporationsCorporationIdMedals200OkList) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonB6f9bc1aDecodeGithubComAntihaxGoesiEsi(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetCorporationsCorporationIdMedals200OkList) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonB6f9bc1aDecodeGithubComAntihaxGoesiEsi(l, v)
}
func easyjsonB6f9bc1aDecodeGithubComAntihaxGoesiEsi1(in *jlexer.Lexer, out *GetCorporationsCorporationIdMedals200Ok) {
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
		case "created_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
		case "creator_id":
			out.CreatorId = int32(in.Int32())
		case "description":
			out.Description = string(in.String())
		case "medal_id":
			out.MedalId = int32(in.Int32())
		case "title":
			out.Title = string(in.String())
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
func easyjsonB6f9bc1aEncodeGithubComAntihaxGoesiEsi1(out *jwriter.Writer, in GetCorporationsCorporationIdMedals200Ok) {
	out.RawByte('{')
	first := true
	_ = first
	if true {
		const prefix string = ",\"created_at\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	if in.CreatorId != 0 {
		const prefix string = ",\"creator_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.CreatorId))
	}
	if in.Description != "" {
		const prefix string = ",\"description\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Description))
	}
	if in.MedalId != 0 {
		const prefix string = ",\"medal_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.MedalId))
	}
	if in.Title != "" {
		const prefix string = ",\"title\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Title))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GetCorporationsCorporationIdMedals200Ok) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonB6f9bc1aEncodeGithubComAntihaxGoesiEsi1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GetCorporationsCorporationIdMedals200Ok) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonB6f9bc1aEncodeGithubComAntihaxGoesiEsi1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GetCorporationsCorporationIdMedals200Ok) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonB6f9bc1aDecodeGithubComAntihaxGoesiEsi1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GetCorporationsCorporationIdMedals200Ok) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonB6f9bc1aDecodeGithubComAntihaxGoesiEsi1(l, v)
}