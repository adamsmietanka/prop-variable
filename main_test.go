package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func Test_readCsvFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want [][]float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readCsvFile(tt.args.filePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readCsvFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertToFloat(t *testing.T) {
	type args struct {
		records [][]string
	}
	tests := []struct {
		name    string
		args    args
		want    [][]float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertToFloat(tt.args.records)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertToFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertToFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadData(t *testing.T) {
	type args struct {
		p props
	}
	tests := []struct {
		name  string
		args  args
		want  [][]float64
		want1 [][]float64
		want2 [][]float64
	}{
		{
			name: "Basic",
			args: args{
				props{
					MaxSpeed:  150,
					StepSize:  10,
					PropSpeed: 20,
					Diameter:  3.902,
					Blades:    3,
					Cp:        0.0902,
					Power:     800,
					Ratio:     0.4,
				},
			},
			want: [][]float64{},
			want1: [][]float64{},
			want2: [][]float64{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := loadData(tt.args.p)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadData() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("loadData() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("loadData() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_getProps(t *testing.T) {
	type args struct {
		request events.APIGatewayProxyRequest
	}
	tests := []struct {
		name string
		args args
		want props
	}{
		{
			name: "Basic",
			args: args{
				request: events.APIGatewayProxyRequest{
					QueryStringParameters: map[string]string{
						"max_speed":  "150",
						"step_size":  "10",
						"prop_speed": "20",
						"diameter":   "3.902",
						"blades":     "3",
						"cp":         "0.0902",
						"power":      "800",
						"ratio":      "0.4",
					},
				},
			},
			want: props{
				MaxSpeed:  150,
				StepSize:  10,
				PropSpeed: 20,
				Diameter:  3.902,
				Blades:    3,
				Cp:        0.0902,
				Power:     800,
				Ratio:     0.4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getProps(tt.args.request); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getProps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_table(t *testing.T) {
	p := props{
		MaxSpeed:  150,
		StepSize:  10,
		PropSpeed: 20,
		Diameter:  3.902,
		Blades:    3,
		Cp:        0.0902,
		Power:     800,
		Ratio:     0.4,
	}
    _, cpXZ, eff := loadData(p)
	type args struct {
		p    props
		cpXZ [][]float64
		eff  [][]float64
	}
	tests := []struct {
		name string
		args args
		want []tableRow
	}{
		{
			name: "Basic",
			args: args{p, cpXZ, eff},
			want: []tableRow{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTable(tt.args.p, tt.args.cpXZ, tt.args.eff); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("table() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handle(t *testing.T) {
	type args struct {
		ctx     context.Context
		request events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		{
			name: "Basic",
			args: args{
				request: events.APIGatewayProxyRequest{
					QueryStringParameters: map[string]string{
						"max_speed":  "150",
						"step_size":  "10",
						"prop_speed": "20",
						"diameter":   "3.902",
						"blades":     "3",
						"cp":         "0.0902",
						"power":      "800",
						"ratio":      "0.4",
					},
				},
			},
			want: events.APIGatewayProxyResponse{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handle(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handle() = %v, want %v", got, tt.want)
			}
		})
	}
}
