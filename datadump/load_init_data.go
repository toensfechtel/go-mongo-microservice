package datadump

import (
	"context"
	"fmt"
	"gongo/entities"
	"gongo/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	DATA_DUMP []entities.Produce = []entities.Produce{
		{
			ProduceName: "Lettuce",
			ProduceCode: "A12T-4GH7-QPL9-3N4M",
			UnitPrice:   3.46,
		},
		{
			ProduceName: "Peach",
			ProduceCode: "E5T6-9UI3-TH15-QR88",
			UnitPrice:   2.99,
		},
		{
			ProduceName: "Green Pepper",
			ProduceCode: "YRT6-72AS-K736-L4AR",
			UnitPrice:   0.79,
		},
		{
			ProduceName: "Gala Apple",
			ProduceCode: "TQ4C-VV6T-75ZX-1RMR",
			UnitPrice:   3.59,
		},
	}
)

type ClientWrapper struct {
	m *mongo.Client
}

func InitiateDataDump(client *mongo.Client) (*ClientWrapper) {
	wrapper := ClientWrapper{
		m: client,
	}
	return &wrapper
}

func (cw *ClientWrapper) CleanDatabase() (*ClientWrapper) {
	//ping to test mongodb database connection
	err := cw.m.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	//drop database if it exists
	err = cw.m.Database(services.SUPERMARKET_DB_NAME).Drop(context.TODO())
	if err != nil {
		panic(err)
	}
	return cw
}

func (cw *ClientWrapper) LoadInitData() (*ClientWrapper) {
	//adds produce IDs and maps to interface
	payload := []interface{}{}
	for _, p := range DATA_DUMP {
		objID := primitive.NewObjectID()
		p.ProduceId = objID
		payload = append(payload, p)
	}
	//batch insert the produce
	insertManyResult, err := cw.m.Database(services.SUPERMARKET_DB_NAME).Collection(services.PRODUCE_COLLECTION_NAME).InsertMany(context.TODO(), payload)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
	return cw
}


