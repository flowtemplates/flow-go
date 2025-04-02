package analyzer

import (
	"fmt"
	"strings"

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

type TypeErrors []TypeError

func (l TypeErrors) Error() string {
	switch len(l) {
	case 0:
		return "no errors"
	case 1:
		return l[0].Error()
	}

	b := []string{}
	for _, e := range l {
		b = append(b, e.Error())
	}

	return strings.Join(b, ", ")
}

// Err returns an error equivalent to this error list.
// If the list is empty, Err returns nil.
func (l TypeErrors) Err() error {
	if len(l) == 0 {
		return nil
	}

	return l
}

// Add adds an [Error] with given position and error message to an [TypeErrors].
func (l *TypeErrors) Add(err *TypeError) {
	*l = append(*l, *err)
}
