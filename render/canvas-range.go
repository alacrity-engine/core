package render

import "github.com/alacrity-engine/core/system/collections"

type Geometric interface {
	collections.Comparable
}

type Range struct {
	Z1 float32
	Z2 float32
}

func (rg Range) Less(other collections.Comparable) bool {
	if point, ok := other.(Point); ok {
		return point.Z < rg.Z1
	}

	otherLS := other.(Range)
	return rg.Z2 < otherLS.Z1
}

func (rg Range) Greater(other collections.Comparable) bool {
	if point, ok := other.(Point); ok {
		return point.Z > rg.Z2
	}

	otherLS := other.(Range)
	return rg.Z1 > otherLS.Z2
}

func (rg Range) Equal(other collections.Comparable) bool {
	if point, ok := other.(Point); ok {
		return point.Z >= rg.Z1 && point.Z <= rg.Z2
	}

	otherLS := other.(Range)
	return rg.Z1 == otherLS.Z1 && rg.Z2 == otherLS.Z2
}

type Point struct {
	Z float32
}

func (p Point) Less(other collections.Comparable) bool {
	if r, ok := other.(Range); ok {
		return p.Z < r.Z1
	}

	otherLS := other.(Point)
	return p.Z < otherLS.Z
}

func (p Point) Greater(other collections.Comparable) bool {
	if r, ok := other.(Range); ok {
		return p.Z > r.Z2
	}

	otherLS := other.(Point)
	return p.Z > otherLS.Z
}

func (p Point) Equal(other collections.Comparable) bool {
	if r, ok := other.(Range); ok {
		return p.Z >= r.Z1 && p.Z <= r.Z2
	}

	otherLS := other.(Point)
	return p.Z == otherLS.Z
}
