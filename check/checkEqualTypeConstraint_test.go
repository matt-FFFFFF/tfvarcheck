package check_test

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/matt-FFFFFF/tfvarcheck/check"
	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
)

func hclExpressionFromString(expr string) hcl.Expression {
	e, diags := hclsyntax.ParseExpression([]byte(expr), "test.tf", hcl.Pos{})
	if diags.HasErrors() {
		panic(diags)
	}
	return e
}

func TestCheckEqualTypeConstraints(t *testing.T) {
	cases := []struct {
		Name   string
		Want   hcl.Expression
		Got    hcl.Expression
		Result bool
	}{
		{
			Name:   "Same",
			Want:   hclExpressionFromString("object({kind = string, name = optional(string, null)})"),
			Got:    hclExpressionFromString("object({kind = string, name = optional(string, null)})"),
			Result: true,
		},
		{
			Name:   "Different",
			Want:   hclExpressionFromString("object({kind = string, name = optional(string, null)})"),
			Got:    hclExpressionFromString("object({kind = string, name = optional(number, null)})"),
			Result: false,
		},
		{
			Name:   "IncorrectDefaults",
			Want:   hclExpressionFromString("object({kind = string, name = optional(number, null)})"),
			Got:    hclExpressionFromString("object({kind = string, name = optional(number, 2)})"),
			Result: false,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			gotType, diags := varcheck.NewTypeConstraintWithDefaultsFromExp(tc.Got)
			if diags.HasErrors() {
				panic(diags)
			}
			wantType, diags := varcheck.NewTypeConstraintWithDefaultsFromExp(tc.Want)
			if diags.HasErrors() {
				panic(diags)
			}
			res := check.CheckEqualTypeConstraints(gotType, wantType)
			if res != tc.Result {
				t.Errorf("Test %s: Expected %v, got %v", tc.Name, tc.Result, res)
			}
		})
	}
}
