package diag

/*
This module owns knowledge about what graphical drawing events should be
triggered as each type of DSL statement is encountered,
*/

import (
	"github.com/peterhoward42/umli"
)

// EventType is the enumerated-type for the constants such as EndBox or
// LifelineLine below.
type EventType int

// These constants comprise the set of values for EventType.
const (
	EndBox EventType = iota + 1
	InteractionLine
	InteractionLabel
	LifelineLine
	LifelineTitleBox
	SelfInteractionLines
	SelfInteractionLabel
	PotentiallyStartFromBox
	PotentiallyStartToBox
)

/*
EventsRequired provides the list of EventType(s) that should be stimulated
in response to each DSL keyword.

The sequence within these lists is significant, because each event
*claims* a certain amount
of vertical room for itself, which then *pushes* everything that follows
further down the diagram.

In this context, the labels for interaction lines and for self interaction
lines, will be drawn above the lines to which they refer, and therefore must
precede the corresponding line events.
*/
var EventsRequired = map[string][]EventType{
	umli.Life: []EventType{
		LifelineTitleBox, // advances tidemark
		LifelineLine,
	},
	umli.Dash: []EventType{
		InteractionLabel,      // advances tidemark
		PotentiallyStartToBox, // no advance
		InteractionLine,       // advances tidemark
	},
	umli.Full: []EventType{
		InteractionLabel,        // advances tidemark
		PotentiallyStartFromBox, // no advance (renders behind tidemark)
		PotentiallyStartToBox,   // no advance
		InteractionLine,         // advances tidemark
	},
	umli.Self: []EventType{
		PotentiallyStartFromBox, // no advance (renders behind tidemark)
		SelfInteractionLines,    // advances tidemark (includes label inside loop)
	},
	umli.Stop: []EventType{
		EndBox,
	},
}
