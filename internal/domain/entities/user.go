package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserData struct {
	Name           string `json:"name" bson:"name"`
	DocumentNumber string `json:"document_number" bson:"document_number"`
	ZipCode        string `json:"zipCode" bson:"zipCode"`
}

type DocumentValidationResponse struct {
	IsValid bool `json:"isValid" bson:"isValid"`
}

type AddressResponse struct {
	Street  string `json:"street" bson:"street"`
	City    string `json:"city" bson:"city"`
	State   string `json:"state" bson:"state"`
	Zipcode string `json:"zipcode" bson:"zipcode"`
}

type ProcessedUser struct {
	ID        primitive.ObjectID    `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string                `json:"name" bson:"name"`
	Document  DocumentUserProcessed `json:"document" bson:"document"`
	Address   AddressResponse       `json:"address" bson:"address"`
	Status    string                `json:"status" bson:"status"`
	Message   string                `json:"message" bson:"message"`
	CreatedAt time.Time             `json:"created_at" bson:"created_at"`
}

type DocumentUserProcessed struct {
	DocumentNumber string `json:"document_number" bson:"document_number"`
	IsValid        bool   `json:"is_valid" bson:"is_valid"`
}
