package repository

import (
	"fmt"

	"github.com/ilovelili/dongfeng/core/model"
)

// Ingredient ingredient repository
type Ingredient struct{}

// NewIngredientRepository new ingredient repository
func NewIngredientRepository() *Ingredient {
	db().AutoMigrate(&model.Menu{}, &model.Recipe{}, &model.RecipeNutrition{}, &model.Ingredient{}, &model.IngredientCategory{})
	return new(Ingredient)
}

// FindAll find all ingredients
func (r *Ingredient) FindAll() ([]*model.Ingredient, error) {
	ingredients := []*model.Ingredient{}
	err := db().Preload("IngredientCategory").Find(&ingredients).Error
	return ingredients, err
}

// FindAllCategories find all ingredient categories
func (r *Ingredient) FindAllCategories() ([]*model.IngredientCategory, error) {
	ingredientCategories := []*model.IngredientCategory{}
	err := db().Find(&ingredientCategories).Error
	return ingredientCategories, err
}

// FindByAlias find ingredients by alias
func (r *Ingredient) FindByAlias(aliases []string) ([]*model.Ingredient, error) {
	ingredients := []*model.Ingredient{}
	err := db().Where("ingredients.alias IN (?)", aliases).Preload("IngredientCategory").Find(&ingredients).Error
	return ingredients, err
}

// FindByPattern find ingredients by pattern
func (r *Ingredient) FindByPattern(pattern string) ([]*model.Ingredient, error) {
	ingredients := []*model.Ingredient{}
	err := db().Where("ingredients.alias LIKE ?", fmt.Sprintf("%%%s%%", pattern)).Preload("IngredientCategory").Find(&ingredients).Error
	return ingredients, err
}

// SaveAll save all ingredients
func (r *Ingredient) SaveAll(ingredients []*model.Ingredient) error {
	tx := db().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	for _, ingredient := range ingredients {
		_ingredient := new(model.Ingredient)
		if tx.Where("ingredients.ingredient = ?", ingredient.Ingredient).Find(&_ingredient).RecordNotFound() {
			// if not found, insert. else update
			ingredient.ID = 0
		} else {
			ingredient.ID = _ingredient.ID
		}

		if err := tx.Save(ingredient).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
