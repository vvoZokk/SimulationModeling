package transaction

import (
	"fmt"
)

type Transaction struct {
	Id int
}

func (transaction *Transaction) Test(id int) {
	transaction.Id = id
}

func (transaction Transaction) String() string {
	return fmt.Sprintf("transaction id is %d", transaction.Id)
}
