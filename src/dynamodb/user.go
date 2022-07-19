package dynamodb

import (
	scraperTypes "scraper/src/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func DeleteUser(client *dynamodb.Client, tableName string, origin string, originID string) (*dynamodb.DeleteItemOutput, error) {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]awsTypes.AttributeValue{
			"origin":   awsTypes.AttributeValueMemberS{Value: origin},   // Partition Key
			"originID": awsTypes.AttributeValueMemberS{Value: originID}, // Sort Key
		},
		TableName: aws.String(tableName),
	}
	return DeleteInput(client, input)
}

func GetUser(client *dynamodb.Client, tableName string, origin string, originID string) (*scraperTypes.User, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]awsTypes.AttributeValue{
			"origin":   awsTypes.AttributeValueMemberS{Value: origin},   // Partition Key
			"originID": awsTypes.AttributeValueMemberS{Value: originID}, // Sort Key
		},
	}
	return GetInput[scraperTypes.User](client, input)
}
