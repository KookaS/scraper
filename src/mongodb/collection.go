package mongodb

import (
	"fmt"
	"scraper/src/types"
	"scraper/src/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

func ImagesCollection[C types.ClientSchema](client *mongo.Client, collection string) (*mongo.Collection, error) {
	switch collection {
	case "wanted":
		return client.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_WANTED_COLLECTION")), nil
	case "pending":
		return client.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_PENDING_COLLECTION")), nil
	case "unwanted":
		return client.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_UNWANTED_COLLECTION")), nil
	default:
		return nil, fmt.Errorf("`%s` does not exist for selecting the images collection. Choose `%s`, `%s` or `%s`",
			collection, "wanted", "pending", "unwanted")
	}
}
