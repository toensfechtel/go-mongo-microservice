package routing

import (
	"errors"
	"gongo/entities"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

var (
	ERR_PRODUCE_PAYLOAD interface{} = ""
	ENTITITY_IDENTIFIER string = "code"
	PRODUCE_NOT_FOUND string = "Err: Produce not found."
	PRODUCE_ADD_ERR     string = "Err: could not add all produce."
	PRODUCE_DELETE_ERR string = "Err: could not delete produce."
	PRODUCE_DELETE_INPUT_ERR string = "Err: please provide the proper delete input path DELETE /produce/{produce-code}."
	PRODUCE_GET_ERR string = "Err: could not get produce."
	PRODUCE_UUID_GEN_ERROR string = "Err: a produce code could not be generated."
	PRODUCE_ADD_INPUT_ERR string = "Err: please provide proper Produce payload { Produce: [] }."
	PRODUCE_ADD_BAD_PARAM string = "Err: please provide a ProduceName for all produce to be added."
)

//get produce result contract
type GetProduceResult struct {
	Produce interface{}
	Err string
}
//get produce route handler
func (m RouteContext) Get_ProduceRoute(c *gin.Context) {
	//get path paramter => this corresponds to code => /produce/:code
	produceCode := c.Param(ENTITITY_IDENTIFIER)
	
	var produce interface{}
	var err error

	//no produce code was found, then get all produce
	if produceCode == "" {
		//get all produce from data source
		produce, err = m.ProduceService.GetProduce()
	}
	
	//a produce code was found, then we send back the specific produce
	//if it is found
	if  produceCode != "" {
		p, e := m.ProduceService.GetProduceByProduceCode(produceCode)
		err = e
		//handling the non transitivity of nil pointer to interface
		if p == nil {
			produce = nil
		}
	} 

	if err != nil && produce == nil {
		//data source has thrown an error
		//the only method call that returns produce as nil is GetProduceByProduceCode
		AbortNow(c, http.StatusNotFound, GetProduceResult{ Produce: entities.Produce{}, Err: PRODUCE_NOT_FOUND}, err)
		return 
	}

	if err != nil {
		//data source has thrown an error
		//a code 500 is thrown with a predefined constant error msg
		AbortNow(c, 404, GetProduceResult{ Produce: &ERR_PRODUCE_PAYLOAD, Err: PRODUCE_GET_ERR}, err)
		return 
	}
	//define response object
	result := GetProduceResult{
		Produce: produce,
		Err: "",
	}
	c.IndentedJSON(http.StatusOK, result)
}

//delete produce result contract
type DeleteProduceResult struct {
	DeletedCount int64
	Err string
}
//delete produce route handler
func (m RouteContext) Delete_ProduceRoute(c *gin.Context) {
	//get path paramter => this corresponds to code => /produce/:code
	produceCode := c.Param(ENTITITY_IDENTIFIER)
	//if a produce code is not provided then a bad request is sent back
	//with a preset constant predefined error msg
	if produceCode == "" {
		AbortNow(c, http.StatusBadRequest, DeleteProduceResult{ DeletedCount: 0, Err: PRODUCE_DELETE_INPUT_ERR}, errors.New(PRODUCE_DELETE_INPUT_ERR))
		return 
	}
	//delete produce by produceCode from data source
	deleted, err := m.ProduceService.DeleteProduce(produceCode)
	if err != nil {
		//data source has thrown an error
		//a code 500 is thrown with a constant err msg
		AbortNow(c, 500, DeleteProduceResult{ DeletedCount: 0, Err: PRODUCE_DELETE_ERR}, err)
		return 
	}
	//define response object
	response := DeleteProduceResult{
		DeletedCount: *deleted,
		Err: "",
	}
	c.IndentedJSON(http.StatusOK, response)
}

//CreateproduceInput defines the json body that
//we are expecting from caller
type CreateproduceInput struct {
	Produce *[]entities.Produce
}
//create produce result contract
type CreateProduceResult struct {
	ProduceAdded *[]interface{}
	Err          string
}
//add produce route handler
func (m RouteContext) Add_ProduceRoute(c *gin.Context) {
	//a list IDs => for all x in produceAdded, then x is an id of a produce
	//that was succesfully added to the datasource
	produceAdded := []interface{}{}
	
	//map json body to the input contract
	var produce *CreateproduceInput
	if err := c.BindJSON(&produce); err != nil || produce.Produce == nil {
		//a mapping error accurred
		//=> a bad request is assumed and a 
			//constant predefined msg is sent back
		AbortNow(c, http.StatusBadRequest, CreateProduceResult{ ProduceAdded: &produceAdded, Err: PRODUCE_ADD_INPUT_ERR}, errors.New(PRODUCE_DELETE_INPUT_ERR))
		return 
	}

	//adds all produce sent by caller
	for _, p := range *produce.Produce {
		//produce does not have a ProduceName
		//=> a bad request is assumed and a 
			//constant predefined msg is sent back
		if p.ProduceName == "" {
			AbortNow(c, http.StatusBadRequest, CreateProduceResult{ ProduceAdded: &produceAdded, Err: PRODUCE_ADD_BAD_PARAM}, errors.New(PRODUCE_ADD_BAD_PARAM))
			return 
		}
		//Load new produce code for new produce
		//If one was provided it will be over written
		//this way we uphold the integrety of the ProduceCode expected values
		err := p.LoadNewProduceCode()
		if err != nil {
			//ProduceCode UUID was not generated  or 
			//a 500 error code is sent to caller with a constant predefined error msg
			AbortNow(c, 500, CreateProduceResult{ProduceAdded: &produceAdded, Err: PRODUCE_UUID_GEN_ERROR }, err)
		}
		//add produce to data source
		pid, err := m.ProduceService.AddProduce(&p)
		if err != nil {
			//data source has thrown an error
			//a code 500 is thrown with a constant err msg
			AbortNow(c, 500, CreateProduceResult{ProduceAdded: &produceAdded, Err: PRODUCE_ADD_ERR }, err)
			return
		}
		//append the added produce's id to the ledger array
		produceAdded = append(produceAdded, pid)
	}
	//define response object
	results := CreateProduceResult{
		ProduceAdded: &produceAdded,
		Err:          "",
	}
	c.IndentedJSON(http.StatusOK, results)
}

//helper function for aborting the route process
func AbortNow(c *gin.Context, code int, response interface{}, err error) {
	log.Println(err.Error())
	// specified response
	c.JSON(500, response)
	// abort the request
	c.Abort()
}
