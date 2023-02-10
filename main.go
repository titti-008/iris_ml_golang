package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/goml/gobrain"
	"github.com/goml/gobrain/persist"
)

const (
	DATA_FILE = "iris.csv"
	MODEL     = "model.json"
)

func loadData() ([][]float64, []string, error) {
	f, err := os.Open(DATA_FILE)
	if err != nil {
		return nil, nil, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	// read header
	scanner.Scan()

	var xResult [][]float64
	var yResult []string
	for scanner.Scan() {
		var f1, f2, f3, f4 float64
		var s string
		n, err := fmt.Sscanf(scanner.Text(), "%f,%f,%f,%f,%s", &f1, &f2, &f3, &f4, &s)
		if n != 5 || err != nil {
			return nil, nil, err
		}
		xResult = append(xResult, []float64{f1, f2, f3, f4})
		yResult = append(yResult, strings.Trim(s, `"`))
	}

	return xResult, yResult, err

}

func shuffle(x [][]float64, y []string) {
	for i := len(x) - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		x[i], x[j] = x[j], x[i]
		y[i], y[j] = y[j], y[i]
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Start ML!")

	X, Y, err := loadData()
	if err != nil {
		log.Fatal(err)
	}

	shuffle(X, Y)

	n := 100
	xTrain, yTrain, xTest, yTest := X[:n], Y[:n], X[n:], Y[n:]

	m := map[string][]float64{
		"Setosa":     {1, 0, 0},
		"Versicolor": {0, 1, 0},
		"Virginica":  {0, 0, 1},
	}

	patterns := [][][]float64{}

	for i, x := range xTrain {
		patterns = append(patterns, [][]float64{x, m[yTrain[i]]})
	}

	ff := &gobrain.FeedForward{}
	ff.Init(4, 3, 3)

	err = persist.Load(MODEL, &ff)
	if err != nil {
		ff.Train(patterns, 100000, 0.6, 0.04, true)
		persist.Save(MODEL, &ff)
	}

	var success int

	for i, x := range xTest {
		result := ff.Update(x)

		var mf float64
		var mj int

		for j, v := range result {
			if v > mf {
				mf = v
				mj = j
			}
		}

		want := yTest[i]
		got := []string{"Setosa", "Versicolor", "Virginica"}[mj]
		fmt.Printf("want: %q, got: %q\n", want, got)

		if want == got {
			success++
		}
	}

	fmt.Printf("%.2f%%\n", float64(success)/float64(len(xTest))*100)
}
