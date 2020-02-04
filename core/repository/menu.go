package repository

import "github.com/ilovelili/dongfeng/core/model"

// Menu menu repository
type Menu struct{}

// NewMenuRepository new menu repository
func NewMenuRepository() *Menu {
	db().AutoMigrate(&model.Menu{}, &model.Recipe{}, &model.RecipeNutrition{}, &model.Ingredient{})
	return new(Menu)
}

// Find find menus
func (r *Menu) Find(from, to, juniorOrSenior, breakfastOrLunch string) ([]*model.Menu, error) {
	menus := []*model.Menu{}
	query := db().Joins("JOIN recipes ON menus.recipe_id = recipes.id")

	if breakfastOrLunch != model.AllMenuType.String() && juniorOrSenior != model.AllClass.String() {
		query = query.Where("menus.date BETWEEN ? AND ? AND recipes.name != '未排菜' AND menus.breakfast_or_lunch = ? AND menus.junior_or_senior = ?", from, to, breakfastOrLunch, juniorOrSenior)
	} else if breakfastOrLunch == model.AllMenuType.String() && juniorOrSenior != model.AllClass.String() {
		query = query.Where("menus.date BETWEEN ? AND ? AND recipes.name != '未排菜' AND menus.junior_or_senior = ?", from, to, juniorOrSenior)
	} else if breakfastOrLunch != model.AllMenuType.String() && juniorOrSenior == model.AllClass.String() {
		query = query.Where("menus.date BETWEEN ? AND ? AND recipes.name != '未排菜' AND menus.breakfast_or_lunch = ?", from, to, breakfastOrLunch)
	} else {
		query = query.Where("menus.date BETWEEN ? AND ? AND recipes.name != '未排菜'", from, to)
	}

	err := query.Preload("Recipe").Find(&menus).Error
	return menus, err
}
