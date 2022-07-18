package dynamodb

import (
	"context"
	"scraper/src/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
)

type wrapperSchema interface {
	types.User | types.Tag | types.Image
}

func InsertItem[I wrapperSchema](client *dynamodb.Client, tableName string, item I) (*dynamodb.PutItemOutput, error) {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return nil, errors.Wrap(err, "failed to DynamoDB marshal Record")
	}
	av2 := map[string]awsTypes.AttributeValue{}
	for key, value := range av {
		av2[key] = value
	}

	putItemOutput, err := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av2,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to put Record to DynamoDB")
	}
	return putItemOutput, nil
}

func RemoveItem[I types.User | types.Image](client *dynamodb.Client, tableName string, item I) (*dynamodb.DeleteItemOutput, error) {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]awsTypes.AttributeValue{
			"origin": awsTypes.AttributeValueMemberS{
				Value: item.Origin,
			},
			"Title": awsTypes.AttributeValueMemberM{
				S: aws.String(""),
			},
		},
		TableName: aws.String(tableName),
	}
	
	deleteItemOutput, err := client.DeleteItem(context.TODO(), input)
	if err != nil {
		return nil, errors.Wrap(err, "Got error calling DeleteItem")
	}
	return deleteItemOutput, nil
}
