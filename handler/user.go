package handler

import (
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/ilovelili/aliyun-client/oss"
	"github.com/ilovelili/dongfeng/core/model"
	"github.com/ilovelili/dongfeng/core/repository"
	"github.com/ilovelili/dongfeng/util"
	"github.com/labstack/echo"
)

// UploadAvatar POST /user/upload
func UploadAvatar(c echo.Context) error {
	aliyunsvc := oss.NewService(config.OSS.APIKey, config.OSS.APISecret)
	aliyunsvc.SetEndPoint(config.OSS.Endpoint)
	aliyunsvc.SetBucket(config.OSS.BucketName)

	form, err := c.MultipartForm()
	if err != nil {
		return util.ResponseError(c, http.StatusBadRequest, "400-100", "failed to parse multipart form", err)
	}

	files := form.File["image"]
	locations := []string{}
	for _, file := range files {
		if !supportedImageMimeType(file.Header["Content-Type"]) {
			return util.ResponseError(c, http.StatusBadRequest, "400-101", "unsupported mimetype", err)
		}

		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.Create(file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		io.Copy(dst, src)

		opts := &oss.UploadOptions{
			Public:     true,
			ObjectName: dst.Name(),
		}

		resp := aliyunsvc.Upload(opts)
		if resp.Error != nil {
			return util.ResponseError(c, http.StatusInternalServerError, "500-102", "failed to upload avatar", err)
		}

		locations = append(locations, resp.Location)
	}

	return c.JSON(http.StatusOK, locations)
}

// UpdateUser POST /user/update
func UpdateUser(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return util.ResponseError(c, http.StatusBadRequest, "400-102", "failed to bind user", err)
	}

	userRepo := repository.NewUserRepository()
	existingUser, err := userRepo.FindByEmail(user.Email)
	if err != nil {
		return util.ResponseError(c, http.StatusInternalServerError, "500-103", "no user", err)
	}

	user.ID = existingUser.ID
	user.Role = existingUser.Role

	if user.Name == "" {
		user.Name = existingUser.Name
	}

	if user.Photo == "" {
		user.Photo = existingUser.Photo
	}

	err = userRepo.Save(user)
	if err != nil {
		return util.ResponseError(c, http.StatusInternalServerError, "500-100", "failed to save user", err)
	}

	return c.NoContent(http.StatusOK)
}

// supportedImageMimeType check if uploaded file is image
func supportedImageMimeType(contenttype []string) bool {
	r := regexp.MustCompile("image/(png|jpeg|gif)")
	for _, ct := range contenttype {
		if r.MatchString(ct) {
			return true
		}
	}

	return false
}
