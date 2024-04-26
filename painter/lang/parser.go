package lang

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"kpi-apz-lab-3/painter"
	"kpi-apz-lab-3/ui"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct{}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	var res []painter.Operation

	for scanner.Scan() {
		cmd := scanner.Text()
		op, err := parseCmd(cmd)
		if err != nil {
			return nil, err
		}
		res = append(res, op)
	}

	return res, nil
}

func parseCmd(cmd string) (painter.Operation, error) {
	fields := strings.Fields(cmd)
	name := fields[0]
	args, err := parseCmdArgs(fields[1:])
	if err != nil {
		return nil, err
	}

	cmdMapper := map[string]struct {
		handler   interface{}
		operation painter.Operation
	}{
		"white": {
			painter.WhiteFill,
			painter.OperationFunc(painter.WhiteFill),
		},
		"green": {
			painter.GreenFill,
			painter.OperationFunc(painter.GreenFill),
		},
		"update": {
			painter.UpdateOp.Do,
			painter.UpdateOp,
		},
		"bgrect": {
			painter.BgRectDraw,
			painter.OperationFunc(func(s ui.StateSetter) {
				painter.BgRectDraw(s, args[0], args[1], args[2], args[3])
			}),
		},
		"figure": {
			painter.FigureDraw,
			painter.OperationFunc(func(s ui.StateSetter) {
				painter.FigureDraw(s, args[0], args[1])
			}),
		},
		"move": {
			painter.FiguresMove,
			painter.OperationFunc(func(s ui.StateSetter) {
				painter.FiguresMove(s, args[0], args[1])
			}),
		},
		"reset": {
			painter.UIStateReset,
			painter.OperationFunc(painter.UIStateReset),
		},
	}

	h := cmdMapper[name].handler
	if h == nil {
		return nil, fmt.Errorf("unknown command name: %s", name)
	}
	if err = checkArgsNum(args, h); err != nil {
		return nil, err
	}

	op := cmdMapper[name].operation
	return op, nil
}

func parseCmdArgs(args []string) ([]float32, error) {
	res := make([]float32, len(args))
	for i, str := range args {
		num, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return nil, fmt.Errorf("agrument %s is not float", str)
		}
		if num < 0.0 || num > 1.0 {
			return nil, fmt.Errorf("agrument %s is not relative value (it can be in [0.0, 1.0])", str)
		}
		res[i] = float32(num)
	}
	return res, nil
}

func checkArgsNum(args []float32, f interface{}) error {
	paramsNum := reflect.TypeOf(f).NumIn()
	if len(args) != paramsNum-1 { // -1 because of ui.StateSetter param
		err := errors.New("wrong number of arguments")
		return err
	}
	return nil
}
