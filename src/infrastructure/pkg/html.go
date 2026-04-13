package pkg

import (
	"fmt"

	rendererDomain "everlasting/src/domain/sharedkernel/renderer"

	"github.com/flosch/pongo2/v6"
)

type HtmlRenderer map[rendererDomain.Template]*pongo2.Template

func NewHtmlRenderer(templateDir string) HtmlRenderer {
	result := make(HtmlRenderer)
	for _, template := range rendererDomain.Templates {
		result[template] = pongo2.Must(pongo2.FromFile(fmt.Sprintf("%s%s", templateDir, string(template))))
	}
	return result
}

func (h HtmlRenderer) Render(template rendererDomain.Template, values map[string]any) (result string, err error) {
	return h[template].Execute(values)
}
