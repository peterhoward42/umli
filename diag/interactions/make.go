package interactions

import (
	"fmt"

	"github.com/peterhoward42/umli"
	"github.com/peterhoward42/umli/diag/lifeline"
	"github.com/peterhoward42/umli/diag/nogozone"
	"github.com/peterhoward42/umli/dsl"
	"github.com/peterhoward42/umli/geom"
	"github.com/peterhoward42/umli/graphics"
	"github.com/peterhoward42/umli/sizer"
)

/*
Maker knows how to make the interaction lines and when to start/stop
their activity boxes.
*/
type Maker struct {
	dependencies  *MakerDependencies
	graphicsModel *graphics.Model
	noGoZones     []nogozone.NoGoZone
}

/*
MakerDependencies encapsulates the prior state of the diagram creation
process at the time that the Make method is called. And includes all
the things the Maker needs from the outside to do its job.
*/
type MakerDependencies struct {
	boxes  map[*dsl.Statement]*lifeline.BoxTracker
	fontHt float64
	sizer  sizer.Sizer
	spacer *lifeline.Spacing
}

// NewMakerDependencies makes a MakerDependencies ready to use.
func NewMakerDependencies(fontHt float64, spacer *lifeline.Spacing,
	sizer sizer.Sizer,
	boxes map[*dsl.Statement]*lifeline.BoxTracker) *MakerDependencies {
	return &MakerDependencies{
		boxes:  boxes,
		fontHt: fontHt,
		sizer:  sizer,
		spacer: spacer,
	}
}

/*
NewMaker initialises a Maker ready to use.
*/
func NewMaker(d *MakerDependencies, gm *graphics.Model) *Maker {
	return &Maker{
		dependencies:  d,
		graphicsModel: gm,
	}
}

/*
ScanInteractionStatements goes through the DSL statements in order, and
works out what graphics are required to represent interaction lines, and
activitiy boxes etc. It advances the tidemark as it goes, and returns the
final resultant tidemark.
*/
func (mkr *Maker) ScanInteractionStatements(
	tidemark float64,
	statements []*dsl.Statement) (newTidemark float64,
	noGoZones []nogozone.NoGoZone, err error) {

	// Build a list of actions to execute depending on the statement
	// keyword.
	actions := []dispatch{}
	for _, s := range statements {
		switch s.Keyword {
		case umli.Dash:
			actions = append(actions, dispatch{mkr.interactionLabel, s})
			actions = append(actions, dispatch{mkr.startToBox, s})
			actions = append(actions, dispatch{mkr.interactionLine, s})
		case umli.Full:
			actions = append(actions, dispatch{mkr.interactionLabel, s})
			actions = append(actions, dispatch{mkr.startFromBox, s})
			actions = append(actions, dispatch{mkr.startToBox, s})
			actions = append(actions, dispatch{mkr.interactionLine, s})
		case umli.Self:
			actions = append(actions, dispatch{mkr.selfLabel, s})
			actions = append(actions, dispatch{mkr.startFromBox, s})
			actions = append(actions, dispatch{mkr.selfLines, s})
		case umli.Stop:
			actions = append(actions, dispatch{mkr.endBox, s})
		}
	}
	var prevTidemark float64 = tidemark
	var updatedTidemark float64
	for _, action := range actions {
		updatedTidemark, err = action.fn(prevTidemark, action.statement)
		if err != nil {
			return -1, nil, fmt.Errorf("actionFn: %v", err)
		}
		prevTidemark = updatedTidemark
	}
	return updatedTidemark, mkr.noGoZones, nil
}

