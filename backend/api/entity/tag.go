package entity

type Tag struct {
	ID   string `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
}

type Tags struct {
	Data  []Tag `json:"data" bson:"data"`
	Total int   `json:"total"`
}
