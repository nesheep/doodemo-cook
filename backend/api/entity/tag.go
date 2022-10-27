package entity

type Tag struct {
	ID        string
	Name      string
	RecipeNum int
}

type Tags []Tag

func (t Tags) Unique() Tags {
	tags := make(Tags, 0, len(t))
	names := make(map[string]bool, len(t))
	for _, v := range t {
		if !names[v.Name] {
			names[v.Name] = true
			tags = append(tags, v)
		}
	}
	return tags
}
