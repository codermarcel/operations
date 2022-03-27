package src

import (
	"errors"
	"io"
	"os"
	"text/template"
)

type BaseTemplate struct {
	tmp  *template.Template
	out  io.Writer
	fmap template.FuncMap
	name string
}

type templateOption func(*BaseTemplate) error

func TemplateWithFuncMap(fm template.FuncMap) templateOption {
	return func(b *BaseTemplate) error {
		b.fmap = fm
		return nil
	}
}

func TemplateFromFile(path string) templateOption {
	return func(b *BaseTemplate) error {
		t := template.New(path).Funcs(b.fmap)
		t, err := t.ParseFiles(path)

		if err != nil {
			return err
		}

		b.tmp = t
		return nil
	}
}

func TemplateFromString(tmpl string) templateOption {
	return func(b *BaseTemplate) error {
		t := template.New(b.name).Funcs(b.fmap)
		t, err := t.Parse(tmpl)

		if err != nil {
			return err
		}

		b.tmp = t
		return nil
	}
}

func TemplateWithOutput(out io.Writer) templateOption {
	return func(b *BaseTemplate) error {
		b.out = out
		return nil
	}
}
func TemplateWithName(name string) templateOption {
	return func(b *BaseTemplate) error {
		b.name = name
		return nil
	}
}

func NewBaseTemplate(templateOptions ...templateOption) (*BaseTemplate, error) {
	b := &BaseTemplate{out: os.Stdout, fmap: map[string]interface{}{}}
	b.fmap["envOr"] = envOr
	b.name = "basetemplName"

	for _, v := range templateOptions {
		err := v(b)
		if err != nil {
			return nil, err
		}
	}

	if b.tmp == nil {
		return nil, errors.New("template has not been set yet")
	}

	return b, nil
}

func (b *BaseTemplate) Execute(data interface{}) error {
	return b.ExecuteTo(data, b.out)
}

func (b *BaseTemplate) ExecuteTo(data interface{}, out io.Writer) error {
	return b.tmp.Execute(out, data)
}
