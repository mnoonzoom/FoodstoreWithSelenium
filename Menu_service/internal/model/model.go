package model

type MenuItem struct {
	ID          string  `bson:"_id,omitempty" json:"id"`
	Name        string  `bson:"name" json:"name"`
	Description string  `bson:"description" json:"description"`
	Price       float64 `bson:"price" json:"price"`
	Available   bool    `bson:"available" json:"available"`
	Category    string  `bson:"category" json:"category"`
	ImageURL    string  `bson:"image_url" json:"image_url"`
}
