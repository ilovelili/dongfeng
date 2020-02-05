package handler

import (
	"net/http"
	"strings"

	"github.com/ilovelili/dongfeng/util"
	"github.com/labstack/echo"
)

// GetRecipes GET /recipes
func GetRecipes(c echo.Context) error {
	recipesParam := c.QueryParam("recipes")
	recipeIDs := []string{}
	if recipesParam != "" {
		recipeIDs = strings.Split(recipesParam, ",")
	}

	recipes, err := recipeRepo.Find(recipeIDs)
	if err != nil {
		return util.ResponseError(c, "500-123", "failed to get recipes", err)
	}

	return c.JSON(http.StatusOK, recipes)
}
