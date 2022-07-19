package types

// the fields in `` are for mapping one object to another with the same structure
// bson is for MongoDB
// json is for JSON used in HTTP requests
// dynamodbav is for DynamoDB

// Structure for an image strored in DB
type Image struct {
	// ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`           // mongodb default id
	Origin       string      `bson:"origin,omitempty" json:"origin,omitempty" dynamodbav:"origin,omitempty"`       // original werbsite
	OriginID     string      `bson:"originID,omitempty" json:"originID,omitempty" dynamodbav:"originID,omitempty"` // id from original website
	User         User        `bson:"user,omitempty" json:"user,omitempty" dynamodbav:"user,omitempty"`
	Extension    string      `bson:"extension,omitempty" json:"extension,omitempty" dynamodbav:"extension,omitempty"` // type of file
	Name         string      `bson:"name,omitempty" json:"name,omitempty" dynamodbav:"name,omitempty"`                // name <originID>.<extension>
	Sizes         []ImageSize `bson:"sizes,omitempty" json:"sizes,omitempty" dynamodbav:"sizes,omitempty"`                // size cropping history
	Title        string      `bson:"title,omitempty" json:"title,omitempty" dynamodbav:"title,omitempty"`
	Description  string      `bson:"description,omitempty" json:"description,omitempty" dynamodbav:"description,omitempty"` // decription of image
	License      string      `bson:"license,omitempty" json:"license,omitempty" dynamodbav:"license,omitempty"`             // type of public license
	CreationDate string      `bson:"creationDate,omitempty" json:"creationDate,omitempty" dynamodbav:"creationDate,omitempty"`
	Tags         []Tag       `bson:"tags,omitempty" json:"tags,omitempty" dynamodbav:"tags,omitempty"`
}

func (item Image) GetOrigin() string   { return item.Origin }
func (item Image) GetOriginID() string { return item.OriginID }

// Structure for a tag strored in DB
type Tag struct {
	// ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"` // mongodb default id
	Name         string    `bson:"name,omitempty" json:"name,omitempty" dynamodbav:"name,omitempty"`
	CreationDate string    `bson:"creationDate,omitempty" json:"creationDate,omitempty" dynamodbav:"creationDate,omitempty"`
	Origin       TagOrigin `bson:"origin,omitempty" json:"origin,omitempty" dynamodbav:"origin,omitempty"` // origin informations
}

func (item Tag) GetName() string { return item.Name }

type User struct {
	// ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`           // mongodb default id
	Origin       string `bson:"origin,omitempty" json:"origin,omitempty" dynamodbav:"origin,omitempty"`       // original website
	Name         string `bson:"name,omitempty" json:",omitempty" dynamodbav:"name,omitempty"`                 // userName
	OriginID     string `bson:"originID,omitempty" json:"originID,omitempty" dynamodbav:"originID,omitempty"` // ID from the original website
	CreationDate string `bson:"creationDate,omitempty" json:"creationDate,omitempty" dynamodbav:"creationDate,omitempty"`
}

func (item User) GetOrigin() string   { return item.Origin }
func (item User) GetOriginID() string { return item.OriginID }

type ImageSize struct {
	// ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"` // mongodb default id
	CreationDate string `bson:"creationDate,omitempty" json:"creationDate,omitempty" dynamodbav:"creationDate,omitempty"`
	Box          Box    `bson:"box,omitempty" json:"box,omitempty" dynamodbav:"box,omitempty"` // absolut reference of the top left of new box based on the original sizes
}

type TagOrigin struct {
	Name    string `bson:"name,omitempty" json:"name,omitempty" dynamodbav:"name,omitempty"`          // name of the origin `gui`, `username` or `detector`
	Model   string `bson:"model,omitempty" json:"model,omitempty" dynamodbav:"model,omitempty"`       // name of the model used for the detector
	Weights string `bson:"weights,omitempty" json:"weights,omitempty" dynamodbav:"weights,omitempty"` // weights of the model used for the detector
	// ImageSizeID primitive.ObjectID `bson:"imageSizeID,omitempty" json:"imageSizeID,omitempty" dynamodbav:"origin,omitempty"` // reference to the anchor point
	ImageSizeDate string `bson:"imageSizeDate,omitempty" json:"imageSizeDate,omitempty" dynamodbav:"imageSizeDate,omitempty"` // reference to the anchor point
	Box           Box    `bson:"box,omitempty" json:"box,omitempty" dynamodbav:"box,omitempty"`                               // reference of the bounding box relative to the anchor
}

type Box struct {
	X      *int `bson:"x,omitempty" json:"x,omitempty" dynamodbav:"x,omitempty"`                // top left x coordinate (pointer because 0 is a possible value)
	Y      *int `bson:"y,omitempty" json:"y,omitempty" dynamodbav:"y,omitempty"`                // top left y coordinate (pointer because 0 is a possible value)
	Width  *int `bson:"width,omitempty" json:"width,omitempty" dynamodbav:"width,omitempty"`    // width (pointer because 0 is a possible value)
	Height *int `bson:"height,omitempty" json:"height,omitempty" dynamodbav:"height,omitempty"` // height (pointer because 0 is a possible value)
}
