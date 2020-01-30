package handler

import (
	"net/http"
	"strconv"

	"github.com/ilovelili/dongfeng/core/model"
	"github.com/ilovelili/dongfeng/util"
	"github.com/labstack/echo"
)

// GetMasters GET /masters/:id
func GetMasters(c echo.Context) error {
	var (
		id     = c.QueryParam("id")
		master interface{}
		err    error
	)

	switch id {
	case strconv.Itoa(int(model.AgeHeightWeightP)):
		master, err = physiqueRepo.SelectAgeHeightWeightPMasters()
		break
	case strconv.Itoa(int(model.AgeHeightWeightSD)):
		master, err = physiqueRepo.SelectAgeHeightWeightSDMasters()
		break
	case strconv.Itoa(int(model.BMI)):
		master, err = physiqueRepo.SelectBMIMasters()
		break
	case strconv.Itoa(int(model.HeightToWeightP)):
		master, err = physiqueRepo.SelectHeightToWeightPMasters()
		break
	case strconv.Itoa(int(model.HeightToWeightSD)):
		master, err = physiqueRepo.SelectHeightToWeightSDMasters()
		break
	default:
		return util.ResponseError(c, "400-112", "invalid master id", err)
	}

	if err != nil {
		return util.ResponseError(c, "500-115", "failed to get physique masters", err)
	}

	return c.JSON(http.StatusOK, master)
}
