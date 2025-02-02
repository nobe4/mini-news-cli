package article

import (
	"log/slog"
	"net/http"
	"os/exec"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	URL = "https://www.newsminimalist.com/?sort=significance"

	separationWidth   = 2
	timeWidth         = 3
	significanceWidth = 5
)

type Article struct {
	Significance string
	Title        string
	Publication  time.Time
	URL          string

	renderedWidth int
	rendered      string
}

func (a Article) Open() error {
	return exec.Command("open", a.URL).Run()
}

type Articles []Article

func Get() (Articles, error) {
	articles := []Article{}

	res, err := http.Get(URL)
	if err != nil {
		return articles, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return articles, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return articles, err
	}

	doc.Find("details").Each(func(_ int, s *goquery.Selection) {
		article, err := parse(s)
		if err != nil {
			slog.Error("failed to parse article", "article", article, "error", err)
			return
		}

		articles = append(articles, article)
	})

	return articles, nil
}

func parse(s *goquery.Selection) (Article, error) {
	a := Article{
		Significance: s.Find("summary > :nth-child(1)").Text(),
		Title:        s.Find("summary > :nth-child(2) > :nth-child(1)").Text(),
	}

	if rawDate, ok := s.Find("summary  > :nth-child(3)").Attr("title"); ok && rawDate != "" {
		if date, err := time.Parse("Mon, Jan 2, 2006, 3 PM", rawDate); err == nil {
			a.Publication = date
		} else {
			return a, err
		}
	}

	if link, ok := s.Find("a").Attr("href"); ok {
		a.URL = link
	}

	return a, nil
}
