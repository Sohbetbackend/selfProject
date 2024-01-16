package api

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/Sohbetbackend/selfProject/internal/app"
	"github.com/Sohbetbackend/selfProject/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Success(c *gin.Context, data gin.H) {
	c.JSON(http.StatusOK, SuccessResponseObject(data))
}

func SuccessResponseObject(data gin.H) gin.H {
	// data["success"] = true
	return data
}

func BindAndValidate(c *gin.Context, r interface{}) (errMessage string, errKey string) {
	if err := c.Bind(r); err != nil {
		errMessage = err.Error()
		return
	}

	v := validator.New()
	if err := v.Struct(r); err != nil {
		err := err.(validator.ValidationErrors)[0]
		errMessage = err.Tag()
		errKey = (err.Field())
		return
	}
	return
}

func handleError(c *gin.Context, err error) {
	if errA, ok := err.(*app.AppError); ok {
		if errA == app.ErrNotFound {
			if utils.GetLoggerDesc() == "" {
				utils.Logger.Error(err)
			}
			c.JSON(http.StatusNotFound, ErrorResponseObject(errA))
		}
	}
}

func ErrorResponseObject(err *app.AppError) gin.H {
	return gin.H{
		"success": false,
		"data":    nil,
		"error": gin.H{
			"code":    err.Code(),
			"key":     err.Key(),
			"comment": err.Comment(),
		},
	}
}

func handleFile(c *gin.Context, handler *multipart.FileHeader, folder string) (string, error) {
	var err error
	path := "web/uploads/" + folder + "/"
	fParts := strings.Split(handler.Filename, ".")
	ext := "." + strings.ToLower(fParts[len(fParts)-1])
	if ext != ".pdf" && ext != ".txt" && ext != ".jpg" && ext != ".png" {
		return "", errors.New("invalid extension: " + ext)
	}
	name := handler.Filename + ext
	err = c.SaveUploadedFile(handler, path+name)
	if err != nil {
		utils.LoggerDesc("HTTP error").Error(err)
		return "", err
	}
	return folder + "/" + name, nil
}

func handleFilesUpload(c *gin.Context, key string, folder string) ([]string, error) {
	c.Request.ParseMultipartForm(100)
	form, err := c.MultipartForm()
	if err != nil {
		return nil, nil
	}
	files := form.File[key]

	paths := []string{}
	for _, f := range files {
		path, err := handleFile(c, f, folder)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}
