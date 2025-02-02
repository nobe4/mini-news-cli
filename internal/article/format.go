package article

import (
	"fmt"
	"time"
)

func fixedString(in string, width int) string {
	if len(in) > width {
		return in[:width-3] + "..."
	}

	return fmt.Sprintf("%*s", -width, in)
}

func relTime(in time.Time) string {
	diff := time.Since(in)

	if diff < 1*time.Hour {
		return "<1h"
	}

	hours := int(diff.Hours())
	if diff < 24*time.Hour {
		return fmt.Sprintf("%2dh", hours)
	}

	if diff < 100*time.Hour*24 {
		return fmt.Sprintf("%2dd", hours/24)
	}

	return "+++"
}

func (a Article) Render(width int) string {
	if width == a.renderedWidth {
		return a.rendered
	}
	a.renderedWidth = width

	titleWidth := a.renderedWidth - timeWidth - significanceWidth - separationWidth

	a.rendered = fmt.Sprintf(
		"%s %s %s",
		relTime(a.Publication),
		a.Significance,
		fixedString(a.Title, titleWidth),
	)

	return a.rendered
}
