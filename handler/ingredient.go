package handler

import (
	"net/http"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/ilovelili/dongfeng/core/model"
	"github.com/ilovelili/dongfeng/util"
	"github.com/labstack/echo"
)

// GetIngredients GET /ingredients
func GetIngredients(c echo.Context) error {
	var (
		ingredientsParam = c.QueryParam("ingredients")
		aliases          = strings.Split(ingredientsParam, ",")
		ingredients      []*model.Ingredient
		err              error
	)

	if ingredientsParam == "" {
		ingredients, err = ingredientRepo.FindAll()
	} else {
		ingredients, err = ingredientRepo.FindByAlias(aliases)
	}

	if err != nil {
		return util.ResponseError(c, "400-103", "failed to get ingredients", err)
	}
	return c.JSON(http.StatusOK, ingredients)
}

// SaveIngredients POST /ingredients
func SaveIngredients(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	file, _, err := c.Request().FormFile("file")
	if err != nil {
		// "400-114": "failed to parse ingredients",
		return util.ResponseError(c, "400-114", "failed to parse ingredients", err)
	}
	defer file.Close()

	ingredients := []*model.Ingredient{}
	if err := gocsv.Unmarshal(file, &ingredients); err != nil {
		return util.ResponseError(c, "400-114", "failed to parse ingredients", err)
	}

	categories, err := ingredientRepo.FindAllCategories()
	if err != nil {
		return util.ResponseError(c, "400-120", "failed to get ingredient categories", err)
	}

	categoryMap := make(map[string]uint)
	for _, category := range categories {
		categoryMap[category.Category] = category.ID
	}

	for _, ingredient := range ingredients {
		id, ok := categoryMap[ingredient.CSVIngredientCategory]
		if !ok {
			return util.ResponseError(c, "400-115", "invalid ingredient category", err)
		}
		ingredient.IngredientCategoryID = id
		ingredient.Alias = ingredient.Ingredient
		ingredient.CreatedBy = userInfo.Email
	}

	if err := ingredientRepo.SaveAll(ingredients); err != nil {
		return util.ResponseError(c, "500-121", "failed to save ingredients", err)
	}

	notify(model.IngredientUpdated(userInfo.Email))
	return c.NoContent(http.StatusOK)
}
