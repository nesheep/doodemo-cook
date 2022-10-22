package handler

import "doodemo-cook/api/entity"

type resTag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func resTagFromTag(tag entity.Tag) resTag {
	return resTag{ID: tag.ID, Name: tag.Name}
}

type resTags struct {
	Data  []resTag `json:"data"`
	Total int      `json:"total"`
}

func resTagsFromTags(tags entity.Tags, total int) resTags {
	data := make([]resTag, 0, len(tags))
	for _, v := range tags {
		data = append(data, resTagFromTag(v))
	}

	return resTags{Data: data, Total: total}
}