// interactionLabel creates the graphics label that belongs to an interaction
// line.
func (mkr *Maker) interactionLabel(
	tidemark float64, s *dsl.Statement) (newTidemark float64, err error) {
	dep := mkr.dependencies
	sourceLifeline := s.ReferencedLifelines[0]
	destLifeline := s.ReferencedLifelines[1]
	fromX, toX, err := mkr.LifelineCentres(sourceLifeline, destLifeline)
	if err != nil {
		return -1, fmt.Errorf("mkr.LifelineCentres: %v", err)
	}
	labelX, horizJustification := NewLabelPosn(fromX, toX).Get()
	mkr.graphicsModel.Primitives.RowOfStrings(
		labelX, tidemark, dep.fontHt, horizJustification, s.LabelSegments)
	newTidemark = tidemark + float64(len(s.LabelSegments))*
		dep.fontHt + dep.sizer.Get("InteractionLineTextPadB")
	noGoZone := nogozone.NewNoGoZone(
		geom.NewSegment(tidemark, newTidemark),
		sourceLifeline, destLifeline)
	mkr.noGoZones = append(mkr.noGoZones, noGoZone)
	return newTidemark, nil
}

// selfLabel creates the graphics label that belongs to a self interaction
// line.
func (mkr *Maker) selfLabel(
	tidemark float64, s *dsl.Statement) (newTidemark float64, err error) {
	dep := mkr.dependencies
	lifelineXCoords, err := dep.spacer.CentreLine(s.ReferencedLifelines[0])
	if err != nil {
		return -1, fmt.Errorf("spacer.CentreLine: %v", err)
	}
	lineStartX := lifelineXCoords.Centre + 0.5*dep.sizer.Get("ActivityBoxWidth")
	lineEndX := lineStartX + dep.sizer.Get("SelfLoopWidthFactor")*dep.spacer.LifelinePitch()
	labelX := 0.5 * (lineStartX + lineEndX)
	mkr.graphicsModel.Primitives.RowOfStrings(
		labelX, tidemark, dep.fontHt, graphics.Centre, s.LabelSegments)
	htOfLabels := float64(len(s.LabelSegments)) * dep.fontHt
	newTidemark = tidemark + htOfLabels + dep.sizer.Get("InteractionLineTextPadB")
	return newTidemark, nil
}

// interactionLine makes an interaction line (and its arrow)
func (mkr *Maker) interactionLine(
	tidemark float64, s *dsl.Statement) (newTidemark float64, err error) {
	dep := mkr.dependencies
	sourceLifeline := s.ReferencedLifelines[0]
	destLifeline := s.ReferencedLifelines[1]
	fromX, toX, err := mkr.LifelineCentres(sourceLifeline, destLifeline)
	if err != nil {
		return -1, fmt.Errorf("mkr.LifelineCentres: %v", err)
	}
	halfActivityBoxWidth := 0.5 * dep.sizer.Get("ActivityBoxWidth")
	geom.ShortenLineBy(halfActivityBoxWidth, &fromX, &toX)
	y := tidemark
	dashed := s.Keyword == umli.Dash
	mkr.graphicsModel.Primitives.AddLine(fromX, y, toX, y, dashed)
	arrowLen := dep.sizer.Get("ArrowLen")
	arrowWidth := dep.sizer.Get("ArrowWidth")
	arrow := geom.MakeArrow(fromX, toX, y, arrowLen, arrowWidth)
	mkr.graphicsModel.Primitives.AddFilledPoly(arrow)
	newTidemark = tidemark + dep.sizer.Get("InteractionLinePadB")
	noGoZone := nogozone.NewNoGoZone(
		geom.NewSegment(tidemark, newTidemark),
		sourceLifeline, destLifeline)
	mkr.noGoZones = append(mkr.noGoZones, noGoZone)
	return newTidemark, nil
}

// selfLines makes the 3 sides of a rectangle for a self interaction
// line (and its arrow).
func (mkr *Maker) selfLines(
	tidemark float64, s *dsl.Statement) (newTidemark float64, err error) {
	dep := mkr.dependencies
	lifelineXCoords, err := dep.spacer.CentreLine(s.ReferencedLifelines[0])
	if err != nil {
		return -1, fmt.Errorf("spacer.CentreLine: %v", err)
	}
	lineStartX := lifelineXCoords.Centre + 0.5*dep.sizer.Get("ActivityBoxWidth")
	lineEndX := lineStartX + dep.sizer.Get("SelfLoopWidthFactor")*dep.spacer.LifelinePitch()
	y := tidemark
	notDashed := false
	bottom := y + dep.sizer.Get("SelfLoopHeight")
	mkr.graphicsModel.Primitives.AddLine(lineStartX, y, lineEndX, y, notDashed)
	mkr.graphicsModel.Primitives.AddLine(lineEndX, y, lineEndX, bottom, notDashed)
	mkr.graphicsModel.Primitives.AddLine(lineEndX, bottom, lineStartX, bottom, notDashed)
	arrowLen := dep.sizer.Get("ArrowLen")
	arrowWidth := dep.sizer.Get("ArrowWidth")
	arrow := geom.MakeArrow(lineEndX, lineStartX, bottom, arrowLen, arrowWidth)
	mkr.graphicsModel.Primitives.AddFilledPoly(arrow)
	newTidemark = bottom + dep.sizer.Get("InteractionLinePadB")
	return newTidemark, nil
}

