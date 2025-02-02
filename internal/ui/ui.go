package ui

import (
	"time"

	"github.com/nobe4/mini-news-cli/internal/article"

	"github.com/gdamore/tcell/v2"
)

func Run() error {
	s, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := s.Init(); err != nil {
		return err
	}

	write(s, 0, 0, "fetching articles...", tcell.StyleDefault)
	s.Show()

	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	current := 0
	articles := article.Articles{}
	go func() {
		var err error
		articles, err = article.Get()
		if err := s.PostEvent(GetEvent{init: time.Now(), err: err}); err != nil {
			panic(err)
		}
	}()

	for {
		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			renderArticles(s, articles, current)
		case *tcell.EventKey:
			if ev.Rune() == 'q' || ev.Key() == tcell.KeyCtrlC {
				return nil
			}

			if ev.Rune() == 'o' {
				if err := articles[current].Open(); err != nil {
					write(s, 0, 0, err.Error(), tcell.StyleDefault.Background(tcell.ColorRed))
					s.Show()
				}
			}

			if ev.Key() == tcell.KeyUp {
				current -= 1
				if current < 0 {
					current = len(articles) - 1
				}
				renderArticles(s, articles, current)
			}
			if ev.Key() == tcell.KeyDown {
				current += 1
				if current > len(articles)-1 {
					current = 0
				}
				renderArticles(s, articles, current)
			}

		case GetEvent:
			if err := ev.Error(); err == nil {
				renderArticles(s, articles, current)
			} else {
				write(s, 0, 0, err.Error(), tcell.StyleDefault.Background(tcell.ColorRed))
				s.Show()
			}
		}
	}
}

func write(s tcell.Screen, x, y int, str string, style tcell.Style) {
	for i, ch := range str {
		s.SetContent(x+i, y, ch, nil, style)
	}
}

func renderArticles(s tcell.Screen, a article.Articles, current int) {
	if len(a) == 0 {
		return
	}

	w, _ := s.Size()

	for i, a := range a {
		style := tcell.StyleDefault
		if i == current {
			style = tcell.StyleDefault.Reverse(true)
		}

		write(s, 0, i, a.Render(w), style)
	}

	s.Show()
}

type GetEvent struct {
	init time.Time
	err  error
}

func (e GetEvent) When() time.Time {
	return e.init
}

func (e GetEvent) Error() error {
	return e.err
}
