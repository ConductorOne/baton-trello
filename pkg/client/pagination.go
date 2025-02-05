package client

import (
	"net/url"
	"strconv"
)

const ItemsPerPage = 50

type PageOptions struct {
	PerPage int `url:"limit,omitempty"`
	Page    int `url:"page,omitempty"`
}

type ReqOpt func(reqURL *url.URL)

// Number of items to return.
func WithPageLimit(pageLimit int) ReqOpt {
	if pageLimit <= 0 || pageLimit > ItemsPerPage {
		pageLimit = ItemsPerPage
	}
	return WithQueryParam("limit", strconv.Itoa(pageLimit))
}

// Number for the page (inclusive). The page number starts with 1.
// If page is 0, first page is assumed.
func WithPage(page int) ReqOpt {
	if page == 0 {
		page = 1
	}
	return WithQueryParam("page", strconv.Itoa(page))
}

func WithQueryParam(key string, value string) ReqOpt {
	return func(reqURL *url.URL) {
		q := reqURL.Query()
		q.Set(key, value)
		reqURL.RawQuery = q.Encode()
	}
}
