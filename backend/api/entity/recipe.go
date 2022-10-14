package entity

type Recipe struct {
	ID    string   `json:"id,omitempty" bson:"_id,omitempty"`
	Title string   `json:"title" bson:"title"`
	URL   string   `json:"url" bson:"url"`
	Tags  []string `json:"tags" bson:"tags"`
}

func NewRecipe() Recipe {
	return Recipe{Tags: []string{}}
}

type Recipes struct {
	Data  []Recipe `json:"data"`
	Total int      `json:"total"`
}

type RecipeWithTags struct {
	ID    string `json:"id,omitempty" bson:"_id,omitempty"`
	Title string `json:"title" bson:"title"`
	URL   string `json:"url" bson:"url"`
	Tags  []Tag  `json:"tags" bson:"tags"`
}

func NewRecipeWithTags() RecipeWithTags {
	return RecipeWithTags{Tags: []Tag{}}
}

type RecipesWithTags struct {
	Data  []RecipeWithTags `json:"data"`
	Total int              `json:"total"`
}
