package fizzbuzz

import (
	"context"
	"fmt"
	"strconv"

	fizzbuzz_v1 "github.com/lahssenk/fizzbuzz-api/protogen/fizzbuzz/v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ fizzbuzz_v1.FizzBuzzServiceServer = (*Service)(nil)

type Service struct {
	fizzbuzz_v1.UnimplementedFizzBuzzServiceServer
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ComputeFizzBuzzRange(ctx context.Context, req *fizzbuzz_v1.ComputeFizzBuzzRangeRequest) (*fizzbuzz_v1.ComputeFizzBuzzRangeResponse, error) {
	if err := validateComputeFizzBuzzRangeRequest(req); err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	m := newMapperFromRequest(req)
	res := m.computeFizzBuzzRange()

	resp := fizzbuzz_v1.ComputeFizzBuzzRangeResponse{
		Output: res,
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

func validateComputeFizzBuzzRangeRequest(req *fizzbuzz_v1.ComputeFizzBuzzRangeRequest) error {
	switch {
	case req.String1 == "":
		return errors.Errorf(ErrString1Required)
	case req.String2 == "":
		return errors.Errorf(ErrString2Required)
	case req.Int1 < 1 || req.Int1 > 100:
		return errors.Errorf(ErrInvalidInt1)
	case req.Int2 < 1 || req.Int2 > 100:
		return errors.Errorf(ErrInvalidInt2)
	case req.Limit < 1 || req.Limit > 100:
		return errors.Errorf(ErrInvalidLimit)
	default:
		return nil
	}
}

func newMapperFromRequest(req *fizzbuzz_v1.ComputeFizzBuzzRangeRequest) *mapper {
	return &mapper{
		fizz:     req.String1,
		buzz:     req.String2,
		fizzbuzz: fmt.Sprintf("%s%s", req.String1, req.String2),
		fizzmod:  int(req.Int1),
		buzzmod:  int(req.Int2),
		limit:    int(req.Limit),
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
