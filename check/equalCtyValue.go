package check

import "github.com/zclconf/go-cty/cty"

func EqualCtyValue(got, want cty.Value) bool {
	ret := got.Equals(want)
	return ret.True()
}
