package entity

type Recipe struct {
	ID    string `json:"id,omitempty" bson:"_id,omitempty"`
	Title string `json:"title" bson:"title"`
	URL   string `json:"url" bson:"url"`
}

type Recipes struct {
	Data  []Recipe `json:"data"`
	Total int      `json:"total"`
}
