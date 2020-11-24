package client

import (
	"github.com/jcmturner/gokrb5/v8/credentials"
)

func LoadCCache(cpath string) (*credentials.CCache, error) {
	return credentials.LoadCCache(cpath)
}
