package entity

type Recipe struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

type Recipes struct {
	Data  []Recipe `json:"data"`
	Total int      `json:"total"`
}
