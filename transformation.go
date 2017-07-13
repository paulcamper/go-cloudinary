package cloudinary

import (
	"bytes"
	"strconv"
)

type Transformation struct {
	Quality  Quality
	Width    float32 // if a value is <1, than the result_width=original_width*Width. Otherwise the number of pixels is assumed.
	Height   float32 // if a value is <1, than the result_width=original_width*Width. Otherwise the number of pixels is assumed.
	CropMode CropMode
}

func (t *Transformation) String() string {
	if t == nil {
		return ""
	}

	var buf bytes.Buffer
	if t.Quality != "" {
		buf.WriteString(string(t.Quality))
	}

	needChainDelimiter := t.Quality != "" // cloudinary requires to put "/" (without quotes) if you need chaining.
	needComma := false

	if t.Width != 0 {
		if needChainDelimiter {
			buf.WriteString("/")
			needChainDelimiter = false
		}
		buf.WriteString("w_" + strconv.FormatFloat(float64(t.Width), 'f', -1, 64))
		needComma = true
	}
	if t.Height != 0 {
		if needChainDelimiter {
			buf.WriteString("/")
			needChainDelimiter = false
		}
		if needComma {
			buf.WriteString(",")
		}
		buf.WriteString("h_" + strconv.FormatFloat(float64(t.Height), 'f', -1, 64))
		needComma = true // to make it obvious from the first sight, although it seems like a redundant stuff.
	}
	if t.CropMode != "" {
		if needComma {
			buf.WriteString(",")
		}
		buf.WriteString(string(t.CropMode))
		needChainDelimiter = true
	}

	return buf.String()
}

type Transformations []*Transformation

func (t *Transformations) String() string {
	var buf bytes.Buffer
	for _, transformation := range *t {
		if buf.Len() > 0 {
			buf.WriteString("/")
		}
		buf.WriteString(transformation.String())
	}
	return buf.String()
}
