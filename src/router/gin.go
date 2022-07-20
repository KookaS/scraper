package router

import (
	"net/http"

	"scraper/src/utils"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router[C *dynamodb.Client](client C) (*gin.Engine){
	router := gin.Default()
	router.Use(cors.Default())

	// routes for one image pending or wanted
	router.Static("/image/file", utils.DotEnvVariable("IMAGE_PATH"))
	router.GET("/image/:id/:collection", wrapperHandlerURI(client, FindImage))
	router.PUT("/image/tags/push", wrapperHandlerBody(client, UpdateImageTagsPush))
	router.PUT("/image/tags/pull", wrapperHandlerBody(client, UpdateImageTagsPull))
	router.PUT("/image/crop", wrapperHandlerBody(client, UpdateImageCrop))
	router.POST("/image/crop", wrapperHandlerBody(client, CreateImageCrop))
	router.POST("/image/transfer", wrapperHandlerBody(client, TransferImage))
	router.DELETE("/image", wrapperHandlerBody(client, RemoveImageAndFile))

	// routes for multiple images pending or wanted
	router.GET("/images/id/:origin/:collection", wrapperHandlerURI(client, FindImagesIDs))

	// routes for one image unwanted
	router.POST("/image/unwanted", wrapperHandlerBody(client, clientdb.InsertImageUnwanted))
	router.DELETE("/image/unwanted/:id", wrapperHandlerURI(client, RemoveImage))

	// routes for multiple images unwanted
	router.GET("/images/unwanted", wrapperHandler(client, FindImagesUnwanted))

	// routes for one tag
	router.POST("/tag/wanted", wrapperHandlerBody(client, clientdb.InsertTagWanted))
	router.POST("/tag/unwanted", wrapperHandlerBody(client, clientdb.InsertTagUnwanted))
	router.DELETE("/tag/wanted/:id", wrapperHandlerURI(client, RemoveTagWanted))
	router.DELETE("/tag/unwanted/:id", wrapperHandlerURI(client, RemoveTagUnwanted))

	// routes for multiple tags
	router.GET("/tags/wanted", wrapperHandler(client, clientdb.TagsWanted))
	router.GET("/tags/unwanted", wrapperHandler(client, clientdb.TagsUnwanted))

	// routes for one user unwanted
	router.POST("/user/unwanted", wrapperHandlerBody(client, clientdb.InsertUserUnwanted))
	router.DELETE("/user/unwanted/:id", wrapperHandlerURI(client, RemoveUserUnwanted))

	// routes for multiplt users unwanted
	router.GET("/users/unwanted", wrapperHandler(client, clientdb.UsersUnwanted))

	// routes for scraping the internet
	router.POST("/search/flickr/:quality", wrapperHandlerURI(client, SearchPhotosFlickr))
	router.POST("/search/unsplash/:quality", wrapperHandlerURI(client, SearchPhotosUnsplash))
	router.POST("/search/pexels/:quality", wrapperHandlerURI(client, SearchPhotosPexels))

	// start the backend
	router.Run("localhost:8080")
	return router
}

// wrapper for the response with argument
func wrapperResponseArg[C *dynamodb.Client, A any, R any](c *gin.Context, f func(client C, arg A) (R, error), client C, arg A) {
	res, err := f(client, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// wrapper for the response
func wrapperResponse[C *dynamodb.Client, R any](c *gin.Context, f func(client C) (R, error), client C) {
	res, err := f(client)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// wrapper for the ginHandler with body with collectionName
func wrapperHandlerBody[C *dynamodb.Client, B any, R any](client C, f func(client C, body B) (R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body B
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		wrapperResponseArg(c, f, client, body)
	}
}

// wrapper for the ginHandler with URI
func wrapperHandlerURI[C *dynamodb.Client, P any, R any](client C, f func(client C, params P) (R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params P
		if err := c.ShouldBindUri(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		wrapperResponseArg(c, f, client, params)
	}
}

// wrapper for the ginHandler
func wrapperHandler[C *dynamodb.Client, R any](client C, f func(client C) (R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		wrapperResponse(c, f, client)
	}
}
