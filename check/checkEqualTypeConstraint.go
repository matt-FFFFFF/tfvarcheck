package check

import (
	"reflect"

	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
)

// EqualTypeConstraints checks if two supplied hcl Expressions are in fact type constraints,
// and if they are that they are equal.
func EqualTypeConstraints(type1, type2 varcheck.TypeConstraintWithDefaults) bool {
	return reflect.DeepEqual(type1, type2)
}
