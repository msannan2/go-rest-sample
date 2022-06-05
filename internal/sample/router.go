package sample

import (
	"github.com/labstack/echo/v4"
)

// loadRestRoutes loads all the routes of this service into the specified REST API router.
func loadRestRoutes(e *echo.Echo) {
	// e.POST("/signup", Signup)
	// e.POST("/login", h.Login)
	// e.POST("/follow/:id", h.Follow)
	// e.POST("/posts", h.CreatePost)
	e.POST(ApiEndpointArticles, CreateNewArticle)
	e.GET(ApiEndpointArticles+ApiPathParamArticleId, GetArticleById)
	e.GET(ApiEndpointArticles, GetAllArticles)
	e.PUT(ApiEndpointArticles+ApiPathParamArticleId, UpdateExistingArticle)
	e.DELETE(ApiEndpointArticles+ApiPathParamArticleId, DeleteExistingArticle)
}
