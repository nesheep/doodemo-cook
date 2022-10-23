package entity

type Recipe struct {
	ID    string
	Title string
	URL   string
	Tags  Tags
}

type Recipes []Recipe
