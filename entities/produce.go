package entities

import (
	"gongo/utils"
)

//produce entity definition and structure
type Produce struct {
	ProduceId   interface{} `bson:"_id"`
	ProduceName string
	ProduceCode string
	UnitPrice   float32
}

//loads new ProduceCode
func (p *Produce) LoadNewProduceCode() (error) {
	newCode, err := utils.UUID()
	if err == nil {
		p.ProduceCode = *newCode
	}
	return err
}