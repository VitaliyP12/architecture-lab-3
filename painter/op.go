package painter

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	"github.com/roman-mazur/architecture-lab-3/ui"
	"golang.org/x/exp/shiny/screen"
)

func conv(args []string) ([]float64, error) {
	parsedValues := make([]float64, len(args))
	for i, str := range args {
		num, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, err
		}
		parsedValues[i] = num
	}
	return parsedValues, nil
}

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture, state *CurState) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(t screen.Texture, state *CurState) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t, state) || ready
	}
	return
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture, state *CurState) bool {
	t.Fill(t.Bounds(), state.background, screen.Src)
	t.Fill(image.Rectangle{state.bgRect[0], state.bgRect[1]}, color.Black, screen.Src)
	for _, item := range state.Figures {
		item.Draw(t)
	}
	return true
}

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(t screen.Texture, state *CurState)

func (f OperationFunc) Do(t screen.Texture, state *CurState) bool {
	f(t, state)
	return false
}

// WhiteFill зафарбовує тестуру у білий колір. Може бути викоистана як Operation через OperationFunc(WhiteFill).
func WhiteFill(t screen.Texture, state *CurState) {
	t.Fill(t.Bounds(), color.White, screen.Src)
}

// GreenFill зафарбовує тестуру у зелений колір. Може бути викоистана як Operation через OperationFunc(GreenFill).
func GreenFill(t screen.Texture, state *CurState) {
	t.Fill(t.Bounds(), color.RGBA{G: 0xff, A: 0xff}, screen.Src)
}

func Reset(t screen.Texture, state *CurState) {
	state.background = color.Black
	state.bgRect = [2]image.Point{{0, 0}, {0, 0}}
	state.Figures = []*ui.MyFigure{}
}

func DrawRectangle(args []string) OperationFunc {
	if len(args) != 4 {
		fmt.Println("Wrong amount of arguments to draw a rectangle")
		return nil
	}
	floatArgs, err := conv(args)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return func(t screen.Texture, state *CurState) {
		cords, err := convToCoords(t.Bounds().Dx(), t.Bounds().Dy(), floatArgs)
		if err == nil && len(cords) == 4 {
			state.bgRect[0] = image.Point{int(cords[0]), int(cords[1])}
			state.bgRect[1] = image.Point{int(cords[2]), int(cords[3])}
		}
	}
}

func Figure(args []string) OperationFunc {
	if len(args) != 2 {
		fmt.Println("Wrong amount of arguments to draw figures")
		return nil
	}
	floatArgs, err := conv(args)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return func(t screen.Texture, state *CurState) {
		cords, err := convToCoords(t.Bounds().Dx(), t.Bounds().Dy(), floatArgs)
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
	floatArgs, err := conv(args)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return func(t screen.Texture, state *CurState) {
		cords, err := convToCoords(t.Bounds().Dx(), t.Bounds().Dy(), floatArgs)
		if err == nil && len(cords) == 2 {
			f := ui.GetFigure(cords[0], cords[1])
			state.Figures = []*ui.MyFigure{f}
		}
	}
}

func convToCoords(width int, height int, floatArgs []float64) ([]int, error) {
	if len(floatArgs)%2 != 0 {
		return nil, fmt.Errorf("Wrong amount of arguments!")
	}

	cords := make([]int, len(floatArgs))

	fWidth := float64(width)
	fHeight := float64(height)

	for index := range floatArgs {
		if index%2 == 0 {
			cords[index] = int(fWidth * floatArgs[index])
		} else {
			cords[index] = int(fHeight * floatArgs[index])
		}
	}

	return cords, nil
}
