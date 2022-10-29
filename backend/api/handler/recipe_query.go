package handler

import (
	"fmt"
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

func (h *Recipe) parseQuery(queryRqw url.Values) (recipeQuery, error) {
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
			return recipeQuery{}, fmt.Errorf("fail 'handler.Recipe.parseQuery': %w", err)
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
			return recipeQuery{}, fmt.Errorf("fail 'handler.Recipe.parseQuery': %w", err)
		}
		if p > 1 {
			page = p
		}
	}

	return recipeQuery{q: q, tags: tags, limit: limit, page: page}, nil
}
