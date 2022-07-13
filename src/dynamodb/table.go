package dynamodb

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func Tables(client *dynamodb.Client) {
	// Build the request with its input parameters
	resp, err := client.ListTables(context.TODO(), &dynamodb.ListTablesInput{
		Limit: aws.Int32(5),
	})
	if err != nil {
		log.Fatalf("failed to list tables, %v", err)
	}
	log.Println("Tables:")
	for _, tableName := range resp.TableNames {
		log.Println(tableName)
	}
}