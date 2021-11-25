package office

import (
	"errors"
	"fmt"
)

type Employee interface {
	GetCurrentLocation() Location
	MoveToLocation(Location) error
}

type baseEmployeeParams struct {
	location Location
	accesses []string
}

func (e baseEmployeeParams) GetCurrentLocation() Location {
	return e.location
}

func (e *baseEmployeeParams) MoveToLocation(loc Location) error {
	sl := e.accesses
	for _, v := range sl {
		if loc.GetLocationTitle() == v{
			e.location = loc
			return nil
		}
	}
	return fmt.Errorf("fail employee move: %w",errors.New("MoveToLocation Error or access denied"))
}

var ErrUnknownEmplType = errors.New("unknown employee type")

func NewEmployeeFactory(title string) (Employee, error) {
	switch title {
	case "hr":
		return newHr(), nil
	case "itSecurity":
		return newItSecurity(), nil
	}
	return nil, fmt.Errorf("newEmployee: %w", ErrUnknownEmplType)
}

// TODO:: must impl Employee
type hr struct {
	baseEmployeeParams
}


func newHr() Employee {
	accesses := []string{"office", "workArea"}
	return &hr{
		baseEmployeeParams{
			accesses: accesses,
		},
	}
}

// TODO:: must impl Employee
type itSecurity struct {
	baseEmployeeParams
}

func newItSecurity() Employee {
	accesses := []string{"office", "workArea", "servers"}	
	return &itSecurity{
		baseEmployeeParams{
			accesses: accesses,
		},
	}
}