// startToBox registers with a lifeline.BoxTracker that an activity box
// on a lifeline should be started ready for an interaction line to arrive at
// the top of it. (If a box is not already in progress for this lifeline.)
func (mkr *Maker) startToBox(
	tidemark float64, s *dsl.Statement) (newTidemark float64, err error) {
	dep := mkr.dependencies
	toLifeline := s.ReferencedLifelines[1]
	boxes := dep.boxes[toLifeline]
	if boxes.HasABoxInProgress() {
		return tidemark, nil
	}
	if err := boxes.AddStartingAt(tidemark); err != nil {
		return -1, fmt.Errorf("boxes.AddStartingAt: %v", err)
	}
	// Return an unchanged tidemark.
	return tidemark, nil
}

// starFromBox registers that an activity box on a lifeline
// should be started (unless it is already) for an activity line to emenate from.
func (mkr *Maker) startFromBox(
	tidemark float64, s *dsl.Statement) (newTidemark float64, err error) {
	dep := mkr.dependencies
	fromLifeline := s.ReferencedLifelines[0]
	boxes := dep.boxes[fromLifeline]
	if boxes.HasABoxInProgress() {
		return tidemark, nil
	}
	// The activity box should start just a tiny bit before the first
	// interaction line leaving from it. This need not claim any vertical
	// space of its own however, because the space already claimed by the interaction
	// line label is sufficient.
	backTrackToStart := dep.sizer.Get("ActivityBoxVerticalOverlap")
	if err := boxes.AddStartingAt(tidemark - backTrackToStart); err != nil {
		return -1, fmt.Errorf("boxes.AddStartingAt: %v", err)
	}
	// Return an unchanged tidemark.
	return tidemark, nil
}

// endBox processes an explicit "stop" statement.
func (mkr *Maker) endBox(
	tidemark float64, s *dsl.Statement) (newTidemark float64, err error) {
	dep := mkr.dependencies
	fromLifeline := s.ReferencedLifelines[0]
	boxes := dep.boxes[fromLifeline]
	err = boxes.TerminateAt(tidemark)
	if err != nil {
		return -1, fmt.Errorf("boxes.TerminateAt: %v", err)
	}
	tidemark += dep.sizer.Get("IndividualStoppedBoxPadB")
	return tidemark, nil
}

// actionFn describes a function that can be called to draw something
// related to statement s. Such a function receives the current tidemark,
// calculates how it should be advanced, and return the updated value.
type actionFn func(
	tideMark float64,
	s *dsl.Statement) (newTidemark float64, err error)

// dispatch is a simple container to hold a binding between  an actionFn and
// the statement to which it refers.
type dispatch struct {
	fn        actionFn
	statement *dsl.Statement
}

/*
LifelineCentres evaluates the X coordinates for the lifelines between which
an interaction line travels.
*/
func (mkr *Maker) LifelineCentres(
	sourceLifeline, destLifeline *dsl.Statement) (fromX, toX float64, err error) {
	fromCoords, err := mkr.dependencies.spacer.CentreLine(sourceLifeline)
	if err != nil {
		return -1.0, -1.0, fmt.Errorf("space.CentreLine: %v", err)
	}
	toCoords, err := mkr.dependencies.spacer.CentreLine(destLifeline)
	if err != nil {
		return -1.0, -1.0, fmt.Errorf("space.CentreLine: %v", err)
	}
	return fromCoords.Centre, toCoords.Centre, nil
}
