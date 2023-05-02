package renders

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"log"
)

func PreRenderMd(mddata []byte) []byte {
	// add empty line before and after math block '$$'
	for i := 0; i < len(mddata)-1; i++ {
		if mddata[i] == '$' && mddata[i+1] == '$' {
			if i-2 >= 0 && mddata[i-2] != '\n' {
				mddata = append(mddata[:i], append([]byte("\n"), mddata[i:]...)...)
				log.Println("add empty line before math block")
			}
			i += 2
			var j int
			for j = i; j < len(mddata)-3; j++ {
				if mddata[j] == '$' && mddata[j+1] == '$' && mddata[j+3] != '\n' {
					mddata = append(mddata[:j+2], append([]byte("\n"), mddata[j+2:]...)...)
					break
				}
			}
			i = j + 2
		}
	}
	return mddata
}

// render markdown file to html file
func RenderMd(mddata []byte) []byte {
	mddata = PreRenderMd(mddata)
	//log.Println("mddata: %s", string(mddata))

	extensions := parser.CommonExtensions | parser.MathJax

	// Create a new parser with the MathJax extension enabled
	parser := parser.NewWithExtensions(extensions)

	// Create a new renderer
	renderer := html.NewRenderer(html.RendererOptions{})

	//renderer := html.NewRenderer(html.RendererOptions{Flags: html.MathJax})
	// renderer := &customHTMLRenderer{html.NewRenderer(html.RendererOptions{})}

	// Convert Markdown to HTML using the custom renderer
	//markdownText := "This text contains the (c) symbol."
	//htmlBytes := markdown.ToHTML([]byte(markdownText), nil, renderer)
	html := markdown.ToHTML(mddata, parser, renderer)
	return html
}
