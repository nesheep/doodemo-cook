package handler

import (
	"net/url"
	"strconv"
	"strings"
)

type recipeQuery struct {
	q     string
	tags  []string
	limit int
	page  int
}

func (q recipeQuery) skip() int {
	return q.limit * (q.page - 1)
}

func (h *Recipe) parseQeury(queryRqw url.Values) (recipeQuery, error) {
	q := queryRqw.Get("q")

	tags := []string{}
	qTags := queryRqw.Get("tags")
	if qTags != "" {
		tags = strings.Split(qTags, ",")
	}

	limit := 10
	qLimit := queryRqw.Get("limit")
	if qLimit != "" {
		l, err := strconv.Atoi(qLimit)
		if err != nil {
			return recipeQuery{}, err
		}
		if l > 0 {
			limit = l
		}
	}

	page := 1
	qPage := queryRqw.Get("page")
	if qPage != "" {
		p, err := strconv.Atoi(qPage)
		if err != nil {
			return recipeQuery{}, err
		}
		if p > 1 {
			page = p
		}
	}

	return recipeQuery{q: q, tags: tags, limit: limit, page: page}, nil
}
