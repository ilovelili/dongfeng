package handler

import (
	"net/http"
	"strconv"

	"github.com/gocarina/gocsv"
	"github.com/ilovelili/dongfeng/core/controller"
	"github.com/ilovelili/dongfeng/core/model"
	"github.com/ilovelili/dongfeng/util"
	"github.com/labstack/echo"
)

// GetPhysiques GET /physiques
func GetPhysiques(c echo.Context) error {
	var (
		err     error
		pupilID uint = 0
		classID uint = 0
		pupils       = []*model.Pupil{}
		year         = c.QueryParam("year")
		class        = c.QueryParam("class")
		pupil        = c.QueryParam("name")
	)

	if class != "" {
		_classID, err := strconv.ParseUint(class, 10, 64)
		if err != nil {
			return util.ResponseError(c, "400-108", "invalid class id", err)
		}
		classID = uint(_classID)
	}
	if pupil != "" {
		_pupilID, err := strconv.ParseUint(pupil, 10, 64)
		if err != nil {
			return util.ResponseError(c, "400-109", "invalid pupil id", err)
		}
		pupilID = uint(_pupilID)
	}

	// retrieve pupils
	if pupilID != 0 {
		pupils, err = pupilRepo.FindByPupilID(pupilID)
	} else if classID != 0 {
		pupils, err = pupilRepo.FindByClassID(classID)
	} else {
		pupils, err = pupilRepo.Find("", year)
	}
	if err != nil {
		return util.ResponseError(c, "500-107", "failed to get pupils", err)
	}

	pupilIDs := []uint{}
	for _, pupil := range pupils {
		pupilIDs = append(pupilIDs, pupil.ID)
	}

	physiques, err := physiqueRepo.Find(pupilIDs)
	if err != nil {
		return util.ResponseError(c, "500-116", "failed to get physiques", err)
	}

	return c.JSON(http.StatusOK, physiques)
}

// UpdatePhysique PUT /physique
func UpdatePhysique(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	physique := new(model.Physique)
	if err := c.Bind(physique); err != nil {
		return util.ResponseError(c, "400-113", "failed to parse physiques", err)
	}

	physiqueCtrl := controller.NewPhysiqueController()
	physique.CreatedBy = userInfo.Email
	// resolve physique
	if err := physiqueCtrl.ResolvePhysique(physique); err != nil {
		return util.ResponseError(c, "500-118", "failed to resolve physiques", err)
	}

	if err := physiqueRepo.Save(physique); err != nil {
		return util.ResponseError(c, "500-117", "failed to save physiques", err)
	}

	notify(model.PhysiqueUpdated(userInfo.Email))
	return c.NoContent(http.StatusOK)
}

// SavePhysiques POST /physiques
func SavePhysiques(c echo.Context) error {
	userInfo, _ := c.Get("userInfo").(model.User)
	file, _, err := c.Request().FormFile("file")
	if err != nil {
		return util.ResponseError(c, "400-113", "failed to parse physiques", err)
	}
	defer file.Close()

	physiques := []*model.Physique{}
	if err := gocsv.Unmarshal(file, &physiques); err != nil {
		return util.ResponseError(c, "400-113", "failed to parse physiques", err)
	}

	physiqueCtrl := controller.NewPhysiqueController()
	for _, physique := range physiques {
		physique.CreatedBy = userInfo.Email
		// resolve physique
		if err := physiqueCtrl.ResolvePhysique(physique); err != nil {
			return util.ResponseError(c, "500-118", "failed to resolve physiques", err)
		}
	}

	if err := physiqueRepo.SaveAll(physiques); err != nil {
		return util.ResponseError(c, "500-117", "failed to save physiques", err)
	}

	notify(model.PhysiqueUpdated(userInfo.Email))
	return c.NoContent(http.StatusOK)
}
