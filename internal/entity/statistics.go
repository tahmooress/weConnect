package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Statistics struct {
	ID              primitive.ObjectID `bson:"_id"`
	SeriesReference string             `bson:"series_reference"`
	Period          string             `bson:"period"`
	DataValue       float64            `bson:"data_value"`
	Suppressed      string             `bson:"suppressed"`
	Status          string             `bson:"status"`
	Units           string             `bson:"Units"`
	Magnitude       int                `bson:"magnitude"`
	Subject         string             `bson:"subject"`
	Group           string             `bson:"group"`
	SeriesTitle1    string             `bson:"Series_title_1"`
	SeriesTitle2    string             `bson:"Series_title_2"`
	SeriesTitle3    string             `bson:"Series_title_3"`
	SeriesTitle4    string             `bson:"Series_title_4"`
	SeriesTitle5    string             `bson:"Series_title_5"`
	CreatedAt       time.Time          `bson:"created_at"`
}
