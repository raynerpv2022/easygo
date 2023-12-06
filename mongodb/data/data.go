package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Podcast struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name,omitempty"`
	Title      string             `bson:"title,omitempty"`
	CreateTime time.Time          `bso:"createTime,omitempty"`
	UpdateTime time.Time          `bson:"updateTime,omitempty"`
}

// check if data is not empty
func (p *Podcast) IsNameEmpty() bool {
	return p.Name == ""
}
