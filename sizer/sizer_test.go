package sizers

import (
	"github.com/peterhoward42/umli/parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSizerSetsItsScalarAttributesCorrectly(t *testing.T) {
	assert := assert.New(t)
	statements := parser.MustCompileParse("lane A foo")
	sizer := NewSizer(2000, 20.0, statements)

	assert.InDelta(20, sizer.DiagramPadT, 0.1)
	assert.InDelta(10, sizer.InteractionLinePadB, 0.1)
	assert.InDelta(10, sizer.InteractionLineTextPadB, 0.1)
	assert.InDelta(30, sizer.ArrowLen, 0.1)
	assert.InDelta(12, sizer.ArrowHeight, 0.1)
	assert.InDelta(10, sizer.DashLineDashLen, 0.1)
	assert.InDelta(5, sizer.DashLineDashGap, 0.1)
}

func TestNewSizerComposesItsDelegatesProperly(t *testing.T) {
	assert := assert.New(t)
	statements := parser.MustCompileParse("lane A foo")
	sizer := NewSizer(2000, 20.0, statements)
	// Have the embedded individual lane sizing data structures been
	// installed?
	assert.Equal(1, len(sizer.Lanes.Individual), "Should be one Lane")
}