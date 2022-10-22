package handler

import "doodemo-cook/api/entity"

type reqRecipe struct {
	Title string   `json:"title"`
	URL   string   `json:"url"`
	Tags  []string `json:"tags"`
}

func (r reqRecipe) toRecipe() entity.Recipe {
	tags := make(entity.Tags, 0, len(r.Tags))
	for _, v := range r.Tags {
		tags = append(tags, entity.Tag{Name: v})
	}

	return entity.Recipe{
		Title: r.Title,
		URL:   r.URL,
		Tags:  tags,
	}
}

type resRecipe struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	URL   string   `json:"url"`
	Tags  []resTag `json:"tags"`
}

func resRecipeFromRecipe(recipe entity.Recipe) resRecipe {
	tags := make([]resTag, 0, len(recipe.Tags))
	for _, v := range recipe.Tags {
		tags = append(tags, resTagFromTag(v))
	}

	return resRecipe{
		ID:    recipe.ID,
		Title: recipe.Title,
		URL:   recipe.URL,
		Tags:  tags,
	}
}

type resRecipes struct {
	Data  []resRecipe `json:"data"`
	Total int         `json:"total"`
}

func resRecipesFromRecipes(recipes entity.Recipes, total int) resRecipes {
	data := make([]resRecipe, 0, len(recipes))
	for _, v := range recipes {
		data = append(data, resRecipeFromRecipe(v))
	}

	return resRecipes{Data: data, Total: total}
}
