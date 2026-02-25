package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const Title string = "Ziemski Blog"

type Metadata struct {
	LastModified    time.Time
	LastModifiedStr string
}

type Article struct {
	Title    string        `json:"title"`
	Content  template.HTML `json:"content"`
	Metadata Metadata      `json:"metadata"`
}

type config struct {
	Title       string
	Article     Article
	ArticleList []Article
}

func GetArticleName(filepath string) string {
	articleSplitted := strings.Split(filepath, "/")
	return articleSplitted[len(articleSplitted)-1]
}

func SaveArticleToJson(articles []Article) error {
	f, err := os.Create("./dist/data.js")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString("const myArticles = ")
	if err != nil {
		fmt.Println("Error writing prefix:", err)
		return err
	}
	encoder := json.NewEncoder(f)

	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "\t")

	// 4. Kodujemy tablicę articles bezpośrednio do pliku
	err = encoder.Encode(articles)
	if err != nil {
		return err
	}

	return nil
}

func SortArticlesByModified(articles []Article) {
	slices.SortFunc(articles, func(a, b Article) int {
		return b.Metadata.LastModified.Compare(a.Metadata.LastModified)
	})
}

func FormatDate(date time.Time) string {
	return date.Format("2006-01-02")
}

// todo: Make it use pointers
func ParseArticleName(articleName string, caser cases.Caser) string {
	articleNameNoFileExt := strings.Join(strings.Split(articleName, ".md"), " ")
	articleNameParsed := caser.String(strings.Join(strings.Split(articleNameNoFileExt, "-"), " "))
	return articleNameParsed
}

func GetFileLastUpdated(filepath string) (time.Time, error) {
	metadata, err := os.Stat(filepath)
	if err != nil {
		log.Fatal(err)
		return time.Time{}, err
	}
	return metadata.ModTime(), nil
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
	fmt.Println(articleDir)
	fmt.Println(filenames)
	for _, a := range filenames {
		articleContent, err := os.ReadFile(a)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		articleContent = mdToHTML(articleContent)
		articleName := ParseArticleName(GetArticleName(a), caser)
		articleLastUpdated, err := GetFileLastUpdated(a)
		fmt.Println(articleName, articleLastUpdated)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		articles = append(articles, Article{
			Title:   articleName,
			Content: template.HTML(articleContent),
			Metadata: Metadata{
				LastModified:    articleLastUpdated,
				LastModifiedStr: FormatDate(articleLastUpdated),
			},
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

	articles, err := GetAllArticles("./articles")
	SortArticlesByModified(articles)
	check(err)

	file, err := os.Create("./dist/index.html")
	check(err)
	defer file.Close()

	SaveArticleToJson(articles)

	tpl := template.Must(template.ParseFiles("./web/static/html/index.tmpl"))

	tpl.Execute(file, &config{
		Title:       Title,
		Article:     articles[0],
		ArticleList: articles,
	})
	check(err)
}
