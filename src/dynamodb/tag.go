package dynamodb

import (
	scraperTypes "scraper/src/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func DeleteTag(client *dynamodb.Client, tableName string, name string) (*dynamodb.DeleteItemOutput, error) {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]awsTypes.AttributeValue{
			"name": awsTypes.AttributeValueMemberS{Value: name}, // Partition Key
		},
		TableName: aws.String(tableName),
	}
	return DeleteInput(client, input)
}

func GetTag(client *dynamodb.Client, tableName string, name string) (*scraperTypes.Tag, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]awsTypes.AttributeValue{
			"name": awsTypes.AttributeValueMemberS{Value: name}, // Partition Key
		},
	}
	return GetInput[scraperTypes.Tag](client, input)
}
