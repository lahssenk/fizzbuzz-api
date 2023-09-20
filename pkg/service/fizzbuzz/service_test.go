package fizzbuzz

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestService_ComputeFizzBuzzRange(t *testing.T) {
	tests := []struct {
		name    string
		req     *ComputeFizzBuzzRangeParams
		resp    *ComputeFizzBuzzRangeOutput
		wantErr bool
	}{
		{
			name: "OK",
			req: &ComputeFizzBuzzRangeParams{
				String1: "fizz",
				String2: "buzz",
				Int1:    3,
				Int2:    5,
				Limit:   16,
			},
			resp: &ComputeFizzBuzzRangeOutput{
				Data: []string{
					"1", "2", "fizz", "4", "buzz", "fizz", "7", "8", "fizz", "buzz", "11", "fizz", "13", "14", "fizzbuzz", "16",
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid input",
			req:     &ComputeFizzBuzzRangeParams{},
			resp:    nil,
			wantErr: true,
		},
	}

	svc := NewService()

	for ind := range tests {
		test := tests[ind]
		t.Run(test.name, func(t *testing.T) {
			resp, err := svc.ComputeFizzBuzzRange(test.req)
			if (err != nil) != test.wantErr {
				t.Errorf("unexpected err: %v", err)
				return
			}

			if err != nil {
				return
			}

			require.Equal(t, resp.Data, test.resp.Data)
		})
	}
}

func Test_validateComputeFizzBuzzRangeParams(t *testing.T) {
	type args struct {
		req *ComputeFizzBuzzRangeParams
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				req: &ComputeFizzBuzzRangeParams{
					String1: "fizz",
					String2: "buzz",
					Int1:    3,
					Int2:    5,
					Limit:   100,
				},
			},
			wantErr: false,
		},
		{
			name: "string1 required",
			args: args{
				req: &ComputeFizzBuzzRangeParams{
					String1: "",
					String2: "buzz",
					Int1:    3,
					Int2:    5,
					Limit:   100,
				},
			},
			wantErr: true,
		},
		{
			name: "string2 required",
			args: args{
				req: &ComputeFizzBuzzRangeParams{
					String1: "fizz",
					String2: "",
					Int1:    3,
					Int2:    5,
					Limit:   100,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid int1",
			args: args{
				req: &ComputeFizzBuzzRangeParams{
					String1: "fizz",
					String2: "buzz",
					Int1:    0,
					Int2:    5,
					Limit:   100,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid int2",
			args: args{
				req: &ComputeFizzBuzzRangeParams{
					String1: "fizz",
					String2: "buzz",
					Int1:    3,
					Int2:    0,
					Limit:   100,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid limit",
			args: args{
				req: &ComputeFizzBuzzRangeParams{
					String1: "fizz",
					String2: "buzz",
					Int1:    3,
					Int2:    5,
					Limit:   0,
				},
			},
			wantErr: true,
		},
	}

	for ind := range tests {
		test := tests[ind]

		t.Run(test.name, func(t *testing.T) {
			if err := validateComputeFizzBuzzRangeParams(test.args.req); (err != nil) != test.wantErr {
				t.Errorf(
					"validateComputeFizzBuzzRangeParams() error = %v, wantErr %v",
					err,
					test.wantErr,
				)
			}
		})
	}
}

func Test_mapper_computeFizzBuzzRange(t *testing.T) {
	type fields struct {
		fizz     string
		buzz     string
		fizzbuzz string
		fizzmod  int
		buzzmod  int
		limit    int
	}

	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "fizzbuzz-3-5-16",
			fields: fields{
				fizz:     "fizz",
				buzz:     "buzz",
				fizzbuzz: "fizzbuzz",
				fizzmod:  3,
				buzzmod:  5,
				limit:    16,
			},
			want: []string{
				"1", "2", "fizz", "4", "buzz", "fizz", "7", "8", "fizz", "buzz", "11", "fizz", "13", "14", "fizzbuzz", "16",
			},
		},
	}

	for ind := range tests {
		test := tests[ind]
		t.Run(test.name, func(t *testing.T) {
			m := mapper{
				fizz:     test.fields.fizz,
				buzz:     test.fields.buzz,
				fizzbuzz: test.fields.fizzbuzz,
				fizzmod:  test.fields.fizzmod,
				buzzmod:  test.fields.buzzmod,
				limit:    test.fields.limit,
			}
			if got := m.computeFizzBuzzRange(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("mapper.computeFizzBuzzRange() = %v, want %v", got, test.want)
			}
		})
	}
}
