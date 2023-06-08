package errsupply

import (
	"errors"

	"github.com/Mines-Little-Theatre/team-trivia-scraper/bot/suppliers"
)

type ErrSupply struct{}

func init() {
	suppliers.RegisterSupplier("errsupply", ErrSupply{})
}

func (ErrSupply) SupplyData(*suppliers.SupplierContext) error {
	return errors.New("jesus, you're making me sound like air supply")
}
