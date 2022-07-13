package utils

import (
	"fmt"
	"scraper/src/types"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.mongodb.org/mongo-driver/mongo"
)

func ImagesCollection[C types.ClientSchema](client *dynamodb.Client, collection string) (*mongo.Collection, error) {
	switch collection {
	case "wanted":
		return client.Database(DotEnvVariable("SCRAPER_DB")).Collection(DotEnvVariable("IMAGES_WANTED_COLLECTION")), nil
	case "pending":
		return client.Database(DotEnvVariable("SCRAPER_DB")).Collection(DotEnvVariable("IMAGES_PENDING_COLLECTION")), nil
	case "unwanted":
		return client.Database(DotEnvVariable("SCRAPER_DB")).Collection(DotEnvVariable("IMAGES_UNWANTED_COLLECTION")), nil
	default:
		return nil, fmt.Errorf("`%s` does not exist for selecting the images collection. Choose `%s`, `%s` or `%s`",
			collection, "wanted", "pending", "unwanted")
	}
}
