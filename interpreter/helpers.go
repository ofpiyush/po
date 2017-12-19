package interpreter

import (
	"strings"
)

// TrueMap has Hacky way of classification :P
// This map should get replaced by a classifier
var TrueMap = map[string]bool{
	"true": true,
	"yes":  true,
	"ok":   true,
	"sure": true,
	"haan": true,
	"si":   true,
}

func isTruthy(o *Object) bool {
	switch o.Type {
	case STRING:
		return TrueMap[strings.ToLower(o.Value.(string))]
	case INT:
		return o.Value.(int) > 0
	case BOOL:
		return o.Value.(bool)
	case NilObj.Type:
		return false
	}
	return false
}
