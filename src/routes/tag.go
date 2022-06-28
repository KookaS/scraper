package routes

import (
	"scrapper/src/mongodb"
	"scrapper/src/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ParamsRemoveTag struct {
	ID string `uri:"id" binding:"required"`
}

func RemoveTagWanted(mongoClient *mongo.Client, params ParamsRemoveTag) (*int64, error) {
	collectionTagsWanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("TAGS_WANTED_COLLECTION"))
	tagID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return nil, err
	}
	return mongodb.RemoveTag(collectionTagsWanted, tagID)
}

func RemoveTagUnwanted(mongoClient *mongo.Client, params ParamsRemoveTag) (*int64, error) {
	collectionTagsUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("TAGS_UNWANTED_COLLECTION"))
	tagID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return nil, err
	}
	return mongodb.RemoveTag(collectionTagsUnwanted, tagID)
}
