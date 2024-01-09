package api

import (
	"net/http"
	"strconv"

	"github.com/Sohbetbackend/selfProject/internal/app"
	"github.com/Sohbetbackend/selfProject/internal/models"
	"github.com/gin-gonic/gin"
)

func AuthorRoutes(api *gin.RouterGroup) {
	authorRoutes := api.Group("/authors")
	{
		authorRoutes.GET("", AuthorList)
		authorRoutes.GET(":id", AuthorsDetail)
		authorRoutes.POST("", AuthorCreate)
		authorRoutes.PUT(":id", AuthorUpdate)
		authorRoutes.DELETE("", AuthorDelete)
	}
}

func AuthorList(c *gin.Context) {
	r := models.AuthorsFilterRequest{}
	if errMsg, errKey := BindAndValidate(c, &r); errMsg != "" || errKey != "" {
		handleError(c, app.NewAppError(errMsg, errKey, ""))
		return
	}
	if r.Limit == nil {
		r.Limit = new(int)
		*r.Limit = 30
	}
	if r.Offset == nil {
		r.Offset = new(int)
		*r.Offset = 0
	}
	authors, total, err := app.AuthorsList(r)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, authors)
}

func AuthorsDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, err)
		return
	}
	idu := uint(id)
	if id == 0 {
		handleError(c, err)
		return
	}

	args := models.AuthorsFilterRequest{
		ID: &idu,
	}
	res, err := app.AuthorDetail(args)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func AuthorUpdate(c *gin.Context) {
	r := models.AuthorRequest{}
	if errMsg, errKey := BindAndValidate(c, &r); errMsg != "" || errKey != "" {
		handleError(c, app.NewAppError(errMsg, errKey, ""))
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, err)
		return
	}
	idp := uint(id)
	r.ID = &idp

	if id == 0 {
		handleError(c, app.ErrRequired.SetKey("id"))
		return
	}
	author, err := app.AuthorsUpdate(r)
	if err != nil {
		handleError(c, err)
		return
	}
	Success(c, gin.H{
		"author": author,
	})
}

func AuthorCreate(c *gin.Context) {
	r := models.AuthorRequest{}
	if errMsg, errKey := BindAndValidate(c, &r); errMsg != "" || errKey != "" {
		handleError(c, app.NewAppError(errMsg, errKey, ""))
		return
	}
	author, err := app.AuthorsCreate(r)
	if err != nil {
		handleError(c, err)
		return
	}
	Success(c, gin.H{
		"author": author,
	})
}

func AuthorDelete(c *gin.Context) {
	var ids []string = c.QueryArray("ids")
	if len(ids) == 0 {
		handleError(c, app.ErrRequired.SetKey("ids"))
		return
	}
	authors, err := app.AuthorsDelete(ids)
	if err != nil {
		handleError(c, err)
		return
	}
	Success(c, gin.H{
		"authors": authors,
	})
}
