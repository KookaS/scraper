package router

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	scraperDynamoDB "scraper/src/dynamodb"
	"scraper/src/types"
	scraperTypes "scraper/src/types"
	"scraper/src/utils"

	awsDynamoDB "github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type ParamsFindImagesIDs struct {
	Origin     string `uri:"origin" binding:"required"`
	Collection string `uri:"collection" binding:"required"`
}

// FindImagesIDs get all the IDs of an image collection
func FindImagesIDs(client ..., params ParamsFindImagesIDs) ([]types.Image, error) {
	collectionImages, err := utils.ImagesCollection(client, params.Collection)
	if err != nil {
		return nil, err
	}
	query := bson.M{"origin": params.Origin}
	options := options.Find().SetProjection(bson.M{"_id": 1})
	return mongodb.FindMany[types.Image](collectionImages, query, options)
}

// FindImagesUnwanted get all the unwanted images
func FindImagesUnwanted(mongoClient *mongo.Client) ([]types.Image, error) {
	collectionImagesUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_UNWANTED_COLLECTION"))
	// no options needed because not much is stored for unwanted images
	return mongodb.FindMany[types.Image](collectionImagesUnwanted, bson.M{})
}

// Body for the RemoveImage request
type ParamsRemoveImage struct {
	ID string `uri:"id" binding:"required"`
}

// RemoveImageAndFile removes in db and file of a pending image
func RemoveImageAndFile(mongoClient *mongo.Client, params ParamsRemoveImage) (*int64, error) {
	collectionImagesPending := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_PENDING_COLLECTION"))
	imageID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return nil, err
	}
	return mongodb.RemoveImageAndFile(collectionImagesPending, imageID)
}

// RemoveImage removes in db an unwanted image
func RemoveImage(mongoClient *mongo.Client, params ParamsRemoveImage) (*int64, error) {
	collectionImagesUnwanted := mongoClient.Database(utils.DotEnvVariable("SCRAPER_DB")).Collection(utils.DotEnvVariable("IMAGES_UNWANTED_COLLECTION"))
	imageID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return nil, err
	}
	return mongodb.RemoveImage(collectionImagesUnwanted, imageID)
}







type ParamsFindImage struct {
	TableName string           `json:"tableName,omitempty"`
	Origin    string           `json:"origin,omitempty"`
	OriginID  string           `json:"originID,omitempty"`
}
// FindImage get a specific image
func FindImage(client *awsDynamoDB.Client, params ParamsFindImage) (*types.Image, error) {
	return scraperDynamoDB.GetImage(client, body.tableName, body.Origin, body.OriginID)
}

type BodyUpdateImageTagsPush struct {
	TableName string           `json:"tableName,omitempty"`
	Origin    string           `json:"origin,omitempty"`
	OriginID  string           `json:"originID,omitempty"`
	Tags   []scraperTypes.Tag              `json:"tags,omitempty"`
}
// UpdateImageTagsPush add tags to a pending image
func UpdateImageTagsPush(client *awsDynamoDB.Client, body BodyUpdateImageTagsPush) (interface{}, error) {
	updateImageOutput, err := scraperDynamoDB.UpdateImageTagsPush(client, body.TableName, body.Origin, body.OriginID, body.Tags[])
	if err != nil {
		return nil, fmt.Errorf("UpdateImageTags has failed: %v", err)
	}
	return updateImageOutput.ResultMetadata, nil
}

type BodyUpdateImageTagsPull struct {
	TableName string           `json:"tableName,omitempty"`
	Origin    string           `json:"origin,omitempty"`
	OriginID  string           `json:"originID,omitempty"`
	Tags   []scraperTypes.Tag              `json:"tags,omitempty"`
}
// UpdateImageTagsPush remove tags to a pending image
func UpdateImageTagsPull(client *awsDynamoDB.Client, body BodyUpdateImageTagsPull) (interface{}, error) {
	updateImageOutput, err := scraperDynamoDB.UpdateImageTagsPull(client, body.TableName, body.Origin, body.OriginID, body.Tags[])
	if err != nil {
		return nil, fmt.Errorf("UpdateImageTags has failed: %v", err)
	}
	return updateImageOutput.ResultMetadata, nil
}

type BodyImageCrop struct {
	TableName string           `json:"tableName,omitempty"`
	Origin    string           `json:"origin,omitempty"`
	OriginID  string           `json:"originID,omitempty"`
	Box       scraperTypes.Box `json:"box,omitempty"`
	File      []byte           `json:"file,omitempty"`
}
// UpdateImageFile update the image with its tags when it is cropped
func UpdateImageCrop(client *awsDynamoDB.Client, body BodyImageCrop) (interface{}, error) {
	// generate new image sizes and tags boxes
	image, err := getNewBoxes(collectionImagesPending, body)
	if err != nil {
		return nil, fmt.Errorf("getNewBoxes has failed: %v", err)
	}

	// replace in db and file of the updated image
	err = replaceImageFile(client, image, body.File)
	if err != nil {
		return nil, fmt.Errorf("replaceImageFile has failed: %v", err)
	}
	
	// update the current image
	updateImageOutput, err := scraperDynamoDB.UpdateImageCrop(client, body.TableName, image)
	if err != nil {
		return nil, fmt.Errorf("UpdateImageCrop has failed: %v", err)
	}
	return updateImageOutput.ResultMetadata, nil
}

