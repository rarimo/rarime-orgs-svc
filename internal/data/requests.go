package data

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

const MaxDigits = 6

func GenerateOTP() string {
	bi, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(10, float64(MaxDigits)))),
	)
	if err != nil {
		panic(errors.Wrap(err, "failed to generate uniform random value"))
	}
	return fmt.Sprintf("%0*d", MaxDigits, bi)
}
