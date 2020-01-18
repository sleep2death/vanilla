package core

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Player Data struct
type Player struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Created  int64              `bson:"created"`

	Crystal int64

	PrivateTiles []Tile
	GlobalTiles  []Tile
	Heroes       []Hero

	Gold  int64 `bson:"gold" json:"gold"`
	Food  int64 `bson:"food" json:"food"`
	Wood  int64 `bson:"wood" json:"wood"`
	Stone int64 `bson:"stone" json:"stone"`
	Iron  int64 `bson:"iron" json:"iron"`
}

// Tile Data struct
type Tile struct {
	ID        string     `bson:"tid" json:"tid"`
	Name      string     `bson:"name" json:"name"`
	Type      int8       `bson:"type" json:"type"`
	Terrain   int8       `bson:"terrain" json:"terrain"`
	Buildings []Building `bson:"buildings" json:"buildings"`
}

// Hero Data struct
type Hero struct {
	ID   string
	Name string
}

// Building Data struct
type Building struct {
	ID   string
	Name string
}
