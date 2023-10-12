package render

import (
	cmath "github.com/alacrity-engine/core/math"
	"github.com/alacrity-engine/core/system/collections"
)

// TODO: use fixed point numbers.

type Geometric interface {
	collections.Comparable
}

type Range struct {
	Z1 cmath.Fixed
	Z2 cmath.Fixed
}

func (rg Range) Less(other collections.Comparable) bool {
	if point, ok := other.(Point); ok {
		return point.Z.Less(rg.Z1)
	}

	otherLS := other.(Range)
	return rg.Z2.Less(otherLS.Z1)
}

func (rg Range) Greater(other collections.Comparable) bool {
	if point, ok := other.(Point); ok {
		return point.Z.Greater(rg.Z2)
	}

	otherLS := other.(Range)
	return rg.Z1.Greater(otherLS.Z2)
}

func (rg Range) Equal(other collections.Comparable) bool {
	if point, ok := other.(Point); ok {
		return point.Z.GreaterOrEqual(rg.Z1) && point.Z.LessOrEqual(rg.Z2)
	}

	otherLS := other.(Range)
	return rg.Z1 == otherLS.Z1 && rg.Z2 == otherLS.Z2
}

type Point struct {
	Z cmath.Fixed
}

func (p Point) Less(other collections.Comparable) bool {
	if r, ok := other.(Range); ok {
		return p.Z.Less(r.Z1)
	}

	otherLS := other.(Point)
	return p.Z.Less(otherLS.Z)
}

func (p Point) Greater(other collections.Comparable) bool {
	if r, ok := other.(Range); ok {
		return p.Z.Greater(r.Z2)
	}

	otherLS := other.(Point)
	return p.Z.Greater(otherLS.Z)
}

func (p Point) Equal(other collections.Comparable) bool {
	if r, ok := other.(Range); ok {
		return p.Z.GreaterOrEqual(r.Z1) && p.Z.LessOrEqual(r.Z2)
	}

	otherLS := other.(Point)
	return p.Z == otherLS.Z
}
