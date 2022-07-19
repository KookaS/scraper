package mongodb

import (
	"errors"
	"strings"

	"scraper/src/types"
	"scraper/src/utils"

	"context"

	"fmt"

	"time"

	"path/filepath"

	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"context"
	"fmt"
	"os"
	"path/filepath"
	"scraper/src/utils"
	"time"

	"github.com/pkg/errors"
)

// InsertImage insert an image in its collection
func InsertImage(collection *mongo.Collection, image types.Image) (primitive.ObjectID, error) {
	res, err := collection.InsertOne(context.TODO(), image)
	if err != nil {
		return primitive.NilObjectID, err
	}
	insertedID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("Safecast of ObjectID did not work")
	}
	return insertedID, nil
}

// RemoveImage remove an image based on its mongodb id
func RemoveImage(collection *mongo.Collection, id primitive.ObjectID) (*int64, error) {
	query := bson.M{"_id": id}
	res, err := collection.DeleteOne(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	return &res.DeletedCount, nil
}

// RemoveImageAndFile remove an image based on its mongodb id and remove its file
func RemoveImageAndFile(collection *mongo.Collection, id primitive.ObjectID) (*int64, error) {
	image, err := FindOne[types.Image](collection, bson.M{"_id": id})
	if err != nil {
		return nil, fmt.Errorf("FindImageByID has failed: %v", err)
	}
	deletedCount, err := RemoveImage(collection, id)
	if err != nil {
		return nil, fmt.Errorf("RemoveImage has failed: %v", err)
	}
	folderDir := utils.DotEnvVariable("IMAGE_PATH")
	path := filepath.Join(folderDir, image.Origin, image.Name)
	err = os.Remove(path)
	// sometimes images can have the same file stored but are present multiple in the search request
	if err != nil && *deletedCount == 0 {
		return nil, fmt.Errorf("os.Remove has failed: %v", err)
	}
	return deletedCount, nil
}

func RemoveImagesAndFiles(mongoClient *mongo.Client, query bson.M, options *options.FindOptions) (*int64, error) {
	collectionImages := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_WANTED_COLLECTION"))
	var deletedCount int64
	images, err := FindMany[types.Image](collectionImages, query, options)
	if err != nil {
		return nil, fmt.Errorf("FindImagesIDs has failed: %v", err)
	}
	for _, image := range images {
		deletedOne, err := RemoveImageAndFile(collectionImages, image.ID)
		if err != nil {
			return nil, fmt.Errorf("RemoveImageAndFile has failed for %s: %v", image.ID.Hex(), err)
		}
		deletedCount += *deletedOne
	}
	return &deletedCount, nil
}

// UpdateImageTags add tags to an image based on its mongodb id
func UpdateImageTagsPush(collection *mongo.Collection, body types.BodyUpdateImageTagsPush) (*int64, error) {
	query := bson.M{"_id": body.ID}
	for i := 0; i < len(body.Tags); i++ {
		tag := &body.Tags[i]
		now := time.Now()
		tag.CreationDate = &now
	}
	update := bson.M{
		"$push": bson.M{
			"tags": bson.M{"$each": body.Tags},
		},
	}
	res, err := collection.UpdateOne(context.TODO(), query, update)
	if err != nil {
		return nil, fmt.Errorf("UpdateOne has failed: %v", err)
	}
	return &res.ModifiedCount, nil
}

// UpdateImageTagsPull removes specific tags from an image
func UpdateImageTagsPull(collection *mongo.Collection, body types.BodyUpdateImageTagsPull) (*int64, error) {
	query := bson.M{
		"_id":    body.ID,
		"origin": body.Origin,
	}
	update := bson.M{
		"$pull": bson.M{
			"tags": bson.M{
				"name": bson.M{
					"$in": body.Names,
				},
			},
		},
	}
	res, err := collection.UpdateOne(context.TODO(), query, update)
	if err != nil {
		return nil, err
	}
	return &res.ModifiedCount, nil
}





// InsertImageUnwanted insert an unwanted image
func InsertImageUnwanted(mongoClient *mongo.Client, body types.Image) (interface{}, error) {
	now := time.Now()
	body.CreationDate = &now
	body.Origin = strings.ToLower(body.Origin)

	// insert the unwanted image
	collectionImagesUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_UNWANTED_COLLECTION"))
	res, err := collectionImagesUnwanted.InsertOne(context.TODO(), body)
	if err != nil {
		return nil, fmt.Errorf("InsertOne has failed: %v", err)
	}
	return res.InsertedID, nil
}
