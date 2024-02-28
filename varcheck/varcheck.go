// Package varcheck is the entry point of the application.
package varcheck

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

// VarCheck is the struct that represents a variable check.
type VarCheck struct {
	Default                cty.Value                  // Default value for the interface as a cty.Value
	Nullable               bool                       // Whether the variable should be nullable.
	TypeConstraintWithDefs TypeConstraintWithDefaults // Strong type representing the type as well as default values. Use NewTypeConstraintWithDefaultsFromExp or NewTypeConstraintWithDefaultsFromBytes to create.
}

// TypeConstraintWithDefaults represents a type constraint and default value for a Terraform variable.
type TypeConstraintWithDefaults struct {
	Type    cty.Type           // The type constraint, this will always have a value
	Default *typeexpr.Defaults // The default value, this will be nil if no defaults are used in the type constraint
}

// NewVarCheck creates a new VarCheck struct.
func NewVarCheck(ty TypeConstraintWithDefaults, def cty.Value, nullable bool) VarCheck {
	return VarCheck{
		Default:                def,
		Nullable:               nullable,
		TypeConstraintWithDefs: ty,
	}
}

// NewVariableTypeFromExpression creates a new TypeConstraintWithDefaults from a given hcl.Expression.
func NewTypeConstraintWithDefaultsFromExp(exp hcl.Expression) (TypeConstraintWithDefaults, hcl.Diagnostics) {
	t, d, diags := typeexpr.TypeConstraintWithDefaults(exp)
	if diags.HasErrors() {
		return TypeConstraintWithDefaults{}, diags
	}
	return TypeConstraintWithDefaults{
		Type:    t,
		Default: d,
	}, nil
}

// NewTypeConstraintWithDefaultsFromString creates a new TypeConstraintWithDefaults from a byte slice.
func NewTypeConstraintWithDefaultsFromBytes(b []byte) (TypeConstraintWithDefaults, hcl.Diagnostics) {
	exp, diags := hclsyntax.ParseExpression(b, "variables.tf", hcl.Pos{})
	if diags.HasErrors() {
		return TypeConstraintWithDefaults{}, diags
	}
	return NewTypeConstraintWithDefaultsFromExp(exp)
}
