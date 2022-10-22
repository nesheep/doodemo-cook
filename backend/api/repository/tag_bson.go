package repository

import "doodemo-cook/api/entity"

type bTag struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
}

func (t *bTag) toTag() entity.Tag {
	return entity.Tag{ID: t.ID, Name: t.Name}
}
