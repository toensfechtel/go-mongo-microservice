package services

import (
	"gongo/entities"
)


//IProduceService is the contract that must be upheld by any ProduceService
//This interface is also leveraged to mock the ProduceService entity
type IProduceService interface {
	GetProduceByProduceCode(code string) (*entities.Produce, error)
	GetProduce() (*[]entities.Produce, error)
	AddProduce(produce *entities.Produce) (interface{}, error)
	DeleteProduce(produceCode string) (*int64, error)
}