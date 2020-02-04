package handler

import (
	"net/http"

	"github.com/ilovelili/dongfeng/util"
	"github.com/labstack/echo"
)

// GetMenus GET /menus
func GetMenus(c echo.Context) error {
	juniorOrSenior, breakfastOrLunch, from, to := c.QueryParam("junior_or_senior"), c.QueryParam("breakfast_or_lunch"), c.QueryParam("from"), c.QueryParam("to")
	menus, err := menuRepo.Find(from, to, juniorOrSenior, breakfastOrLunch)
	if err != nil {
		return util.ResponseError(c, "500-122", "failed to get menus", err)
	}
	return c.JSON(http.StatusOK, menus)
}
