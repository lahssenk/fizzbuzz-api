package fizzbuzz

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

// Service handling the fizzbuzz feature
type Service struct{}

// NewService constructs a fizzbuzz Service
func NewService() *Service {
	return &Service{}
}

// ComputeFizzBuzzRangeParams contains the input data defining the fizzbuzz task
type ComputeFizzBuzzRangeParams struct {
	// custom value for fizz case
	String1 string
	// custom value for buzz case
	String2 string
	// modulo that will match fizz case
	Int1 int
	// modulo that wil match buzz case
	Int2 int
	// upper bound of the range on which to compute fizzbuzz
	Limit int
}

// ComputeFizzBuzzRangeOutput contains the output as a list of strings where each
// item is the result of the confugured fizzbuzz applied to the matching int value
// in the range from 1 to <limit>
type ComputeFizzBuzzRangeOutput struct {
	Data []string `json:"data"`
}

// ComputeFizzBuzzRange applies fizzbuzz confugured with the params to a range from 1 to <limit>
func (s *Service) ComputeFizzBuzzRange(
	params *ComputeFizzBuzzRangeParams,
) (*ComputeFizzBuzzRangeOutput, error) {
	if err := validateComputeFizzBuzzRangeParams(params); err != nil {
		return nil, err
	}

	m := newMapperFromParams(params)
	res := m.computeFizzBuzzRange()

	resp := ComputeFizzBuzzRangeOutput{
		Data: res,
	}

	return &resp, nil
}

func (m mapper) computeFizzBuzzRange() []string {
	res := make([]string, m.limit)

	for i := range res {
		res[i] = m.Map(i + 1)
	}

	return res
}

func validateComputeFizzBuzzRangeParams(params *ComputeFizzBuzzRangeParams) error {
	switch {
	case params.String1 == "":
		return errors.Errorf(ErrString1Required)
	case params.String2 == "":
		return errors.Errorf(ErrString2Required)
	case params.Int1 < 1 || params.Int1 > 100:
		return errors.Errorf(ErrInvalidInt1)
	case params.Int2 < 1 || params.Int2 > 100:
		return errors.Errorf(ErrInvalidInt2)
	case params.Limit < 1 || params.Limit > 100:
		return errors.Errorf(ErrInvalidLimit)
	default:
		return nil
	}
}

func newMapperFromParams(params *ComputeFizzBuzzRangeParams) *mapper {
	return &mapper{
		fizz:     params.String1,
		buzz:     params.String2,
		fizzbuzz: fmt.Sprintf("%s%s", params.String1, params.String2),
		fizzmod:  int(params.Int1),
		buzzmod:  int(params.Int2),
		limit:    int(params.Limit),
	}
}

type mapper struct {
	fizz, buzz, fizzbuzz    string
	fizzmod, buzzmod, limit int
}

func (m mapper) Map(val int) string {
	res1, res2 := val%m.fizzmod, val%m.buzzmod

	switch {
	case res1 == 0 && res2 == 0:
		return m.fizzbuzz
	case res1 == 0:
		return m.fizz
	case res2 == 0:
		return m.buzz
	default:
		return strconv.Itoa(val)
	}
}
