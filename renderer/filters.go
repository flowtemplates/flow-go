package renderer

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/flowtemplates/flow-go/value"
	"github.com/iancoleman/strcase"
)

type filter func(value.Valuable) value.Valuable

var filtersMap = map[string]filter{
	"upper": func(v value.Valuable) value.Valuable {
		return value.StringValue(strings.ToUpper(v.AsString()))
	},
	"lower": func(v value.Valuable) value.Valuable {
		return value.StringValue(strings.ToLower(v.AsString()))
	},
	"pascal": func(v value.Valuable) value.Valuable {
		return value.StringValue(strcase.ToCamel(v.AsString()))
	},
	"camel": func(v value.Valuable) value.Valuable {
		return value.StringValue(strcase.ToLowerCamel(v.AsString()))
	},
	"kebab": func(v value.Valuable) value.Valuable {
		return value.StringValue(strcase.ToKebab(v.AsString()))
	},
	"snake": func(v value.Valuable) value.Valuable {
		return value.StringValue(strcase.ToSnake(v.AsString()))
	},
	"capitalize": func(v value.Valuable) value.Valuable {
		s := v.AsString()

		return value.StringValue(string(unicode.ToUpper(rune(s[0]))) + s[1:])
	},
	"title": func(v value.Valuable) value.Valuable {
		var sb strings.Builder
		prevSpace := true
		for _, c := range v.AsString() {
			if unicode.IsSpace(c) {
				prevSpace = true
			} else if prevSpace {
				prevSpace = false
				c = unicode.ToUpper(c)
			}

			sb.WriteRune(c)
		}

		return value.StringValue(sb.String())
	},
	"length": func(v value.Valuable) value.Valuable {
		return value.NumberValue(len(v.AsString()))
	},
	"trim": func(v value.Valuable) value.Valuable {
		return value.StringValue(strings.TrimSpace(v.AsString()))
	},
}

func callFilter(name string, v value.Valuable) (value.Valuable, error) {
	f, ok := filtersMap[name]
	if !ok {
		return nil, fmt.Errorf("filter %s is not declared", name)
	}

	return f(v), nil
}
