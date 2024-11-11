package urls

import (
	"fmt"
	"github.com/maddalax/htmgo/framework/h"
)

func WithQs(url string, pairs ...string) string {
	if len(pairs) == 0 {
		return url
	}
	qs := h.NewQs(pairs...).ToString()
	return fmt.Sprintf("%s?%s", url, qs)
}

func ResourceUrl(id string) string {
	return WithQs("/resource", "id", id)
}

func NewResourceUrl() string {
	return WithQs("/resource/new")
}
