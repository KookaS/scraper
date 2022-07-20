package dynamodb

import (
	"context"
	"os"
	"path/filepath"
	scraperTypes "scraper/src/types"
	"scraper/src/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
)

func DeleteImage(client *dynamodb.Client, tableName string, origin string, originID string) (*dynamodb.DeleteItemOutput, error) {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]awsTypes.AttributeValue{
			"origin":   awsTypes.AttributeValueMemberS{Value: origin},   // Partition Key
			"originID": awsTypes.AttributeValueMemberS{Value: originID}, // Sort Key
		},
		TableName: aws.String(tableName),
	}
	return DeleteInput(client, input)
}

func GetImage(client *dynamodb.Client, tableName string, origin string, originID string) (*scraperTypes.Image, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]awsTypes.AttributeValue{
			"origin":   awsTypes.AttributeValueMemberS{Value: origin},   // Partition Key
			"originID": awsTypes.AttributeValueMemberS{Value: originID}, // Sort Key
		},
	}
	return GetInput[scraperTypes.Image](client, input)
}

func UpdateImageCrop(client *dynamodb.Client, tableName string, image scraperTypes.Image) (*dynamodb.UpdateItemOutput, error) {
	upd := expression.
		Set(expression.Name("tags"), expression.Value(image.Tags)).
		Set(expression.Name("sizes"), expression.Value(image.Sizes))
	expr, err := expression.NewBuilder().WithUpdate(upd).Build()
	if err != nil {
		return nil, errors.Wrap(err, "Got error building expression")
	}
	return updateExpression(client, tableName, image.Origin, image.OriginID, expr)
}

func UpdateImageTagsPush(client *dynamodb.Client, tableName string, origin string, originID string, tags []scraperTypes.Tag) (*dynamodb.UpdateItemOutput, error) {
	upd := expression.Add(expression.Name("tags"), expression.Value(tags))
	expr, err := expression.NewBuilder().WithUpdate(upd).Build()
	if err != nil {
		return nil, errors.Wrap(err, "Got error building expression")
	}
	return updateExpression(client, tableName, origin, originID, expr)
}

func UpdateImageTagsPull(client *dynamodb.Client, tableName string, origin string, originID string, tags []scraperTypes.Tag) (*dynamodb.UpdateItemOutput, error) {
	upd := expression.Delete(expression.Name("tags"), expression.Value(tags))
	expr, err := expression.NewBuilder().WithUpdate(upd).Build()
	if err != nil {
		return nil, errors.Wrap(err, "Got error building expression")
	}
	return updateExpression(client, tableName, origin, originID, expr)
}

func updateExpression(client *dynamodb.Client, tableName string, origin string, originID string, expr expression.Expression) (*dynamodb.UpdateItemOutput, error){
	input := &dynamodb.UpdateItemInput{
		Key: map[string]awsTypes.AttributeValue{
			"origin":   awsTypes.AttributeValueMemberS{Value: origin},   // Partition Key
			"originID": awsTypes.AttributeValueMemberS{Value: originID}, // Sort Key
		},
		TableName:                 aws.String(tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ConditionExpression:       expr.Condition(),
	}

	updateItemOutput, err := client.UpdateItem(context.TODO(), input)
	if err != nil {
		return nil, errors.Wrap(err, "Got error calling UpdateItem")
	}
	return updateItemOutput, nil
}

func RemoveImagesAndFiles(client *dynamodb.Client, tableName string, origin string, originID string) (interface{}, error) {
	var deletedCount int64
	images, err := ScanItems[scraperTypes.Image](client, tableName)
	if err != nil {
		return nil, errors.Wrap(err, "Got error calling ScanItems")
	}

	for _, image := range images {
		deletedOne, err := RemoveImageAndFile(client, tableName, origin, originID)
		if err != nil {
			return nil, errors.Wrap(err, "Got error calling RemoveImageAndFile")
		}
		deletedCount += *deletedOne
	}
	return &deletedCount, nil
}

// RemoveImageAndFile remove an image based on its mongodb id and remove its file
func RemoveImageAndFile(client *dynamodb.Client, tableName string, origin string, originID string) (*int, error) {
	// get image for removing the file
	image, err := GetImage(client, tableName, origin, originID)
	if err != nil {
		return nil, errors.Wrap(err, "Got error calling GetImage")
	}

	// remove image in DB
	_, err = DeleteImage(client, tableName, origin, originID)
	if err != nil {
		return nil, errors.Wrap(err, "Got error calling DeleteImage")
	}

	// TODO: S3
	folderDir := utils.DotEnvVariable("IMAGE_PATH")
	path := filepath.Join(folderDir, image.Origin, image.Name)
	err = os.Remove(path)
	// sometimes images can have the same file stored but are present multiple in the search request
	// TODO: if err != nil && deleteCount == 0 
	if err != nil {
		return nil, errors.Wrap(err, "Got error calling os.Remove")
	}
	deleteCount := 1
	return &deleteCount, nil
}