package math

import (
	"github.com/alacrity-engine/core/system/collections"
	"github.com/robaho/fixed"
)

type Fixed fixed.Fixed

func (fixedObj Fixed) Less(other collections.Comparable) bool {
	otherF := other.(Fixed)
	return fixed.Fixed(fixedObj).LessThan(fixed.Fixed(otherF))
}

func (fixedObj Fixed) LessOrEqual(other collections.Comparable) bool {
	otherF := other.(Fixed)
	return fixed.Fixed(fixedObj).LessThanOrEqual(fixed.Fixed(otherF))
}

func (fixedObj Fixed) Greater(other collections.Comparable) bool {
	otherF := other.(Fixed)
	return fixed.Fixed(fixedObj).GreaterThan(fixed.Fixed(otherF))
}

func (fixedObj Fixed) GreaterOrEqual(other collections.Comparable) bool {
	otherF := other.(Fixed)
	return fixed.Fixed(fixedObj).GreaterThanOrEqual(fixed.Fixed(otherF))
}

func (fixedObj Fixed) Equal(other collections.Comparable) bool {
	otherF := other.(Fixed)
	return fixed.Fixed(fixedObj).Equal(fixed.Fixed(otherF))
}

func FixedFromFloat32(num float32) Fixed {
	return Fixed(fixed.NewF(float64(num)))
}

func FixedFromFloat64(num float64) Fixed {
	return Fixed(fixed.NewF(num))
}
