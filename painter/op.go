package painter

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	"github.com/roman-mazur/architecture-lab-3/ui"
	"golang.org/x/exp/shiny/screen"
)

func conv(args []string) ([]int, error) {
	parsedValues := make([]int, len(args))
	for i, str := range args {
		num, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		parsedValues[i] = int(num)
	}
	return parsedValues, nil
}

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture, pv *ui.Visualizer) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(t screen.Texture, state *ui.Visualizer) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t, state) || ready
	}
	return
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture, state *ui.Visualizer) bool {
	state.Update(t)
	return true
}

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(t screen.Texture, state *ui.Visualizer)

func (f OperationFunc) Do(t screen.Texture, state *ui.Visualizer) bool {
	f(t, state)
	return false
}

// WhiteFill зафарбовує тестуру у білий колір. Може бути викоистана як Operation через OperationFunc(WhiteFill).
func WhiteFill(t screen.Texture, state *ui.Visualizer) {
	state.Background = color.White
}

// GreenFill зафарбовує тестуру у зелений колір. Може бути викоистана як Operation через OperationFunc(GreenFill).
func GreenFill(t screen.Texture, state *ui.Visualizer) {
	state.Background = color.RGBA{0, 255, 0, 255}
}

func Reset(t screen.Texture, state *ui.Visualizer) {
	state.Background = color.Black
	state.Rect = [2]image.Point{{0, 0}, {0, 0}}
	state.Figures = []*ui.MyFigure{}
}

func DrawRectangle(args []string) OperationFunc {
	if len(args) != 4 {
		fmt.Println("Wrong amount of arguments to draw a rectangle")
		return nil
	}
	cords, err := conv(args)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return func(t screen.Texture, state *ui.Visualizer) {
		if err == nil && len(cords) == 4 {
			state.Rect[0] = image.Point{int(cords[0]), int(cords[1])}
			state.Rect[1] = image.Point{int(cords[2]), int(cords[3])}
		}
	}
}

func Figure(args []string) OperationFunc {
	if len(args) != 2 {
		fmt.Println("Wrong amount of arguments to draw figures")
		return nil
	}
	cords, err := conv(args)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return func(t screen.Texture, state *ui.Visualizer) {
		if err == nil && len(cords) == 2 {
			f := ui.GetFigure(cords[0], cords[1])
			state.Figures = append(state.Figures, f)
		}
	}
}

func Move(args []string) OperationFunc {
	if len(args) != 2 {
		fmt.Println("Wrong amount of arguments to move figures")
		return nil
	}
	cords, err := conv(args)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return func(t screen.Texture, state *ui.Visualizer) {
		if err == nil && len(cords) == 2 {
			f := ui.GetFigure(cords[0], cords[1])
			state.Figures = []*ui.MyFigure{f}
		}
	}
}
