package renders
import (
    "io/ioutil"
    "github.com/gomarkdown/markdown"
)
// render markdown file to html file

func RenderMd(file_name string) ([]byte, error) {
    // read file
    file, err := ioutil.ReadFile(file_name)

    // render use gomarkdown
    html := markdown.ToHTML(file, nil, nil)
    // append mathjax script to html
    // html = append(html, []byte(`<script type="text/javascript" async src="https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.5/MathJax.js?config=TeX-MML-AM_CHTML"></script>`)...)
    // write file
    // ioutil.WriteFile("test.html", html, 0644)
    return html, err
}
