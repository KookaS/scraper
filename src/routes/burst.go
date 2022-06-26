package routes

import (
	"os"
	"path/filepath"
	"scrapper/src/mongodb"
	"scrapper/src/utils"
	"sort"

	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

func SearchPhotosBurst(mongoClient *mongo.Client) (interface{}, error) {

	accessToken , err := loginBurst()
	if err != nil {
		return nil, err
	}
	// return accessToken, nil

	// If path is already a directory, MkdirAll does nothing and returns nil
	folderDir := utils.DotEnvVariable("IMAGE_PATH")
	origin := "unsplash"
	err = os.MkdirAll(filepath.Join(folderDir, origin), os.ModePerm)
	if err != nil {
		return nil, err
	}

	// collectionBurst := mongoClient.Database(utils.DotEnvVariable("SCRAPPER_DB")).Collection(utils.DotEnvVariable("UNSPLASH_COLLECTION"))

	unwantedTags, err := mongodb.TagsUnwantedNames(mongoClient)
	if err != nil {
		return nil, err
	}
	sort.Strings(unwantedTags)

	wantedTags, err := mongodb.TagsWantedNames(mongoClient)
	if err != nil {
		return nil, err
	}
	sort.Strings(wantedTags)

	for _, wantedTag := range wantedTags {
		page := 1

		body, err := searchPhotosPerPageBurst(*accessToken, wantedTag, page)
		if err != nil {
			return nil, err
		}
		return utils.ToJSON(body), nil
	}
	return nil, nil
}

func loginBurst() (*string, error) {
	r := &Request{
		Host: fmt.Sprintf("https://burst.myshopify.com/admin/oauth/authorize?client_id=%s", utils.DotEnvVariable("BURST_PUBLIC_KEY")),
		Args: map[string]string{},
		Header: map[string][]string{},
	}
	fmt.Println(r.URL())

	body, err := r.Execute()
	if err != nil {
		return nil, err
	}
	accessToken := fmt.Sprint(body)
	return &accessToken, nil
}

// https://shopify.dev/api/admin-rest
func searchPhotosPerPageBurst(acessToken string, tag string, page int) (interface{}, error) {
	version := "2022-04"
	r := &Request{
		Host: fmt.Sprintf("https://burst.myshopify.com/admin/api/%s/inventory_items.json?", version),
		Args: map[string]string{},
		Header: map[string][]string{
			"X-Shopify-Access-Token": {acessToken},
		},
	}
	fmt.Println(r.URL())

	body, err := r.Execute()
	if err != nil {
		return nil, err
	}
	return body, nil

	// var searchPerPage unsplash.PhotoSearchResult
	// err = json.Unmarshal(body, &searchPerPage)
	// if err != nil {
	// 	return nil, err
	// }
	// return &searchPerPage, nil
}
