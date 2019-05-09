package graphics

// Line represents a line, optionally dashed, and optionally with an arrow at
// the (X2, Y2) end.
type Line struct {
	X1, Y1, X2, Y2 float64
	Arrow          bool
	Dashed         bool // vs. Full
}

// The values that may be used in label justification.
const (
	Left   = "Left"
	Right  = "Right"
	Top    = "Top"
	Bottom = "Bottom"
	Centre = "Centre"
)

// Label encapsulates a (potentially) multi-line label in terms of a position,
// justification and its consituent lines of text.
type Label struct {
	LinesOfText []string
	// Anchor point about which the justifications are applied
	X, Y  float64
	HJust string
	VJust string
}

// Primitives is a container for a set of Line(s) and a set of Label(s).
type Primitives struct {
	Lines  []*Line
	Labels []*Label
}

// NewPrimitives constructs a Primitives ready to use.
func NewPrimitives() *Primitives {
	return &Primitives{[]*Line{}, []*Label{}}
}

// AddLine adds the given line to the Primitive's line store.
func (p *Primitives) AddLine(x1 float64, y1 float64, x2 float64, y2 float64,
	dashed bool, arrow bool) {
	line := &Line{x1, y1, x2, y2, dashed, arrow}
	p.Lines = append(p.Lines, line)
}

// AddLabel adds a Lable to the Primitive's Lable store.
func (p *Primitives) AddLabel(linesOfText []string, x float64, y float64,
	hJust string, vJust string) {
	label := &Label{linesOfText, x, y, hJust, vJust}
	p.Labels = append(p.Labels, label)
}

// AddRect adds 4 lines to the Primitive's line store to represent
// the rectangle of the given opposite corners.
func (p *Primitives) AddRect(
	left float64, top float64, right float64, bot float64) {
	p.AddLine(left, top, right, top, false, false)
	p.AddLine(right, top, right, bot, false, false)
	p.AddLine(right, bot, left, bot, false, false)
	p.AddLine(left, bot, left, top, false, false)
}

// Add adds the Primitives given to those already held in the model.
func (p *Primitives) Add(newPrims *Primitives) {
	p.Lines = append(p.Lines, newPrims.Lines...)
	p.Labels = append(p.Labels, newPrims.Labels...)
}
