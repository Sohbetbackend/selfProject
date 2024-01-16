package api

import (
	"net/http"
	"strconv"

	"github.com/Sohbetbackend/selfProject/internal/app"
	"github.com/Sohbetbackend/selfProject/internal/models"
	"github.com/gin-gonic/gin"
)

func BookRoutes(api *gin.RouterGroup) {
	bookRoutes := api.Group("/books")
	{
		bookRoutes.GET("", BookList)
		bookRoutes.GET(":id", BooksDetail)
		bookRoutes.POST("", BookCreate)
		bookRoutes.PUT(":id", BookUpdate)
		bookRoutes.DELETE("", BookDelete)
	}
}

func BookList(c *gin.Context) {
	r := models.BookFilterRequest{}
	if errMsg, errKey := BindAndValidate(c, &r); errMsg != "" || errKey != "" {
		handleError(c, app.NewAppError(errMsg, errKey, ""))
		return
	}
	books, total, err := app.BooksList(r)
	if err != nil {
		handleError(c, err)
		return
	}
	c.Header("X-Total-Count", strconv.Itoa(total))
	c.JSON(http.StatusOK, books)
}

func BooksDetail(c *gin.Context) {
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

	args := models.BookFilterRequest{
		ID: &idu,
	}
	res, err := app.BookDetail(args)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func BookUpdate(c *gin.Context) {
	r := models.BookRequest{}
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
	book, err := app.BookUpdate(r)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, book)
}

func BookCreate(c *gin.Context) {
	r := models.BookRequest{}
	if errMsg, errKey := BindAndValidate(c, &r); errMsg != "" || errKey != "" {
		handleError(c, app.NewAppError(errMsg, errKey, ""))
		return
	}

	files, err := handleFilesUpload(c, "files", "books")
	if err != nil {
		app.NewAppError(err.Error(), "files", "")
		return
	}

	if files != nil {
		r.Files = &files
	}

	book, err := app.BookCreate(r)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, book)
}

func BookDelete(c *gin.Context) {
	var ids []string = c.QueryArray("ids")
	if len(ids) == 0 {
		handleError(c, app.ErrRequired.SetKey("ids"))
		return
	}
	books, err := app.BookDelete(ids)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, books)
}
