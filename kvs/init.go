package kvs

import logxi "github.com/aa-ar/logxi"

var logger logxi.Logger

func init() {
	logger = logxi.New("dat.cache")
}
