package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	chart "github.com/wcharczuk/go-chart/v2"
)

func newStat(c config) (Cmd, error) {
	return func() (string, error) {
		b, err := c.Read(userIDJSON(c.userID))
		if err != nil {
			return "", err
		}

		var r Registry
		if err := json.Unmarshal(b, &r); err != nil {
			return "", err
		}

		png, err := statPNG(r)
		if err != nil {
			return "", err
		}

		if _, err := c.SendPNG("stat.png", png); err != nil {
			return "", err
		}

		return "", nil
	}, nil
}

func statPNG(r Registry) (io.Reader, error) {
	var xs []time.Time
	var ys []float64
	for _, i := range r {
		xs = append(xs, i.Time)
		ys = append(ys, float64(i.Mood))
	}

	graph := chart.Chart{
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: xs,
				YValues: ys,
			},
		},
	}

	b := &bytes.Buffer{}
	if err := graph.Render(chart.PNG, b); err != nil {
		return nil, err
	}

	return b, nil
}
