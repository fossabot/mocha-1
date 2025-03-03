package mocha

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/vitorsalgado/mocha/internal/header"
	"github.com/vitorsalgado/mocha/internal/mime"
)

type (
	BodyParser interface {
		CanParse(content string, r *http.Request) bool
		Parse(r *http.Request) (any, error)
	}
)

func ParseRequestBody(r *http.Request, parsers []BodyParser) (any, error) {
	if r.Body != nil && r.Method != http.MethodGet && r.Method != http.MethodHead {
		var content = r.Header.Get(header.ContentType)

		for _, parse := range parsers {
			if parse.CanParse(content, r) {
				body, err := parse.Parse(r)
				if err != nil {
					return nil, err
				}

				return body, nil
			}
		}
	}

	return nil, nil
}

type JSONBodyParser struct{}

func (parser JSONBodyParser) CanParse(content string, _ *http.Request) bool {
	return strings.Contains(content, mime.ContentTypeJSON)
}

func (parser JSONBodyParser) Parse(r *http.Request) (any, error) {
	var d any
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

type FormURLEncodedParser struct{}

func (parser FormURLEncodedParser) CanParse(content string, _ *http.Request) bool {
	return strings.Contains(content, mime.ContentTypeFormURLEncoded)
}

func (parser *FormURLEncodedParser) Parse(r *http.Request) (any, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	return r.Form.Encode(), nil
}
