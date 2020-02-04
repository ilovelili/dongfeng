package repository

import (
	"github.com/ilovelili/dongfeng/core/model"
)

// Recipe recipe repository
type Recipe struct{}

// NewRecipeRepository new recipe repository
func NewRecipeRepository() *Recipe {
	db().AutoMigrate(&model.Menu{}, &model.Recipe{}, &model.RecipeNutrition{}, &model.Ingredient{})
	return new(Recipe)
}

// Find find recipes by names
func (r *Recipe) Find(names []string) ([]*model.Recipe, error) {
	recipes := []*model.Recipe{}
	err := db().Where("recipes.name IN (?)", names).Preload("Ingredient").Preload("IngredientNutrition").Find(&recipes).Error
	return recipes, err
}

// SaveAll save all recipes
func (r *Recipe) SaveAll(recipes []*model.Recipe) error {
	tx := db().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	for _, recipe := range recipes {
		if err := tx.Create(recipe).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
