package lang

import (
	"bufio"
	"io"
	"kpi-apz-lab-3/ui"
	"strconv"
	"strings"

	"kpi-apz-lab-3/painter"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct{}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	var res []painter.Operation

	for scanner.Scan() {
		cmd := scanner.Text()
		op := parseCmd(cmd)
		res = append(res, op)
	}

	return res, nil
}

func parseCmd(cmd string) painter.Operation {
	fields := strings.Fields(cmd)
	name := fields[0]
	args := parseCmdArgs(fields[1:])

	return map[string]painter.Operation{
		"white":  painter.OperationFunc(painter.WhiteFill),
		"green":  painter.OperationFunc(painter.GreenFill),
		"update": painter.UpdateOp,
		"bgrect": painter.OperationFunc(func(s ui.State) {
			painter.BgRectDraw(s, args[0], args[1], args[2], args[3])
		}),
		"figure": painter.OperationFunc(func(s ui.State) {
			painter.FigureDraw(s, args[0], args[1])
		}),
		"move": painter.OperationFunc(func(s ui.State) {
			painter.FiguresMove(s, args[0], args[1])
		}),
		"reset": painter.OperationFunc(painter.UiStateReset),
	}[name]
}

func parseCmdArgs(args []string) []float32 {
	if len(args) > 0 {
		res := make([]float32, len(args))
		for i, str := range args {
			num, _ := strconv.ParseFloat(str, 32)
			res[i] = float32(num)
		}
		return res
	}
	return nil
}
