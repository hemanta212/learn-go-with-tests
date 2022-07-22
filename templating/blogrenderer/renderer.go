package blogrenderer

import (
	"embed"
	"html/template"
	"io"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

var (
	//go:embed "templates/*"
	postTemplates embed.FS
)

type Post struct {
	Title, Description, Body string
	Tags                     []string
}

type PostRenderer struct {
	templ    *template.Template
	mdParser *parser.Parser
}

func NewPostRenderer() (*PostRenderer, error) {
	templ, err := template.ParseFS(postTemplates, "templates/*.gohtml")
	if err != nil {
		return nil, err
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)

	return &PostRenderer{templ: templ, mdParser: parser}, nil
}

func (r *PostRenderer) Render(w io.Writer, post Post) error {
	if err := r.templ.ExecuteTemplate(w, "blog.gohtml", newPostVM(post, r)); err != nil {
		return err
	}
	return nil
}

func (r *PostRenderer) RenderIndex(w io.Writer, posts []Post) error {
	if err := r.templ.ExecuteTemplate(w, "index.gohtml", posts); err != nil {
		return err
	}
	return nil
}

func (p Post) SanitizedTitle() string {
	return strings.ToLower(strings.Replace(p.Title, " ", "-", -1))
}

type postViewModel struct {
	Post
	HTMLBody template.HTML
}

func newPostVM(post Post, r *PostRenderer) postViewModel {
	vm := postViewModel{Post: post}
	md := []byte(post.Body)
	vm.HTMLBody = template.HTML(markdown.ToHTML(md, r.mdParser, nil))
	return vm
}
