package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
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

		for _, e := range []struct {
			name string
			f    func(Registry) (io.Reader, error)
		}{
			{name: "time.png", f: timePNG},
			{name: "label.png", f: labelPNG},
			{name: "freq.png", f: freqPNG},
		} {
			png, err := e.f(r)
			if err != nil {
				return "", err
			}

			if _, err := c.SendPNG(e.name, png); err != nil {
				return "", err
			}
		}

		return "Enjoy!", nil
	}, nil
}

func timePNG(r Registry) (io.Reader, error) {
	var xs []time.Time
	var ys []float64
	for _, i := range r {
		xs = append(xs, i.Time)
		ys = append(ys, float64(i.Mood))
	}

	graph := chart.Chart{
		Title: "Mood Values",
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

func labelPNG(r Registry) (io.Reader, error) {
	m := map[string]int{}
	for _, i := range r {
		for _, l := range i.Labels {
			m[l]++
		}
	}

	var vals []chart.Value
	for k, v := range m {
		vals = append(vals, chart.Value{
			Label: k,
			Value: float64(v),
		})
	}

	graph := chart.BarChart{
		Title: "Mood Labels",
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:   512,
		BarWidth: 60,
		Bars:     vals,
	}

	b := &bytes.Buffer{}
	if err := graph.Render(chart.PNG, b); err != nil {
		return nil, err
	}

	return b, nil
}

func freqPNG(r Registry) (io.Reader, error) {
	m := map[int]int{
		0: 0,
		1: 0,
		2: 0,
		3: 0,
		4: 0,
		5: 0,
	}
	for _, i := range r {
		m[i.Mood]++
	}

	var vals []chart.Value
	for k, v := range m {
		vals = append(vals, chart.Value{
			Label: strconv.Itoa(k),
			Value: float64(v),
		})
	}

	graph := chart.BarChart{
		Title: "Mood Frequencies",
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:   512,
		BarWidth: 60,
		Bars:     vals,
	}

	b := &bytes.Buffer{}
	if err := graph.Render(chart.PNG, b); err != nil {
		return nil, err
	}

	return b, nil
}
