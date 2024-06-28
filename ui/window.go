package ui

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

var yellowColor = color.RGBA{R: 0xff, G: 0xff, B: 0xf, A: 0xff}

type MyFigure struct {
	x, y, w, h int
}

func GetFigure(myX int, myY int) *MyFigure {
	return &MyFigure{x: myX, y: myY, w: 100, h: 300}
}

func (f *MyFigure) Draw(t screen.Texture) {
	x1 := f.x - f.w/2
	y1 := f.y - f.h/2
	x2 := f.x + f.w/2
	y2 := f.y + f.h/2
	t.Fill(image.Rect(x1, y1, x2, y2), yellowColor, draw.Src)

	x1 = f.x - f.h/2
	y1 = f.y - f.w/2
	x2 = f.x + f.h/2
	y2 = f.y + f.w/2
	t.Fill(image.Rect(x1, y1, x2, y2), yellowColor, draw.Src)
}

func (f *MyFigure) VisualizeFigure(pw *Visualizer) {
	x1 := f.x - f.w/2
	y1 := f.y - f.h/2
	x2 := f.x + f.w/2
	y2 := f.y + f.h/2
	pw.w.Fill(image.Rect(x1, y1, x2, y2), yellowColor, draw.Src)

	x1 = f.x - f.h/2
	y1 = f.y - f.w/2
	x2 = f.x + f.h/2
	y2 = f.y + f.w/2
	pw.w.Fill(image.Rect(x1, y1, x2, y2), yellowColor, draw.Src)
}

func (f *MyFigure) Move(x, y int) {
	f.x = x
	f.y = y
}

func (pw* Visualizer) MoveAllFigures(x, y int) {
	for _, f := range pw.Figures {
		f.Move(x, y)
	}
}

func (pw* Visualizer) AddFigure(t screen.Texture) {
	pw.Figures = append(pw.Figures, GetFigure(100, 100))
}

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	w    screen.Window
	tx   chan screen.Texture
	done chan struct{}

	sz  size.Event
	pos image.Rectangle
	Figures   []*MyFigure
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})
	pw.pos.Max.X = 200
	pw.pos.Max.Y = 200
	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.tx <- t
}

func (pw *Visualizer) run(s screen.Screen) {
	w, err := s.NewWindow(&screen.NewWindowOptions{
		Title: pw.Title,
		Width: 800,
		Height: 800,
	})
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		w.Release()
		close(pw.done)
	}()

	if pw.OnScreenReady != nil {
		pw.OnScreenReady(s)
	}

	pw.w = w
	pw.Figures = []*MyFigure{GetFigure(400, 400)}

	events := make(chan any)
	go func() {
		for {
			e := w.NextEvent()
			if pw.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var t screen.Texture

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, t)

		case t = <-pw.tx:
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true // Window destroy initiated.
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true // Esc pressed.
		}
	}
	return false
}

func (pw *Visualizer) handleEvent(e any, t screen.Texture) {
	switch e := e.(type) {

	case size.Event: // Оновлення даних про розмір вікна.
		pw.sz = e

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if t != nil { return; }
		if e.Button != mouse.ButtonRight { return; }
		if e.Direction != mouse.DirPress { return; }
		
		pw.MoveAllFigures(int(e.X), int(e.Y))
		pw.w.Send(paint.Event{})

	case paint.Event:
		// Малювання контенту вікна.
		if t == nil {
			pw.drawDefaultUI()
		} else {
			// Використання текстури отриманої через виклик Update.
			pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
		}
		pw.w.Publish()
	}
}

func (pw *Visualizer) drawDefaultUI() {
	pw.w.Fill(pw.sz.Bounds(), color.Black, draw.Src) // Фон.

	for _, f := range pw.Figures {
		f.VisualizeFigure(pw)
	}

	// Малювання білої рамки.
	for _, br := range imageutil.Border(pw.sz.Bounds(), 10) {
		pw.w.Fill(br, color.White, draw.Src)
	}
}
