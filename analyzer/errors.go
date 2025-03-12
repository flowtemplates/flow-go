package analyzer

import (
	"fmt"

	"github.com/flowtemplates/flow-go/types"
)

type TypeError struct {
	ExpectedType types.Type
	Name         string
	Val          string
}

//	func (e *TypeError) String() string {
//		return fmt.Sprintf("TypeError: Variable '%s' expected type '%s'", e.Name, e.ExpectedType)
//	}

func (e TypeError) Error() string {
	return fmt.Sprintf("TypeError: Variable '%s' expected type '%s'", e.Name, e.ExpectedType)
}
