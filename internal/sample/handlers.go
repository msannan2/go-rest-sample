package sample

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func CreateNewCategory(c echo.Context) (err error) {

	u := new(Category)
	if err = c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := CreateRecord(u); err != nil {
		return echo.NewHTTPError(409, "A resource with the ID already exists.")
	}
	return c.JSON(http.StatusOK, u)
}

func CreateNewArticle(c echo.Context) (err error) {
	u := new(Article)
	if err = c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.Validate(u); err != nil {
		return err
	}
	if err := CreateArticle(u); err != nil {
		return echo.NewHTTPError(409, "A resource with the ID already exists.")
	}
	return c.JSON(http.StatusOK, u)
}

func UpdateExistingArticle(c echo.Context) (err error) {
	u := new(Article)
	id := c.Param("articleId")
	cv, err := strconv.Atoi(id)

	if err = c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.Validate(u); err != nil {
		return err
	}
	u.Id = int64(cv)
	if err := UpdateArticle(u); err != nil {
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, u)
}

func GetArticleById(c echo.Context) (err error) {

	var u *Article
	id := c.Param("articleId")
	cv, err := strconv.Atoi(id)
	if err != nil {
		return echo.ErrBadRequest
	}
	u, err = GetArticle(int64(cv))
	if err != nil {
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, u)
}

func DeleteExistingArticle(c echo.Context) (err error) {

	id := c.Param("articleId")
	cv, err := strconv.Atoi(id)
	if err != nil {
		return echo.ErrBadRequest
	}
	err = DeleteArticle(int64(cv))
	if err != nil {
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, id)
}

func GetAllArticles(c echo.Context) (err error) {

	u, err := GetArticles()
	if err != nil {
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, u)
}
