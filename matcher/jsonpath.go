package matcher

import "github.com/vitorsalgado/mocha/internal/jsonpath"

func JSONPath[V any](p string, m Matcher[V]) Matcher[any] {
	return func(v any, params Args) (bool, error) {
		value, err := jsonpath.Get(p, v)
		if err != nil || value == nil {
			return false, err
		}

		return m(value.(V), params)
	}
}
