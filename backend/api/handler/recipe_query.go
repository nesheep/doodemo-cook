package handler

import (
	"net/url"
	"strconv"
)

type recipeQuery struct {
	q     string
	limit int
	skip  int
}

func (h *Recipe) parseQeury(queryRqw url.Values) (recipeQuery, error) {
	q := queryRqw.Get("q")

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

	skip := 0
	qSkip := queryRqw.Get("skip")
	if qSkip != "" {
		s, err := strconv.Atoi(qSkip)
		if err != nil {
			return recipeQuery{}, err
		}
		if s > 0 {
			skip = s
		}
	}

	return recipeQuery{q: q, limit: limit, skip: skip}, nil
}