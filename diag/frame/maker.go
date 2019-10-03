package frame

import (
	"github.com/peterhoward42/umli/graphics"
	"github.com/peterhoward42/umli/sizer"
)

/*
maker knows how to draw the diagram outer frame and title box into a
graphics.Primitives structure.
*/
type Maker struct {
	sizer    sizer.Sizer
	frameTop float64
	prims    *graphics.Primitives
}

// newMaker provides a lifelineBoxes ready to use.
func NewMaker(s sizer.Sizer, prims *graphics.Primitives) *Maker {
	return &Maker{
		sizer: s,
		prims: prims,
	}
}

/*
initFrameAndMakeTitleBox is responsible capturing the Y coordinate at which
the diagram's frame rectangle should start, and then drawing the diagram title
in an enclosing rectangle just below it. Then advancing the tidemark
accordingly.
*/
func (fm *Maker) InitFrameAndMakeTitleBox(titleSegments []string,
	frameTop float64) (newTideMark float64) {
	fm.frameTop = frameTop
	tideMark := frameTop + fm.sizer.Get("FrameTitleTextPadT")
	topOfTitleTextY := tideMark
	leftOfText := fm.sizer.Get("FramePadLR") + fm.sizer.Get("FrameTitleTextPadL")
	fm.prims.RowOfStrings(leftOfText, topOfTitleTextY,
		fm.sizer.Get("FontHeight"), graphics.Left, titleSegments)
	tideMark += float64(len(titleSegments)) * fm.sizer.Get("FontHeight")
	tideMark += fm.sizer.Get("FrameTitleTextPadB")
	rightOfFrameTitleBox := fm.sizer.Get("FrameTitleBoxWidth")
	fm.prims.AddRect(fm.sizer.Get("FramePadLR"), fm.frameTop,
		rightOfFrameTitleBox, tideMark)
	tideMark += fm.sizer.Get("FrameTitleRectPadB")
	return tideMark
}

/*
finalizeFrame claims a little space below the diagram vertical extent so far,
and draws the enclosing frame. It is not responsible for reserving space,
below the frame - that is handled externally.
*/
func (fm *Maker) FinalizeFrame(currentTideMark float64,
	diagWidth float64) (newTideMark float64) {
	currentTideMark += fm.sizer.Get("FrameInternalPadB")
	frameBottom := currentTideMark
	left := fm.sizer.Get("FramePadLR")
	right := float64(diagWidth) - fm.sizer.Get("FramePadLR")
	fm.prims.AddRect(left, fm.frameTop, right, frameBottom)
	return currentTideMark
}
