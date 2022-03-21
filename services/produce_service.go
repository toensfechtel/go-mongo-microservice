package services

import (
	"context"
	"gongo/entities"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	PRODUCE_COLLECTION_NAME string = "produce"
	SUPERMARKET_DB_NAME string = "supermarket"
)

//MongoConstruct encapsolates the mongo collection that maps to this service
type ProduceMongoConstruct struct {
	CollectionGrunt *mongo.Collection
	CollectionName  string
}

//constructs the ProduceMongoConstruct
func NewProduceService(client *mongo.Client) *ProduceMongoConstruct {
	//gets the produce collection
	//we use the collection entity to
	//add, delete, find, etc.. - produce
	collection := client.Database(SUPERMARKET_DB_NAME).Collection(PRODUCE_COLLECTION_NAME)
	cosntruct := ProduceMongoConstruct{
		CollectionGrunt: collection,
		CollectionName:  PRODUCE_COLLECTION_NAME,
	}
	return &cosntruct
}

//get a specific produce by ProduceCode
func (mc *ProduceMongoConstruct) GetProduceByProduceCode(code string) (*entities.Produce, error)  {
	var produce entities.Produce
	//serach for Produce
	if err := mc.CollectionGrunt.FindOne(context.TODO(), bson.M{ "producecode": code }).Decode(&produce); err != nil {
		return nil, err
	}
	return &produce, nil
}

func (mc *ProduceMongoConstruct) GetProduce() (*[]entities.Produce, error) {
	
	//we define the filter by which we search for all produce
	filter := bson.M{}
	//we find produce
	cursor, err := mc.CollectionGrunt.Find(context.TODO(), filter)
	if err != nil {
		log.Println("Err: could not find produce.")
		return nil, err
	}
	//decodes each document into produce
	var produce []entities.Produce
	if cerr := cursor.All(context.TODO(), &produce); cerr != nil {
		log.Println("Err: could not map found produce.")
		return nil, cerr
	}
	//if no produce is found, we send back an empty slice
	if produce == nil {
		produce = []entities.Produce{}
	}
	return &produce, nil
}

func (mc *ProduceMongoConstruct) AddProduce(produce *entities.Produce) (interface{}, error) {
	//add a new id to to the Produce that we are about to insert
	produce.ProduceId = primitive.NewObjectID()
	//insert produce
	reponseVal, err := mc.CollectionGrunt.InsertOne(context.TODO(), produce, options.InsertOne())
	if err != nil {
		log.Println("Err: could not insert produce.")
		return nil, err
	}
	
	//we return the inserted Produce id
	return &reponseVal.InsertedID, nil
}


func (mc *ProduceMongoConstruct) DeleteProduce(produceCode string) (*int64, error) {
	//we delete the Produce by ProduceCode
	result, err := mc.CollectionGrunt.DeleteOne(context.TODO(), bson.M{"producecode": produceCode})
	if err != nil {
		log.Println("Err: could not delete produce.")
		return nil, err
	}
	//we return the ammount of Produce that was deleted
	return &result.DeletedCount, nil
}
