package main

import (
	"context"
	"fmt"
	"scraper/src/mongodb"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	// "scraper/src/router"
)

func main() {

	fmt.Print("Hello AWS")

	// client := mongodb.ConnectMongoDB()
	client := mongodb.ConnectDynamoDB()
	// _ = router.Router(client)

	var params *dynamodb.GetItemInput = {
		
	}
	
	client.GetItem(context.TODO())
}
