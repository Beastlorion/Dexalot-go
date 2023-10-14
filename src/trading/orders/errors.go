package orders

import "errors"

var NoOrderError = errors.New("nil order")
var NoExchangeError = errors.New("nil exchange metadata")
var NoTimeInForceError = errors.New("nil time in force")
var NonPositivePriceError = errors.New("non-positive order price")
var NonPositiveQtyError = errors.New("non-positive order qty")
var NoSideError = errors.New("nil order side")
