package must

import (
	"net/url"
)

func Url(r string) *url.URL {
	u, err := url.Parse(r)
	if err != nil {
		panic(err)
	}
	return u
}
