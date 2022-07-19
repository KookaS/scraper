package dynamodb

import (
	"context"
	scraperTypes "scraper/src/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
)

type wrapperSchema interface {
	scraperTypes.User | scraperTypes.Tag | scraperTypes.Image
}

func InsertItem[I wrapperSchema](client *dynamodb.Client, tableName string, item I) (*dynamodb.PutItemOutput, error) {
	// transform structure to map for dynamoDB format
	data, err := attributevalue.MarshalMap(item)
	if err != nil {
		return nil, errors.Wrap(err, "failed to MarshalMap")
	}

	putItemOutput, err := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      data,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to put Record to DynamoDB")
	}
	return putItemOutput, nil
}

func DeleteInput(client *dynamodb.Client, input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	deleteItemOutput, err := client.DeleteItem(context.TODO(), input)
	if err != nil {
		return nil, errors.Wrap(err, "Got error calling DeleteItem")
	}
	return deleteItemOutput, nil
}

// getInput retrieve an item
func GetInput[I wrapperSchema](client *dynamodb.Client, input *dynamodb.GetItemInput) (*I, error) {
	result, err := client.GetItem(context.TODO(), input)
	if err != nil {
		return nil, errors.Wrap(err, "Got error calling GetItem")
	}
	if result.Item == nil {
		return nil, errors.Wrap(err, "Could not find Item")
	}

	// writes the result into a struct
	var item I
	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal Item")
	}
	return &item, nil
}

// GetItems retrieve all the items from a table
func ScanItems[I wrapperSchema](client *dynamodb.Client, tableName string) (*[]I, error) {
	data, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Scan")
	}

	var items []I
	err = attributevalue.UnmarshalListOfMaps(data.Items, &items)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to UnmarshalListOfMaps")
	}
	return &items, nil
}

func QueryItems[I wrapperSchema](client *dynamodb.Client, tableName string) (*[]I, error) {
	data, err := client.Scan(context.TODO(), &dynamodb.QueryInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Scan")
	}

	var items []I
	err = attributevalue.UnmarshalListOfMaps(data.Items, &items)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to UnmarshalListOfMaps")
	}
	return &items, nil
}
