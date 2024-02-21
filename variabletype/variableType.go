package variabletype

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/zclconf/go-cty/cty"
)

// VariableType represents a type constraint and default value for a Terraform variable.
type VariableType struct {
	Type    cty.Type
	Default *typeexpr.Defaults
}

// NewVariableTypeFromExpression creates a new VariableType from a given hcl Expression.
func NewVariableTypeFromExpression(exp hcl.Expression) (VariableType, hcl.Diagnostics) {
	t, d, diags := typeexpr.TypeConstraintWithDefaults(exp)
	if diags.HasErrors() {
		return VariableType{}, diags
	}
	return VariableType{
		Type:    t,
		Default: d,
	}, nil
}
