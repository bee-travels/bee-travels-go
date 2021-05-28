package database

import (
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type Datatype string

const (
	Float     Datatype = "float"
	Int       Datatype = "int"
	String    Datatype = "string"
	Timestamp Datatype = "timestamp"
	Uint      Datatype = "uint"
)

func ValueByDataType(value, datatype string) (interface{}, error) {
	switch Datatype(datatype) {
	case Float:
		return strconv.ParseFloat(value, 64)
	case Int:
		return strconv.ParseInt(value, 10, 64)
	case String:
		return value, nil
	case Timestamp:
		t, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return time.Unix(t, 0), nil
	case Uint:
		return strconv.ParseUint(value, 10, 64)
	}
	return nil, errors.Errorf("Unknown datatype: %s", datatype)
}

func UuidArray(array []uuid.UUID) []string {
	a := make([]string, len(array))
	for i := 0; i < len(array); i++ {
		a[i] = array[i].String()
	}
	return a
}
