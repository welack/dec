package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Site struct {
	Title       string
	Subtitle    string
	Description string
	URL         string
}

type Template struct {
	Title       string
	Date        string
	Body        string
	Description string
	Canonical   string
	OGType      string
	Site        Site
	Submenu     []SubmenuConfig
	Theme       template.HTML
}

type IndexData struct {
	Site        Site
	Posts       []Post
	Pagination  Pagination
	Submenu     []SubmenuConfig
	Theme       template.HTML
	Description string
	Canonical   string
}

type Post struct {
	Title string
	Date  string
	URL   string
}

type Pagination struct {
	CurrentPage int
	TotalPages  int
	HasPrev     bool
	HasNext     bool
	PrevURL     string
	NextURL     string
}

var htmlTagRE = regexp.MustCompile(`<[^>]*>`)

func stripHTML(s string) string {
	text := htmlTagRE.ReplaceAllString(s, " ")
	return strings.Join(strings.Fields(text), " ")
}

func truncateText(s string, max int) string {
	runes := []rune(strings.TrimSpace(s))
	if len(runes) <= max {
		return string(runes)
	}
	return strings.TrimSpace(string(runes[:max])) + "..."
}

func metaDescription(page *Page, body string) string {
	if page != nil && page.Meta != nil {
		if desc, ok := page.Meta["description"].(string); ok && strings.TrimSpace(desc) != "" {
			return strings.TrimSpace(desc)
		}
	}
	if page != nil {
		if plain := stripHTML(page.Content); plain != "" {
			return truncateText(plain, 160)
		}
	}
	if plain := stripHTML(body); plain != "" {
		return truncateText(plain, 160)
	}
	if strings.TrimSpace(conf.Main.Description) != "" {
		return strings.TrimSpace(conf.Main.Description)
	}
	return strings.TrimSpace(conf.Main.Subtitle)
}

func canonicalURL(path string) string {
	base := strings.TrimRight(conf.Main.Url, "/")
	if base == "" {
		return path
	}
	if path == "" {
		return base + "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return base + path
}

func themeStyle() template.HTML {
	c := conf.Colorscheme
	defaults := map[string]string{
		"bg":               "#ffffff",
		"text":             "#000000",
		"primary":          "#507538",
		"primary-hover":    "#3c5a2a",
		"primary-text":     "#ffffff",
		"muted":            "#666666",
		"muted-on-primary": "#d8e2c8",
		"border":           "#dddddd",
	}

	values := map[string]string{
		"bg":               c.BG,
		"text":             c.Text,
		"primary":          c.Primary,
		"primary-hover":    c.PrimaryHover,
		"primary-text":     c.PrimaryText,
		"muted":            c.Muted,
		"muted-on-primary": c.MutedOnPrimary,
		"border":           c.Border,
	}

	var overrides []string
	for k, def := range defaults {
		if v := strings.TrimSpace(values[k]); v != "" && v != def {
			overrides = append(overrides, fmt.Sprintf("--%s: %s;", k, v))
		}
	}
	if len(overrides) == 0 {
		return ""
	}
	sort.Strings(overrides)
	return template.HTML("<style>:root{" + strings.Join(overrides, "") + "}</style>")
}

func (t *Template) render(file, header, footer, outpath, htmlContent string, page *Page) error {

	funcMap := template.FuncMap{
		"dateFormat": func(format string, t time.Time) string { return t.Format(format) },
		"safeHTML":   func(s string) template.HTML { return template.HTML(s) },
		"urlize":     func(s string) string { return strings.ToLower(strings.ReplaceAll(s, " ", "-")) },
		"now":        func() time.Time { return time.Now() },
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFiles(file, header, footer)
	if err != nil {
		return fmt.Errorf("%w\n", err)
	}

	title, ok := page.Meta["title"].(string)
	if !ok {
		return fmt.Errorf("missing title in frontmatter")
	}

	var dateStr string
	switch d := page.Meta["date"].(type) {
	case string:
		dateStr = d
	case time.Time:
		dateStr = d.Format("2006-01-02")
	case nil:
		dateStr = ""
	default:
		return fmt.Errorf("invalid date type in frontmatter")
	}

	pagePath := "/" + strings.TrimSuffix(filepath.Base(page.FilePath), ".md") + "/"
	if strings.TrimSuffix(filepath.Base(page.FilePath), ".md") == "index" {
		pagePath = "/"
	}

	data := Template{
		Title:       title,
		Date:        dateStr,
		Body:        htmlContent,
		Description: metaDescription(page, htmlContent),
		Canonical:   canonicalURL(pagePath),
		OGType:      "article",
		Site: Site{
			Title:       conf.Main.Title,
			Subtitle:    conf.Main.Subtitle,
			Description: conf.Main.Description,
			URL:         strings.TrimRight(conf.Main.Url, "/"),
		},
		Submenu: conf.Submenu,
		Theme:   themeStyle(),
	}

	if err := os.MkdirAll(filepath.Dir(outpath), 0755); err != nil {
		return fmt.Errorf("%w", err)
	}
	outFile, err := os.Create(outpath)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer outFile.Close()

	return tmpl.ExecuteTemplate(outFile, filepath.Base(file), data)
}
