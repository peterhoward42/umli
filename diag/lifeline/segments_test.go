package lifeline

import (
	"testing"

	"github.com/peterhoward42/umli/diag/nogozone"
	"github.com/peterhoward42/umli/dsl"
	"github.com/peterhoward42/umli/geom"
	"github.com/stretchr/testify/assert"
)

func TestWithJustOneNoGoZoneGap(t *testing.T) {
	assert := assert.New(t)

	// We will check the segments generated for lifeline b.

	// Create a noGoZone that affects b
	a := &dsl.Statement{}
	b := &dsl.Statement{}
	c := &dsl.Statement{}
	allLifelines := []*dsl.Statement{a, b, c}
	noGoZones := []nogozone.NoGoZone{
		nogozone.NewNoGoZone(geom.NewSegment(10, 50), a, c),
	}

	topOfLifeline := 1.0
	bottomOfLifeline := 100.0

	segments := LifelineSegments{}
	minSegLen := 0.1
	segments.Assemble(b, topOfLifeline, bottomOfLifeline, minSegLen,
		noGoZones, *NewBoxTracker(), allLifelines)

	assert.Len(segments.Segs, 2)
	assert.Equal(segments.Segs[0], geom.NewSegment(1.0, 10))
	assert.Equal(segments.Segs[1], geom.NewSegment(50.0, 100))
}

func TestWithJustOneBoxGap(t *testing.T) {
	assert := assert.New(t)

	boxes := *NewBoxTracker()
	err := boxes.AddStartingAt(70)
	assert.NoError(err)
	err = boxes.TerminateAt(90)
	assert.NoError(err)

	topOfLifeline := 1.0
	bottomOfLifeline := 100.0

	segments := LifelineSegments{}
	minSegLen := 0.1
	var lifeline *dsl.Statement = nil
	segments.Assemble(lifeline, topOfLifeline, bottomOfLifeline, minSegLen,
		[]nogozone.NoGoZone{}, boxes, []*dsl.Statement{})

	assert.Len(segments.Segs, 2)
	assert.Equal(segments.Segs[0], geom.NewSegment(1.0, 70))
	assert.Equal(segments.Segs[1], geom.NewSegment(90.0, 100))
}

func TestConsumesBothNoGoZonesAndActivityBoxesWhenBothPresent(t *testing.T) {
	assert := assert.New(t)

	a := &dsl.Statement{}
	b := &dsl.Statement{}
	c := &dsl.Statement{}
	allLifelines := []*dsl.Statement{a, b, c}
	noGoZones := []nogozone.NoGoZone{
		nogozone.NewNoGoZone(geom.NewSegment(10, 50), a, c),
	}

	boxes := NewBoxTracker()
	err := boxes.AddStartingAt(70)
	assert.NoError(err)
	err = boxes.TerminateAt(90)
	assert.NoError(err)

	topOfLifeline := 1.0
	bottomOfLifeline := 100.0

	segments := LifelineSegments{}
	minSegLen := 0.1
	segments.Assemble(b, topOfLifeline, bottomOfLifeline, minSegLen,
		noGoZones, *boxes, allLifelines)

	assert.Len(segments.Segs, 3)
	assert.Equal(segments.Segs[0], geom.NewSegment(1.0, 10))
	assert.Equal(segments.Segs[1], geom.NewSegment(50.0, 70))
	assert.Equal(segments.Segs[2], geom.NewSegment(90.0, 100))
}

func TestItSortsTheGaps(t *testing.T) {
	assert := assert.New(t)

	a := &dsl.Statement{}
	b := &dsl.Statement{}
	c := &dsl.Statement{}
	allLifelines := []*dsl.Statement{a, b, c}
	noGoZones := []nogozone.NoGoZone{
		nogozone.NewNoGoZone(geom.NewSegment(50, 60), a, c),
		nogozone.NewNoGoZone(geom.NewSegment(10, 20), a, c),
	}

	boxes := *NewBoxTracker()

	topOfLifeline := 1.0
	bottomOfLifeline := 100.0

	segments := LifelineSegments{}
	minSegLen := 0.1
	segments.Assemble(b, topOfLifeline, bottomOfLifeline, minSegLen,
		noGoZones, boxes, allLifelines)

	assert.Len(segments.Segs, 3)
	assert.Equal(segments.Segs[0], geom.NewSegment(1.0, 10))
	assert.Equal(segments.Segs[1], geom.NewSegment(20.0, 50))
	assert.Equal(segments.Segs[2], geom.NewSegment(60.0, 100))
}

func TestItMergesTheGaps(t *testing.T) {
	assert := assert.New(t)

	a := &dsl.Statement{}
	b := &dsl.Statement{}
	c := &dsl.Statement{}
	allLifelines := []*dsl.Statement{a, b, c}
	noGoZones := []nogozone.NoGoZone{
		nogozone.NewNoGoZone(geom.NewSegment(50, 60), a, c),
		nogozone.NewNoGoZone(geom.NewSegment(50, 70), a, c),
	}

	boxes := *NewBoxTracker()

	topOfLifeline := 1.0
	bottomOfLifeline := 100.0

	segments := LifelineSegments{}
	minSegLen := 0.1
	segments.Assemble(b, topOfLifeline, bottomOfLifeline, minSegLen,
		noGoZones, boxes, allLifelines)

	assert.Len(segments.Segs, 2)
	assert.Equal(segments.Segs[0], geom.NewSegment(1.0, 50))
	assert.Equal(segments.Segs[1], geom.NewSegment(70.0, 100))
}

func TestItDiscardsTinySegmentsGaps(t *testing.T) {
	assert := assert.New(t)

	a := &dsl.Statement{}
	b := &dsl.Statement{}
	c := &dsl.Statement{}
	allLifelines := []*dsl.Statement{a, b, c}
	noGoZones := []nogozone.NoGoZone{
		nogozone.NewNoGoZone(geom.NewSegment(50, 99.99), a, c),
	}

	boxes := *NewBoxTracker()

	topOfLifeline := 1.0
	bottomOfLifeline := 100.0

	segments := LifelineSegments{}
	minSegLen := 0.1
	segments.Assemble(b, topOfLifeline, bottomOfLifeline, minSegLen,
		noGoZones, boxes, allLifelines)

	assert.Len(segments.Segs, 1)
	assert.Equal(segments.Segs[0], geom.NewSegment(1.0, 50))
}
