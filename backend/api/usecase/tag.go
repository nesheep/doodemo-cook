package usecase

import (
	"context"
	"doodemo-cook/api/entity"
	"fmt"
)

type Tag struct {
	r tagRepository
}

func NewTag(r tagRepository) *Tag {
	return &Tag{r: r}
}

func (u *Tag) Find(ctx context.Context) (entity.Tags, int, error) {
	tags, err := u.r.Find(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("fail 'usecase.Tag.Find': %w", err)
	}

	cnt, err := u.r.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("fail 'usecase.Tag.Find': %w", err)
	}

	return tags, cnt, nil
}
