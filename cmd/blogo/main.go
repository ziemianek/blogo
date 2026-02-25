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

const Title string = "Test strona"

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
	// 1. Tworzymy plik (dodajemy też obsługę błędów, to zawsze dobra praktyka)
	f, err := os.Create("output.json")
	if err != nil {
		return err
	}
	defer f.Close()

	// 2. Tworzymy enkoder i podpinamy go BEZPOŚREDNIO pod nasz plik 'f'
	encoder := json.NewEncoder(f)

	// 3. Ustawiamy opcje enkodera
	encoder.SetEscapeHTML(false) // Brak "krzaczków" zamiast tagów HTML
	encoder.SetIndent("", "\t")  // Ładne wcięcia z użyciem tabulatora

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

	articles, err := GetAllArticles("./web/static/articles")
	SortArticlesByModified(articles)
	check(err)

	file, err := os.Create("output.html")
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
