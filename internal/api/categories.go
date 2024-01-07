package api

import (
	"strconv"

	"github.com/Sohbetbackend/selfProject/internal/app"
	"github.com/Sohbetbackend/selfProject/internal/models"
	"github.com/gin-gonic/gin"
)

func CategoriesRoutes(api *gin.RouterGroup) {
	categoryRoutes := api.Group("/categories")
	{
		categoryRoutes.GET("", CategoryList)
		categoryRoutes.POST("", CategoryCreate)
		categoryRoutes.PUT(":id", CategoryUpdate)
		categoryRoutes.DELETE("", CategoryDelete)
	}
}

func CategoryList(c *gin.Context) {
	r := models.CategoryFilterRequest{}
	if errMsg, errKey := BindAndValidate(c, &r); errMsg != "" || errKey != "" {
		handleError(c, app.NewAppError(errMsg, errKey, ""))
		return
	}
	categories, total, err := app.CategoryList(r)
	if err != nil {
		handleError(c, err)
		return
	}
	Success(c, gin.H{
		"categories": categories,
		"total":      total,
	})
}

func CategoryUpdate(c *gin.Context) {
	r := models.CategoryRequest{}
	if errMsg, errKey := BindAndValidate(c, &r); errMsg != "" || errKey != "" {
		handleError(c, app.NewAppError(errMsg, errKey, ""))
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	idp := uint(id)
	r.ID = &idp

	if id == 0 {
		handleError(c, app.ErrRequired.SetKey("id"))
		return
	}
	category, err := app.CategoryUpdate(r)
	if err != nil {
		handleError(c, err)
		return
	}
	Success(c, gin.H{
		"category": category,
	})
}

func CategoryCreate(c *gin.Context) {
	r := models.CategoryRequest{}
	if errMsg, errKey := BindAndValidate(c, &r); errMsg != "" || errKey != "" {
		handleError(c, app.NewAppError(errMsg, errKey, ""))
		return
	}
	category, err := app.CategoryCreate(r)
	if err != nil {
		handleError(c, err)
		return
	}
	Success(c, gin.H{
		"category": category,
	})
}

func CategoryDelete(c *gin.Context) {
	var ids []string = c.QueryArray("ids")
	if len(ids) == 0 {
		handleError(c, app.ErrRequired.SetKey("ids"))
		return
	}
	categories, err := app.CategoryDelete(ids)
	if err != nil {
		handleError(c, err)
		return
	}
	Success(c, gin.H{
		"categories": categories,
	})
}
