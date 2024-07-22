package trading

import "math/big"

type Event struct {
    Amount0Out *big.Int
    Amount1Out *big.Int
    To         string
}