// CreateImageCrop update the image with its tags when it is cropped
func CreateImageCrop(client *awsDynamoDB.Client, body BodyImageCrop) (interface{}, error) {
	// generate new image sizes and tags boxes
	image, err := getNewBoxes(client, body.TableName, body.Origin, body.OriginID, body.Box)
	if err != nil {
		return nil, fmt.Errorf("getNewBoxes has failed: %v", err)
	}

	// add the current date and time to the originID
	image.OriginID = fmt.Sprintf("%s_%s.%s", image.OriginID, time.Now().Format(time.RFC3339), image.Extension)

	// replace in db and file of the updated image
	err = replaceImageFile(client, image, body.File)
	if err != nil {
		return nil, fmt.Errorf("replaceImageFile has failed: %v", err)
	}

	// create a new image in table
	putItemOutput, err := scraperDynamoDB.InsertItem(client, body.TableName, image)
	if err != nil {
		return nil, fmt.Errorf("InsertItem has failed: %v", err)
	}
	return putItemOutput.ResultMetadata, nil
}

type BodyTransferImage struct {
	Origin    string           `json:"origin,omitempty"`
	OriginID  string           `json:"originID,omitempty"`
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
}
func TransferImage(client *awsDynamoDB.Client, body types.BodyTransferImage) (interface{}, error) {
	collectionImagesFrom, err := utils.ImagesCollection(mongoClient, body.From)
	if err != nil {
		return nil, err
	}
	collectionImagesTo, err := utils.ImagesCollection(mongoClient, body.To)
	if err != nil {
		return nil, err
	}
	query := bson.M{"originID": body.OriginID}
	image, err := FindOne[types.Image](collectionImagesFrom, query)
	if err != nil {
		return nil, fmt.Errorf("FindOne[Image] has failed: %v", err)
	}
	image.ID = primitive.NilObjectID
	res, err := collectionImagesTo.InsertOne(context.TODO(), *image)
	if err != nil {
		return nil, fmt.Errorf("InsertOne has failed: %v", err)
	}
	_, err = collectionImagesFrom.DeleteOne(context.TODO(), query)
	if err != nil {
		return nil, fmt.Errorf("DeleteOne has failed: %v", err)
	}
	return res.InsertedID, nil
}



func replaceImageFile(client *awsDynamoDB.Client, imageReplace *scraperTypes.Image, imageFile []byte) error {
	// replace or create the file
	folderDir := utils.DotEnvVariable("IMAGE_PATH")
	path := filepath.Join(folderDir, imageReplace.Origin, imageReplace.Name)
	return os.WriteFile(path, imageFile, 0644)
}

func getNewBoxes(client *awsDynamoDB.Client, tableName string, origin string, originID string, box scraperTypes.Box) (*scraperTypes.Image, error) {
	imageFound, err := scraperDynamoDB.GetImage(client, tableName, origin, originID)
	if err != nil {
		return nil, errors.Wrap(err, "GetImage failed")
	}

	// new size creation
	creationDate := time.Now().Format(time.RFC3339)
	size := scraperTypes.ImageSize{
		CreationDate: creationDate,
		Box:          box, // absolute position
	}
	imageFound.Sizes = append(imageFound.Sizes, size)

	i := 0
	for {
		if i >= len(imageFound.Tags) {
			break
		}
		tag := imageFound.Tags[i]
		if (scraperTypes.Box{}) != tag.Origin.Box {
			// relative position of tags
			tlx := *tag.Origin.Box.X
			tly := *tag.Origin.Box.Y
			width := *tag.Origin.Box.Width
			height := *tag.Origin.Box.Height

			// box outside on the image right
			if tlx > *box.X+*box.Width {
				// last element removed
				if i == len(imageFound.Tags)-1 {
					imageFound.Tags = imageFound.Tags[:i]
				} else { // not last element removed
					imageFound.Tags = append(imageFound.Tags[:i], imageFound.Tags[i+1:]...)
				}
				continue
			}
			// box left outside on the image left
			if tlx < *box.X {
				// box outside on the image left
				if tlx+width < *box.X {
					width = 0
				} else { // box right inside the image
					width = width - *box.X + tlx
				}
				tlx = *box.X
			} else { // box left inside image
				// box right outside on the image right
				if tlx+width > *box.X+*box.Width {
					width = *box.X + *box.Width - tlx
				}
				tlx = tlx - *box.X
			}
			// box width too small
			if width < 50 {
				// last element removed
				if i == len(imageFound.Tags)-1 {
					imageFound.Tags = imageFound.Tags[:i]
				} else { // not last element removed
					imageFound.Tags = append(imageFound.Tags[:i], imageFound.Tags[i+1:]...)
				}
				continue
			}

			// box outside at the image bottom
			if tly > *box.Y+*box.Height {
				// last element removed
				if i == len(imageFound.Tags)-1 {
					imageFound.Tags = imageFound.Tags[:i]
				} else { // not last element removed
					imageFound.Tags = append(imageFound.Tags[:i], imageFound.Tags[i+1:]...)
				}
				continue
			}
			// box top outside on the image top
			if tly < *box.Y {
				// box outside on the image top
				if tly+height < *box.Y {
					height = 0
				} else { // box bottom inside the image
					height = height - *box.Y + tly
				}
				tly = *box.Y
			} else { // box top inside image
				// box bottom outside on the image bottom
				if tly+height > *box.Y+*box.Height {
					height = *box.Y + *box.Height - tly
				}
				tly = tly - *box.Y
			}
			// box height too small
			if height < 50 {
				// last element removed
				if i == len(imageFound.Tags)-1 {
					imageFound.Tags = imageFound.Tags[:i]
				} else { // not last element removed
					imageFound.Tags = append(imageFound.Tags[:i], imageFound.Tags[i+1:]...)
				}
				continue
			}

			// set the new relative reference to the newly cropped image
			tag.Origin.ImageSizeDate = creationDate
			tag.Origin.Box.X = &tlx
			tag.Origin.Box.Y = &tly
			tag.Origin.Box.Width = &width
			tag.Origin.Box.Height = &height
		}
		i++
	}
	return imageFound, nil
}
