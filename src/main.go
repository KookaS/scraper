package main

import (
	"fmt"
	"time"
	// "scraper/src/mongodb"
	// "scraper/src/router"
	// "scraper/src/utils"
)

func main() {
	fmt.Print("Hello AWS")
	time.Sleep(15 * time.Minute)
	// fmt.Println(utils.DotEnvVariable("SCRAPER_DB"))
	// fmt.Println(utils.DotEnvVariable("IMAGE_PATH"))
	// fmt.Println(utils.DotEnvVariable("TAGS_UNWANTED_COLLECTION"))
	// fmt.Println(utils.DotEnvVariable("TAGS_WANTED_COLLECTION"))
	// mongoClient := mongodb.ConnectMongoDB()
	// _ = router.Router(mongoClient)
}
