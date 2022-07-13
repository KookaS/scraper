package router

import (
	"net/http"
	"scraper/src/mongodb"
	"scraper/src/utils"
	"scraper/src/types"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router[C types.ClientSchema](mongoClient C) (*gin.Engine){
	router := gin.Default()
	router.Use(cors.Default())

	// routes for one image pending or wanted
	router.Static("/image/file", utils.DotEnvVariable("IMAGE_PATH"))
	router.GET("/image/:id/:collection", wrapperHandlerURI(mongoClient, FindImage))
	router.PUT("/image/tags/push", wrapperHandlerBody(mongoClient, UpdateImageTagsPush))
	router.PUT("/image/tags/pull", wrapperHandlerBody(mongoClient, UpdateImageTagsPull))
	router.PUT("/image/crop", wrapperHandlerBody(mongoClient, mongodb.UpdateImageCrop))
	router.POST("/image/crop", wrapperHandlerBody(mongoClient, mongodb.CreateImageCrop))
	router.POST("/image/transfer", wrapperHandlerBody(mongoClient, mongodb.TransferImage))
	router.DELETE("/image/:id", wrapperHandlerURI(mongoClient, RemoveImageAndFile))

	// routes for multiple images pending or wanted
	router.GET("/images/id/:origin/:collection", wrapperHandlerURI(mongoClient, FindImagesIDs))

	// routes for one image unwanted
	router.POST("/image/unwanted", wrapperHandlerBody(mongoClient, mongodb.InsertImageUnwanted))
	router.DELETE("/image/unwanted/:id", wrapperHandlerURI(mongoClient, RemoveImage))

	// routes for multiple images unwanted
	router.GET("/images/unwanted", wrapperHandler(mongoClient, FindImagesUnwanted))

	// routes for one tag
	router.POST("/tag/wanted", wrapperHandlerBody(mongoClient, mongodb.InsertTagWanted))
	router.POST("/tag/unwanted", wrapperHandlerBody(mongoClient, mongodb.InsertTagUnwanted))
	router.DELETE("/tag/wanted/:id", wrapperHandlerURI(mongoClient, RemoveTagWanted))
	router.DELETE("/tag/unwanted/:id", wrapperHandlerURI(mongoClient, RemoveTagUnwanted))

	// routes for multiple tags
	router.GET("/tags/wanted", wrapperHandler(mongoClient, mongodb.TagsWanted))
	router.GET("/tags/unwanted", wrapperHandler(mongoClient, mongodb.TagsUnwanted))

	// routes for one user unwanted
	router.POST("/user/unwanted", wrapperHandlerBody(mongoClient, mongodb.InsertUserUnwanted))
	router.DELETE("/user/unwanted/:id", wrapperHandlerURI(mongoClient, RemoveUserUnwanted))

	// routes for multiplt users unwanted
	router.GET("/users/unwanted", wrapperHandler(mongoClient, mongodb.UsersUnwanted))

	// routes for scraping the internet
	router.POST("/search/flickr/:quality", wrapperHandlerURI(mongoClient, SearchPhotosFlickr))
	router.POST("/search/unsplash/:quality", wrapperHandlerURI(mongoClient, SearchPhotosUnsplash))
	router.POST("/search/pexels/:quality", wrapperHandlerURI(mongoClient, SearchPhotosPexels))

	// start the backend
	router.Run("localhost:8080")
	return router
}

// wrapper for the response with argument
func wrapperResponseArg[C types.ClientSchema, A any, R any](c *gin.Context, f func(mongo C, arg A) (R, error), mongo C, arg A) {
	res, err := f(mongo, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// wrapper for the response
func wrapperResponse[C types.ClientSchema, R any](c *gin.Context, f func(mongo C) (R, error), mongo C) {
	res, err := f(mongo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// wrapper for the ginHandler with body with collectionName
func wrapperHandlerBody[C types.ClientSchema, B any, R any](mongoClient C, f func(mongo C, body B) (R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body B
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		wrapperResponseArg(c, f, mongoClient, body)
	}
}

// wrapper for the ginHandler with URI
func wrapperHandlerURI[C types.ClientSchema, P any, R any](mongoClient C, f func(mongo C, params P) (R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params P
		if err := c.ShouldBindUri(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}
		wrapperResponseArg(c, f, mongoClient, params)
	}
}

// wrapper for the ginHandler
func wrapperHandler[C types.ClientSchema, R any](mongoClient C, f func(mongo C) (R, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		wrapperResponse(c, f, mongoClient)
	}
}
