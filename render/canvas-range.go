package render

import "github.com/alacrity-engine/core/system/collections"

type Range struct {
	A float32
	B float32
}

func (rg Range) Less(other collections.Comparable) bool {
	if point, ok := other.(Point); ok {
		return point.Num < rg.A
	}

	otherLS := other.(Range)
	return rg.B < otherLS.A
}

func (rg Range) Greater(other collections.Comparable) bool {
	if point, ok := other.(Point); ok {
		return point.Num > rg.B
	}

	otherLS := other.(Range)
	return rg.A > otherLS.B
}

func (rg Range) Equal(other collections.Comparable) bool {
	if point, ok := other.(Point); ok {
		return point.Num >= rg.A && point.Num <= rg.B
	}

	otherLS := other.(Range)
	return rg.A == otherLS.A && rg.B == otherLS.B
}

type Point struct {
	Num float32
}

func (p Point) Less(other collections.Comparable) bool {
	if r, ok := other.(Range); ok {
		return p.Num < r.A
	}

	otherLS := other.(Point)
	return p.Num < otherLS.Num
}

func (p Point) Greater(other collections.Comparable) bool {
	if r, ok := other.(Range); ok {
		return p.Num > r.B
	}

	otherLS := other.(Point)
	return p.Num > otherLS.Num
}

func (p Point) Equal(other collections.Comparable) bool {
	if r, ok := other.(Range); ok {
		return p.Num >= r.A && p.Num <= r.B
	}

	otherLS := other.(Point)
	return p.Num == otherLS.Num
}
