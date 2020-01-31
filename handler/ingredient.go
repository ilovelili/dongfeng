package handler

import (
	"net/http"
	"strings"

	"github.com/ilovelili/dongfeng/util"
	"github.com/labstack/echo"
)

// GetIngredients GET /ingredients
func GetIngredients(c echo.Context) error {
	var aliases = strings.Split(c.QueryParam("ingredients"), ",")
	ingredients, err := ingredientRepo.FindByAlias(aliases)
	if err != nil {
		return util.ResponseError(c, "400-103", "failed to get ingredients", err)
	}
	return c.JSON(http.StatusOK, ingredients)
}
