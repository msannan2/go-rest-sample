package sample

// Constants related to the sample REST service routes
const (
	ApiBase                = "/api/v1"
	ApiPathParamArticleId  = "/:articleId"
	ApiPathParamCategoryId = "/:categoryId"
)

const ApiEndpointArticles = ApiBase + "/articles"
const ApiEndpointCategories = ApiBase + "/categories"
