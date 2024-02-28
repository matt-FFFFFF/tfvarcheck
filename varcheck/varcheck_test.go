package varcheck_test

import (
	"testing"

	"github.com/matt-FFFFFF/tfvarcheck/varcheck"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestNewTypeConstraintWithDefaultsFromBytes(t *testing.T) {
	cases := []struct {
		name      string
		got       string
		wantTyd   varcheck.TypeConstraintWithDefaults
		hasErrors bool
	}{
		{
			name: "simple",
			got: `object({
			foo = object({
				bar = optional(string, "baz")
			})
		})`,
			wantTyd: varcheck.TypeConstraintWithDefaults{
				Type:    cty.Object(map[string]cty.Type{"foo": cty.ObjectWithOptionalAttrs(map[string]cty.Type{"bar": cty.String}, []string{"bar"})}),
				Default: nil,
			},
			hasErrors: false,
		},
		{
			name: "invalid",
			got: `obect({
			foo = object({
				bar = optional(string, "baz")
			})
		})`,
			wantTyd: varcheck.TypeConstraintWithDefaults{
				Type:    cty.Type{},
				Default: nil,
			},
			hasErrors: true,
		},
	}

	for _, tc := range cases {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotTyd, diags := varcheck.NewTypeConstraintWithDefaultsFromBytes([]byte(tt.got))
			assert.Equal(t, tt.wantTyd.Type, gotTyd.Type)
			assert.Equal(t, tt.hasErrors, diags.HasErrors())
		})
	}
}

func TestNewTypeConstraintWithDefaultsFromString(t *testing.T) {
	_, diags := varcheck.NewTypeConstraintWithDefaultsFromBytes([]byte(""))
	assert.True(t, diags.HasErrors())
}
