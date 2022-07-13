package types

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClientSchema interface {
	*mongo.Client | *dynamodb.Client
}
