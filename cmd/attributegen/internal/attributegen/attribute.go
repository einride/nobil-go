package attributegen

import (
	"strings"

	"github.com/stoewer/go-strcase"
)

type Attribute struct {
	ID             int
	Name           string
	TranslationID  int
	UserSelectable bool
	Required       bool
	Values         []*AttributeValue
}

func (a *Attribute) GoIdent() string {
	return strcase.UpperCamelCase(a.Name)
}

type AttributeValue struct {
	ID         int
	Name       string
	Translated string
	Key        string
}

func (a *AttributeValue) GoIdent() string {
	if a.Key != "" {
		return strcase.UpperCamelCase(a.Key)
	}
	return strcase.UpperCamelCase(
		strings.NewReplacer(
			",", "Point",
			"+", "And",
			"(", "_",
			")", "_",
			"/", "_",
		).Replace(a.Name),
	)
}
