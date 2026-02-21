package main

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const Title string = "Test strona"

type Article struct {
	Title   string
	Content []byte
}

type config struct {
	Title string
	// Article     string
	ArticleList []Article
}

func GetArticleName(filepath string) string {
	articleSplitted := strings.Split(filepath, "/")
	return articleSplitted[len(articleSplitted)-1]
}

// todo: Make it use pointers
func ParseArticleName(articleName string, caser cases.Caser) string {
	articleNameNoFileExt := strings.Join(strings.Split(articleName, ".md"), " ")
	articleNameParsed := caser.String(strings.Join(strings.Split(articleNameNoFileExt, "-"), " "))
	return articleNameParsed
}

// todo: refactor this shit
func GetAllArticles(articleDir string) ([]Article, error) {
	articles := []Article{}
	caser := cases.Title(language.English)
	filenames, err := filepath.Glob(fmt.Sprintf("%v/*.md", articleDir))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for _, a := range filenames {
		articleContent, err := os.ReadFile(a)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		articleContent = mdToHTML(articleContent)
		articleName := ParseArticleName(GetArticleName(a), caser)
		articles = append(articles, Article{
			Title:   articleName,
			Content: articleContent,
		})
	}
	return articles, nil
}

// todo: make it use a poiunter to article content
func mdToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	return markdown.Render(doc, renderer)
}

func main() {
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	articles, err := GetAllArticles("./web/static/articles")
	check(err)

	file, err := os.Create("output.html")
	check(err)
	defer file.Close()

	tpl, err := template.ParseFiles("./web/static/html/index.tmpl")
	check(err)

	tpl.Execute(file, &config{
		Title: Title,
		// Article:     string(html),
		ArticleList: articles,
	})
	check(err)
}
