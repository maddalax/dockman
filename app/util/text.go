package util

import (
	"github.com/microcosm-cc/bluemonday"
	"sync"
)

var once sync.Once
var p *bluemonday.Policy

func GetPolicy() *bluemonday.Policy {
	once.Do(func() {
		p = bluemonday.UGCPolicy()
		p.AllowStyling()
	})
	return p
}

func Sanitize(text string) string {
	return GetPolicy().Sanitize(text)
}
