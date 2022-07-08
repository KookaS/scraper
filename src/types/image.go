package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// elements here prevent cyclical imports in other folders

// BodyUpdateImage is the body for the update of an image.
type BodyUpdateImageTagsPush struct {
	Origin string             `bson:"origin,omitempty" json:"origin,omitempty"`
	ID     primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"`
	Tags   []Tag              `bson:"tags,omitempty" json:"tags,omitempty"`
}

type BodyUpdateImageTagsPull struct {
	Origin string             `bson:"origin,omitempty" json:"origin,omitempty"`
	ID     primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"`
	Names  []string           `bson:"names,omitempty" json:"names,omitempty"`
}

type BodyUpdateImageFile struct {
	Origin string `bson:"origin,omitempty" json:"origin,omitempty"`
	Name   string `bson:"name,omitempty" json:"name,omitempty"`
	Box    Box    `bson:"box,omitempty" json:"box,omitempty"`
	File   []byte `bson:"file,omitempty" json:"file,omitempty"`
}
