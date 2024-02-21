// Package main is the entry point of the application.
package main

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

type VarCheck struct {
	Default  cty.Value       // Default value for the interface as a cty.Value
	Enabled  bool            // Whether to test this interface interface.
	Link     string          // Link to the interface documentation.
	Name     string          // Name of the interface.
	Nullable bool            // Whether the interface is nullable.
	Type     string          // String representing the type value.
	Severity tflint.Severity // Severity of the rule.
	//TypeStrong
}

// TypeExpression returns an HCL expression that represents the interface type.
// If the interface cannot be correctly parsed, this function will panic.
func (v VarCheck) TypeExpression() hcl.Expression {
	e, d := hclsyntax.ParseExpression([]byte(v.Type), "variables.tf", hcl.Pos{})
	if d.HasErrors() {
		panic(d.Error())
	}
	return e
}

// TerraformVar returns a string that represents the interface as the
// minimum required Terraform variable definition for testing.
func (v VarCheck) TerraformVar() (string, error) {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	varBlock := rootBody.AppendNewBlock("variable", []string{v.Name})
	varBody := varBlock.Body()
	// check the Type constraint is valid and panic if not
	if _, _, diags := typeexpr.TypeConstraintWithDefaults(v.TypeExpression()); diags.HasErrors() {
		return "", diags
	}
	// I couldn't get the hclwrite to work with the type constraint so I'm just adding it as a string
	// using SetSAttributeRaw and hclWrite.Token.
	varBody.SetAttributeRaw("type", hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenStringLit,
			Bytes: []byte(" " + v.Type),
		},
	})
	varBody.SetAttributeValue("default", v.Default)
	// If the interface is not nullable, set the nullable attribute to false.
	// the default is true so we only need to set it if it's false.
	if !v.Nullable {
		varBody.SetAttributeValue("nullable", cty.False)
	}
	return string(f.Bytes()), nil
}
