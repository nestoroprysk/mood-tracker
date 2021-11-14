package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"sort"
	"strconv"
	"time"

	chart "github.com/wcharczuk/go-chart/v2"
	"golang.org/x/sync/errgroup"
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

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
		defer cancel()
		g, ctx := errgroup.WithContext(ctx)

		for _, e := range []struct {
			name string
			f    func(Registry) (io.Reader, error)
		}{
			{name: "time.png", f: timePNG},
			{name: "label.png", f: labelPNG},
			{name: "freq.png", f: freqPNG},
		} {
			f, name := e.f, e.name
			g.Go(func() error {
				png, err := f(r)
				if err != nil {
					return err
				}

				if _, err := c.SendPNG(name, png); err != nil {
					return err
				}

				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return "", err
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

	sort.Slice(vals, func(i, j int) bool { return vals[i].Label < vals[j].Label })

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

	sort.Slice(vals, func(i, j int) bool { return vals[i].Label < vals[j].Label })

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
