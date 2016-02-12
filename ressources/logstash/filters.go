package logstash

import (
	"bytes"
	"strconv"
	"text/template"
)

const (
	FILTERS_TEMPLATE string = "filters_template"
)

var (
	filtersTemplate *template.Template
)

func init() {
	var err error
	if filtersTemplate, err = template.New(FILTERS_TEMPLATE).Parse(`filter {
    {{ range .Filters }}
    if {{ .Cond }} {
        mutate {
            {{ range .Mut }}
                add_field => {
                {{ range $key, $value := .AddField }}
                    {{ $key }} => {{ $value }}
                {{ end }}
                }
            {{ end }}
        }
    }
    {{ end }}
}`); err != nil {
		panic(err)
	}
}

type FilterMutate struct {
	AddField map[string]string
}

func NewMutate(assocs ...string) *FilterMutate {

	fm := &FilterMutate{
		AddField: make(map[string]string),
	}

	for i := 0; i < len(assocs); i += 2 {
		fm.AddField[strconv.Quote(assocs[i])] = assocs[i+1]
	}

	return fm
}

type Filter struct {
	Cond string
	Mut  []*FilterMutate
}

func (f *Filter) Mutate(m *FilterMutate) *Filter {
	f.Mut = append(f.Mut, m)
	return f
}

func (f *Filter) Condition(cond string) *Filter {
	f.Cond = cond
	return f
}

func NewFilter(cond string) *Filter {
	return &Filter{
		Cond: cond,
	}
}

type TagsFilters struct {
	Filters []*Filter
}

func NewTagFilters() *TagsFilters {
	return new(TagsFilters)
}

func (tf *TagsFilters) AddFilter(f *Filter) {
	tf.Filters = append(tf.Filters, f)
}

func (tf *TagsFilters) Marshal() ([]byte, error) {

	buf := &bytes.Buffer{}

	if err := filtersTemplate.Execute(buf, tf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
