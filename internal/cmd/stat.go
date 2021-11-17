package cmd

import (
	"bytes"
	"io"
	"sort"
	"strconv"
	"time"

	"github.com/nestoroprysk/mood-tracker/internal/registry"
	chart "github.com/wcharczuk/go-chart/v2"
	"golang.org/x/sync/errgroup"
)

func newStat(e env) (Cmd, error) {
	return func() (string, error) {
		b, err := e.Read(userIDJSON(e.userID))
		if err != nil {
			return "", err
		}

		r, err := registry.Make(b, e.userID)
		if err != nil {
			return "", err
		}

		var g errgroup.Group

		for _, i := range []struct {
			name string
			f    func(registry.T) (io.Reader, error)
		}{
			{name: "time.png", f: timePNG},
			{name: "label.png", f: labelPNG},
			{name: "freq.png", f: freqPNG},
		} {
			f, name := i.f, i.name
			g.Go(func() error {
				png, err := f(r)
				if err != nil {
					return err
				}

				if _, err := e.SendPNG(name, png); err != nil {
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

func timePNG(r registry.T) (io.Reader, error) {
	var xs []time.Time
	var ys []float64
	for _, i := range r.Items {
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

func labelPNG(r registry.T) (io.Reader, error) {
	m := map[string]int{}
	for _, i := range r.Items {
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

	sort.Slice(vals, func(i, j int) bool { return vals[i].Value >= vals[j].Value })

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

func freqPNG(r registry.T) (io.Reader, error) {
	m := map[int]int{
		1: 0,
		2: 0,
		3: 0,
		4: 0,
		5: 0,
	}
	for _, i := range r.Items {
		m[i.Mood]++
	}

	var vals []chart.Value
	for k, v := range m {
		vals = append(vals, chart.Value{
			Label: strconv.Itoa(k),
			Value: float64(v),
		})
	}

	sort.Slice(vals, func(i, j int) bool { return vals[i].Label >= vals[j].Label })

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
