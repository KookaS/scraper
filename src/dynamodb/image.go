package dynamodb

import (
	"context"
	scraperTypes "scraper/src/types"

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
	// tags, err := attributevalue.MarshalWithOptions(image.Tags)
	// sizes , err := attributevalue.MarshalWithOptions(image.Sizes)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "failed to MarshalWithOptions")
	// }
	// input := &dynamodb.UpdateItemInput{
	// 	ExpressionAttributeValues: map[string]awsTypes.AttributeValue{
	// 		":tags": tags,
	// 		":sizes": sizes,
	// 	},
	// 	TableName: aws.String(tableName),
	// 	Key: map[string]awsTypes.AttributeValue{
	// 		"origin":   awsTypes.AttributeValueMemberS{Value: image.Origin},   // Partition Key
	// 		"originID": awsTypes.AttributeValueMemberS{Value: image.OriginID}, // Sort Key
	// 	},
	// 	ReturnValues:     "UPDATED_NEW",
	// 	UpdateExpression: aws.String("set Tags = :tags, Sizes = :sizes"),
	// }
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
