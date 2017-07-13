package cloudinary

import (
	"bytes"
	"sort"
	"strconv"
	"time"
)

// TODO: distinguish common fields into base struct to follow DRY and not repeat ToParams for them.

// ExplicitParams is a basic set of fields needed for "explicit" api method.
type ExplicitParams struct {
	PublicID   string
	Type       string
	Eager      Transformations
	Invalidate bool
}

func (p *ExplicitParams) ToParams() params {
	var params params
	params.set("public_id", p.PublicID)
	params.set("type", p.Type)
	params.set("eager", p.Eager.String())
	params.set("invalidate", strconv.FormatBool(p.Invalidate))
	return params
}

type UploadParams struct {
	Transformation *Transformation
	PublicID       string
	Invalidate     bool
}

func (p *UploadParams) ToParams() params {
	var params params
	params.set("transformation", p.Transformation.String())
	params.set("public_id", p.PublicID)
	params.set("invalidate", strconv.FormatBool(p.Invalidate))
	return params
}

// params is a key-value store of data which is then passed to api request as body or
// is used to make a signature.
type params map[string]string

// stringForSignature works the same as url.Values.Encode, but don't apply Escaping and ignore
// "api_key" entry in the underlying p.
func (p params) stringForSignature() string {
	var buf bytes.Buffer
	keys := make([]string, 0, len(p))
	for pName := range p {
		if pName == "api_key" {
			continue
		}
		keys = append(keys, pName)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := p[k]
		if v == "" {
			continue
		}
		prefix := k + "="
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(prefix)
		buf.WriteString(v)
	}
	return buf.String()
}

func (p *params) guaranteeParams() {
	if p == nil || *p == nil {
		*p = make(params)
	}
}

func (p *params) set(key, value string) {
	p.guaranteeParams()
	(*p)[key] = value
}

func (p *params) join(p2 params) {
	p.guaranteeParams()

	for key, value := range p2 {
		(*p)[key] = value
	}
}

func (p *params) publicID() (string, bool) {
	p.guaranteeParams()
	publicID, ok := (*p)["public_id"]
	return publicID, ok
}

func (s *Service) basicParams() params {
	return params{
		"api_key":   s.apiKey,
		"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
	}
}