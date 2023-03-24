package selector

import (
	"github.com/consensys/gnark/frontend"
)

var api frontend.API

func Configure(a frontend.API) {
	api = a
}
