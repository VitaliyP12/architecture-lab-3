package lang

import (
	"bufio"
	"io"
	"strings"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
}

func recOp(operation painter.OperationFunc) painter.Operation {
	if operation == nil {
		return nil
	}
	return painter.OperationFunc(operation)
}

func (p *Parser) ParseCommands(name string, args []string) painter.Operation {
	switch name {
	case "white":
		return painter.OperationFunc(painter.WhiteFill)
	case "green":
		return painter.OperationFunc(painter.GreenFill)
	case "bgrect":
		return recOp(painter.DrawRectangle(args))
	case "figure":
		return recOp(painter.Figure(args))
	case "move":
		return recOp(painter.Move(args))
	case "update":
		return painter.UpdateOp
	case "reset":
		return painter.OperationFunc(painter.Reset)
	}
	return nil
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	var res []painter.Operation

	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		sliced := strings.Split(line, " ")
		args := sliced[1:]
		com := p.ParseCommands(sliced[0], args)

		if com == nil {
			continue
		}

		res = append(res, com)
	}

	//res = append(res, painter.OperationFunc(painter.WhiteFill))
	//res = append(res, painter.UpdateOp)

	return res, nil
}
