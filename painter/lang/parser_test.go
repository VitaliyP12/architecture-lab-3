package lang

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

func TestParser_CommandParser(t *testing.T) {
	p := &Parser{}

	tests := []struct {
		name         string
		commandName  string
		args         []string
		expectedFunc painter.OperationFunc
	}{
		{
			name:         "white fill",
			commandName:  "white",
			expectedFunc: painter.OperationFunc(painter.WhiteFill),
		},
		{
			name:         "green fill",
			commandName:  "green",
			expectedFunc: painter.OperationFunc(painter.GreenFill),
		},
		{
			name:         "draw rectangle",
			commandName:  "bgrect",
			args:         []string{"5", "6", "30", "40"},
			expectedFunc: painter.DrawRectangle([]string{"5", "6", "30", "40"}),
		},
		{
			name:         "figure",
			commandName:  "figure",
			args:         []string{"25", "25"},
			expectedFunc: painter.Figure([]string{"25", "25"}),
		},
		{
			name:         "move",
			commandName:  "move",
			args:         []string{"56", "57"},
			expectedFunc: painter.Move([]string{"56", "57"}),
		},
		{
			name:         "reset",
			commandName:  "reset",
			expectedFunc: painter.OperationFunc(painter.Reset),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fmt.Println(test.commandName)
			result := p.ParseCommands(test.commandName, test.args)
			if result == nil {
				t.Errorf("Expected non-nil result for command %s", test.commandName)
			}
			if func1, ok := result.(painter.OperationFunc); ok {
				if !painterOperationFuncEquals(func1, test.expectedFunc) {
					t.Errorf("Expected function %v, got %v", test.expectedFunc, func1)
				}
			}
		})
	}
}

func TestParser_Parse(t *testing.T) {
	p := &Parser{}

	tests := []struct {
		name          string
		input         string
		expectedCount int
	}{
		{
			name:          "single command",
			input:         "white\n",
			expectedCount: 1,
		},
		{
			name:          "multiple commands",
			input:         "white\ngreen\n",
			expectedCount: 2,
		},
		{
			name:          "command with args",
			input:         "bgrect 55 45 30 50\n",
			expectedCount: 1,
		},
		{
			name:          "no command",
			input:         "",
			expectedCount: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := strings.NewReader(test.input)
			ops, err := p.Parse(r)
			if err != nil {
				t.Fatalf("Error parsing input: %v", err)
			}
			if len(ops) != test.expectedCount {
				t.Errorf("Expected %d operations, got %d", test.expectedCount, len(ops))
			}
		})
	}
}

func painterOperationFuncEquals(f1, f2 painter.OperationFunc) bool {
	return reflect.ValueOf(f1).Pointer() == reflect.ValueOf(f2).Pointer()
}
