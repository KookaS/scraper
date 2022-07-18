package dynamodb

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
)

func TablesList(client *dynamodb.Client) error {
	// Build the request with its input parameters
	resp, err := client.ListTables(context.TODO(), &dynamodb.ListTablesInput{
		Limit: aws.Int32(5),
	})
	if err != nil {
		return errors.Wrap(err, "failed to list tables")
	}
	log.Println("Tables:")
	for _, tableName := range resp.TableNames {
		log.Println(tableName)
	}
	return nil
}

func TableWait(db *dynamodb.Client, tn string) error {
	w := dynamodb.NewTableExistsWaiter(db)
	err := w.Wait(context.TODO(),
		&dynamodb.DescribeTableInput{
			TableName: aws.String(tn),
		},
		2*time.Minute,
		func(o *dynamodb.TableExistsWaiterOptions) {
			o.MaxDelay = 5 * time.Second
			o.MinDelay = 5 * time.Second
		})
	if err != nil {
		return errors.Wrap(err, "timed out while waiting for table to become active")
	}
	return nil
}

func TableCreateImage(client *dynamodb.Client, tableName string) *awsTypes.TableDescription {
	out, err := client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []awsTypes.AttributeDefinition{
			{
				AttributeName: aws.String("origin"),
				AttributeType: awsTypes.ScalarAttributeTypeS, // (S | N | B) for string, number, binary
			},
			{
				AttributeName: aws.String("originID"),
				AttributeType: awsTypes.ScalarAttributeTypeS, // (S | N | B) for string, number, binary
			},
		},
		KeySchema: []awsTypes.KeySchemaElement{
			{
				AttributeName: aws.String("origin"),
				KeyType:       awsTypes.KeyTypeHash,
			},
			{
				AttributeName: aws.String("originID"),
				KeyType:       awsTypes.KeyTypeRange,
			},
		},
		TableName:   aws.String(tableName),
		BillingMode: awsTypes.BillingModePayPerRequest,
	})
	if err != nil {
		log.Fatalf("Cannot create Table `%s`, %v", tableName, err)
	}
	return out.TableDescription
}

func TableCreateUser(client *dynamodb.Client, tableName string) *awsTypes.TableDescription {
	out, err := client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []awsTypes.AttributeDefinition{
			{
				AttributeName: aws.String("origin"),
				AttributeType: awsTypes.ScalarAttributeTypeS, // (S | N | B) for string, number, binary
			},
			{
				AttributeName: aws.String("originID"),
				AttributeType: awsTypes.ScalarAttributeTypeS, // (S | N | B) for string, number, binary
			},
		},
		KeySchema: []awsTypes.KeySchemaElement{
			{
				AttributeName: aws.String("origin"),
				KeyType:       awsTypes.KeyTypeHash,
			},
			{
				AttributeName: aws.String("originID"),
				KeyType:       awsTypes.KeyTypeRange,
			},
		},
		TableName:   aws.String(tableName),
		BillingMode: awsTypes.BillingModePayPerRequest,
	})
	if err != nil {
		log.Fatalf("Cannot create Table `%s`, %v", tableName, err)
	}
	return out.TableDescription
}

func TableCreateTag(client *dynamodb.Client, tableName string) *awsTypes.TableDescription {
	out, err := client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []awsTypes.AttributeDefinition{
			{
				AttributeName: aws.String("name"),
				AttributeType: awsTypes.ScalarAttributeTypeS, // (S | N | B) for string, number, binary
			},
		},
		KeySchema: []awsTypes.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
				KeyType:       awsTypes.KeyTypeHash,
			},
		},
		TableName:   aws.String(tableName),
		BillingMode: awsTypes.BillingModePayPerRequest,
	})
	if err != nil {
		log.Fatalf("Cannot create Table `%s`, %v", tableName, err)
	}
	return out.TableDescription
}
