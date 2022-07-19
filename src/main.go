package main

import (
	"log"
	"scraper/src/dynamodb"
	"scraper/src/utils"
	// "scraper/src/router"
)

func main() {

	// client := mongodb.ConnectMongoDB()
	client := dynamodb.ConnectDynamoDB("us-east-1")
	_ = dynamodb.TableCreateImage(client, utils.DotEnvVariable("IMAGES_WANTED_COLLECTION"))
	_ = dynamodb.TableCreateImage(client, utils.DotEnvVariable("IMAGES_PENDING_COLLECTION"))
	_ = dynamodb.TableCreateImage(client, utils.DotEnvVariable("IMAGES_UNWANTED_COLLECTION"))
	_ = dynamodb.TableCreateImage(client, utils.DotEnvVariable("USERS_UNWANTED_COLLECTION"))
	_ = dynamodb.TableCreateImage(client, utils.DotEnvVariable("TAGS_UNWANTED_COLLECTION"))
	_ = dynamodb.TableCreateImage(client, utils.DotEnvVariable("TAGS_WANTED_COLLECTION"))
	err := dynamodb.TablesList(client)
	log.Fatal(err)

	// _ = router.Router(client)
}
