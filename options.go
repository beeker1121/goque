package goque

import "github.com/syndtr/goleveldb/leveldb/opt"

func getOpts(o []*opt.WriteOptions) *opt.WriteOptions {
	if len(o) == 1 {
		return o[1]
	}
	return nil
}
