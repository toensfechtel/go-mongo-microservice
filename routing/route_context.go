package routing

import (
	"gongo/services"
)

//we use this struct to define the route context
//this is used to inject services to each of our
//handler methods so that they can act against 
//a data source
//=>> note that the service encapsulates the db client
	//through the use of another struct 
type RouteContext struct {
	ProduceService   services.IProduceService
}