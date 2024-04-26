package lang

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"kpi-apz-lab-3/painter"
	"kpi-apz-lab-3/ui"
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	for _, tc := range []struct {
		name        string
		input       string
		expectedRes painter.Operation
		expectedErr error
	}{
		{
			name:        "White fill",
			input:       "white\n",
			expectedRes: painter.OperationFunc(painter.WhiteFill),
			expectedErr: nil,
		},
		{
			name:        "Green fill",
			input:       "green\n",
			expectedRes: painter.OperationFunc(painter.GreenFill),
			expectedErr: nil,
		},
		{
			name:        "Figure",
			input:       "figure 0.5 0.6\n",
			expectedRes: painter.OperationFunc(func(s ui.StateSetter) { painter.FigureDraw(s, 0.5, 0.6) }),
			expectedErr: nil,
		},
		{
			name:        "Move",
			input:       "move 0.7 0.8\n",
			expectedRes: painter.OperationFunc(func(s ui.StateSetter) { painter.FiguresMove(s, 0.7, 0.8) }),
			expectedErr: nil,
		},
		{
			name:        "Background rectangle",
			input:       "bgrect 0.1 0.2 0.3 0.4\n",
			expectedRes: painter.OperationFunc(func(s ui.StateSetter) { painter.BgRectDraw(s, 0.1, 0.2, 0.3, 0.4) }),
			expectedErr: nil,
		},
		{
			name:        "Reset",
			input:       "reset\n",
			expectedRes: painter.OperationFunc(painter.UIStateReset),
			expectedErr: nil,
		},
		{
			name:        "Update",
			input:       "update\n",
			expectedRes: painter.UpdateOp,
			expectedErr: nil,
		},
		{
			name:        "Unknown command",
			input:       "123command",
			expectedRes: nil,
			expectedErr: errors.New("unknown command name: 123command"),
		},
		{
			name:        "Command argument is not float",
			input:       "move num 0.3\n",
			expectedRes: nil,
			expectedErr: errors.New("argument num is not float"),
		},
		{
			name:        "Command argument is not relative value",
			input:       "move 3 0.4\n",
			expectedRes: nil,
			expectedErr: errors.New("argument 3 is not relative value (it can be in [0.0, 1.0])"),
		},
		{
			name:        "Wrong number of arguments",
			input:       "figure 0.5 0.6 0.3\n",
			expectedRes: nil,
			expectedErr: errors.New("wrong number of arguments"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			parser := &Parser{}
			res, err := parser.Parse(strings.NewReader(tc.input))
			if tc.expectedErr == nil {
				require.NoError(t, err)
				require.Len(t, res, 1)
				assert.IsType(t, tc.expectedRes, res[0])
			} else {
				assert.Nil(t, res)
				assert.EqualError(t, err, tc.expectedErr.Error())
			}
		})
	}
}
