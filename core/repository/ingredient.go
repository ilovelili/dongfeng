package repository

import (
	"fmt"

	"github.com/ilovelili/dongfeng/core/model"
)

// Ingredient ingredient repository
type Ingredient struct{}

// NewIngredientRepository new ingredient repository
func NewIngredientRepository() *Ingredient {
	db().AutoMigrate(&model.Ingredient{}, &model.IngredientCategory{})
	return new(Ingredient)
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
