package handler

import (
	"doodemo-cook/lib/response"
	"log"
	"net/http"
)

type Tag struct {
	u tagUsecase
}

func NewTag(u tagUsecase) *Tag {
	return &Tag{u: u}
}

func (h *Tag) Find(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tags, cnt, err := h.u.Find(ctx)
	if err != nil {
		response.FromStatusCode(ctx, w, http.StatusInternalServerError)
		log.Printf("fail 'handler.Tag.Find': %v", err)
		return
	}

	res := resTagsFromTags(tags, cnt)
	response.JSON(ctx, w, res, http.StatusOK)
}
