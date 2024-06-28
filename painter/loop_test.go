package painter

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"

	"github.com/roman-mazur/architecture-lab-3/ui"
	"golang.org/x/exp/shiny/screen"
)

func TestLoop_Post(t *testing.T) {
	var (
		l  Loop
		tr testReceiver
	)
	l.Receiver = &tr
	pw := &ui.Visualizer{
		Background: color.RGBA{0, 0, 0, 0}, 
		Rect: [2]image.Point{{0, 0}, {0, 0}}, 
		Figures: []*ui.MyFigure{},
	}

	

	l.Connect(pw)

	var testOps []string

	l.Start(mockScreen{})
	l.Post(logOp(t, "do green fill", GreenFill))
	l.Post(logOp(t, "do white fill", WhiteFill))
	for i := 0; i < 3; i++ {
		go l.Post(logOp(t, "do green fill", GreenFill))
	}

	l.Post(OperationFunc(func(screen.Texture, *ui.Visualizer) {
		testOps = append(testOps, "op 1")
		l.Post(OperationFunc(func(screen.Texture, *ui.Visualizer) {
			testOps = append(testOps, "op 2")
		}))
	}))
	l.Post(OperationFunc(func(screen.Texture, *ui.Visualizer) {
		testOps = append(testOps, "op 3")
	}))

	l.StopAndWait()

	if !reflect.DeepEqual(testOps, []string{"op 1", "op 3", "op 2"}) {
		t.Error("Bad order:", testOps)
	}

	fmt.Println("Closed")
}

func logOp(t *testing.T, msg string, op OperationFunc) OperationFunc {
	return func(tx screen.Texture, state *ui.Visualizer) {
		t.Log(msg)
		op(tx, state)
	}
}

type testReceiver struct {
	lastTexture screen.Texture
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.lastTexture = t
}

type mockScreen struct{}

// NewBuffer implements screen.Screen.
func (m mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	panic("unimplemented")
}

// NewWindow implements screen.Screen.
func (m mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	panic("unimplemented")
}

func (m mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return new(mockTexture), nil
}

type mockTexture struct {
	Colors []color.Color
}

func (m *mockTexture) Release() {}

func (m *mockTexture) Size() image.Point { return size }

func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: m.Size()}
}

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}
func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.Colors = append(m.Colors, src)
}
