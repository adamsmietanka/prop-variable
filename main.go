package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type props struct {
	MaxSpeed float64 `json:"max_speed"`
	StepSize float64 `json:"step_size"`
	PropSpeed float64 `json:"prop_speed"`
	Diameter float64 `json:"diameter"`
	Blades int `json:"blades"`
	Cp float64 `json:"cp"`
	Power float64 `json:"power"`
	Ratio float64 `json:"ratio"`
}

type tableRow struct {
	V, J, Cp, Rpm, Angle, Eff, Power float64
}

type surfaceChart struct {
	Colorscale [][]interface{} `json:"colorscale"`
	Hovertemplate string `json:"hovertemplate"`
	Opacity float64 `json:"opacity"`
	Showscale bool `json:"showscale"`
	Type string `json:"type"`
	X []float64 `json:"x"`
	Y []float64 `json:"y"`
	Z [][]float64 `json:"z"`
}

type scatterChart struct {
	Marker marker `json:"marker"`
	Hovertemplate string `json:"hovertemplate"`
	Line line `json:"line"`
	Opacity float64 `json:"opacity"`
	Type string `json:"type"`
	X []float64 `json:"x"`
	Y []float64 `json:"y"`
	Z []float64 `json:"z"`
}

type marker struct {
	Color []string `json:"color"`
	Size float64 `json:"size"`
}

type line struct {
	Color string `json:"color"`
}

func readCsvFile(filePath string) [][]float64 {
    f, err := os.Open(filePath)
    if err != nil {
        log.Fatal("Unable to read input file " + filePath, err)
    }
    defer f.Close()

    csvReader := csv.NewReader(f)
    records, err := csvReader.ReadAll()
    if err != nil {
        log.Fatal("Unable to parse file as CSV for " + filePath, err)
    }
	values, _ := convertToFloat(records)

    return values
}

func convertToFloat(records [][]string) ([][]float64, error) {
    values := make([][]float64, len(records))
    for i := range values {
        values[i] = make([]float64, len(records[i]))
    }
    for i := range records {
        for j := range records[i] {
			if records[i][j] != "" {
				val, err := strconv.ParseFloat(records[i][j], 64)
				values[i][j] = val
				if err != nil {
					log.Fatal("Unable to convert array to float", err)
				}
			}
		}
	}
	return values, nil
}

var maxZ = map[int]float64{
	2: 0.44356605705079777,
	3: 0.619386,
	4: 0.8102169239607058,
}

func loadData(p props) ([][]float64, [][]float64, [][]float64) {
	cpY := make([]float64, 51)
	Z := make([]float64, 101)
	effY := make([]float64, 101)
	for i := range cpY {
		cpY[i] = 10 + float64(i)
	}
	for i := range Z {
		Z[i] = float64(i) * maxZ[p.Blades] / 100
	}
	for i := range effY {
		effY[i] = 10 + 0.5 * float64(i)
	}
	cpFile := fmt.Sprintf("data/clark_%d_cp.csv", p.Blades)
	cpXZFile := fmt.Sprintf("data/xz_clark_%d_cp.csv", p.Blades)
	effFile := fmt.Sprintf("data/clark_%d_eff.csv", p.Blades)
	cp := readCsvFile(cpFile)
	cpXZ := readCsvFile(cpXZFile)
	eff := readCsvFile(effFile)
    return append(cp, cpY), append(cpXZ, Z), append(eff, effY) 
}

func getTable(p props, cpXZ, eff [][]float64) []tableRow {
	length := int(1.2 * p.MaxSpeed / p.StepSize)
	table := make([]tableRow, length)
	v := float64(0)
	for i := range table {
		j := v / (p.PropSpeed * p.Diameter)
		cp := p.Cp
		rpm := p.PropSpeed * 60 / p.Ratio
		angle := InterpolateY(cpXZ[0], cpXZ[len(cpXZ) - 1], cpXZ[1:len(cpXZ) - 1], j, cp)
		eff := 	 InterpolateZ(eff[0] , eff[len(eff) - 1]  , eff[1:len(eff) - 1]  , j, angle)
		power := p.Power * eff
		table[i] = tableRow{v, j, cp, rpm, angle, eff, power}
		v += p.StepSize
	}
	return table
}

func getCharts(p props, table []tableRow, cp, eff [][]float64) ([]interface{}, []interface{}) {
	J, Angle := make([]float64, len(table)), make([]float64, len(table))
	Cp, Eff := make([]float64, len(table)), make([]float64, len(table))
	color := make([]string, len(table))
	for i, v := range table {
		J[i], Angle[i] = v.J, v.Angle
		Cp[i], Eff[i] = v.Cp, v.Eff
		color[i] = "#0275d8"
	}
	colorscale := [][]interface{}{{0, "rgb(23, 55, 35)"}, {0.5, "rgb(37, 157, 81)"}, {1, "rgb(186, 228, 174)"}}
	opacity := 0.9
	hoverCp := "<b>J</b>: %{x}<br><b>Angle</b>: %{y}°<br><b>Cp</b>: %{z}<extra></extra>"
	hoverEff := "<b>J</b>: %{x}<br><b>Angle</b>: %{y}°<br><b>Eff</b>: %{z}<extra></extra>"
	show := false
	chartCp := []interface{}{
		surfaceChart{
			colorscale, hoverCp, opacity, show, "surface", cp[0], cp[len(cp) - 1], cp[1:len(cp) - 1],
		},
		scatterChart{
			marker{color, 6}, hoverCp, line{"#0275d8"}, opacity, "scatter3d", J, Angle, Cp,
		},
	}
	chartEff := []interface{}{
		surfaceChart{
			colorscale, hoverEff, opacity, show, "surface", eff[0], eff[len(eff) - 1], eff[1:len(eff) - 1],
		},
		scatterChart{
			marker{color, 6}, hoverEff, line{"#0275d8"}, opacity, "scatter3d", J, Angle, Eff,
		},
	}
	return chartCp, chartEff
}

func getProps (request events.APIGatewayProxyRequest) props {
	maxSpeed, _ := strconv.ParseFloat(request.QueryStringParameters["max_speed"], 64)
	stepSize, _ := strconv.ParseFloat(request.QueryStringParameters["step_size"], 64)
	propSpeed, _ := strconv.ParseFloat(request.QueryStringParameters["prop_speed"], 64)
	diameter, _ := strconv.ParseFloat(request.QueryStringParameters["diameter"], 64)
	blades, _ := strconv.Atoi(request.QueryStringParameters["blades"])
	Cp, _ := strconv.ParseFloat(request.QueryStringParameters["cp"], 64)
	power, _ := strconv.ParseFloat(request.QueryStringParameters["power"], 64)
	ratio, _ := strconv.ParseFloat(request.QueryStringParameters["ratio"], 64)
	return props{maxSpeed, stepSize, propSpeed, diameter, blades, Cp, power, ratio}
}

type charts struct {
	Cp []interface{} `json:"cp"`
	Eff []interface{} `json:"eff"`
}

func handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	p := getProps(request)
    cp, cpXZ, eff := loadData(p)
	table := getTable(p, cpXZ, eff)
	chartCp, chartEff := getCharts(p, table, cp, eff)
	body := struct{
		Table []tableRow `json:"table"`
		Charts charts `json:"charts"`
	}{
		table,
		charts{
			chartCp,
			chartEff,
		},
	}
	js, _ := json.Marshal(body)
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body:       string(js),
    }, nil
}

func main() {
        lambda.Start(handle)
}