package main

import (
	"context"
	"gongo/routing"
	"gongo/services"
	"log"
	"os"
	"time"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gongo/datadump"
)

//this is a user defined method to close resources.
//this method closes mongoDB connection and cancel context.
func close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {
	//cancelFunc to cancel to context
	defer cancel()
	//client provides a method to close
	//a mongoDB connection.
	defer func() {
		//client.Disconnect method also has deadline.
		//returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func main() {

	//***********************
	//ctx will be used to set deadline for process, here
	//deadline will of 30 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//mongo.Connect return mongo.Client method
	//the URI corresponds to the mongo host => the mongo container spun up
	//if you wish to set a URI to a non local instance, then I recommend
	//making the connection string somesort of secrete, never have your
	//database connection strings inline
	//That being said, this value is usually made into an env variable otherwise,
	//as this is just an example, I added it inline
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017"))
	//release resource when the main
	//function is returned.
	if err != nil {
		panic(err.Error())
	}
	//defers the closing of the client connection upon the return of main
	defer close(client, ctx, cancel)

	//load default data
	datadump.InitiateDataDump(client).CleanDatabase().LoadInitData()

	//inject ProduceService to the route handlers
	routeContext := routing.RouteContext{
		ProduceService: services.NewProduceService(client),
	}

	//load gin
	router := gin.Default()
	//hooks on routes mapped to handlers
	router.GET("/produce", routeContext.Get_ProduceRoute)
	router.GET("/produce/:code", routeContext.Get_ProduceRoute)
	router.DELETE("/produce/:code", routeContext.Delete_ProduceRoute)
	router.POST("/produce", routeContext.Add_ProduceRoute)
	//handle error response when a route is not defined
    router.NoRoute(func(c *gin.Context) {
        c.JSON(404, gin.H{"message": "Not found"})
    })
	ginServerAddress := getServerAddress()
	log.Println(ginServerAddress)
	router.Run(ginServerAddress)
}

func getServerAddress() string {
	log.Println("updated instance")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "8080"
	}
	addr := host + ":" + port
	return addr
}
