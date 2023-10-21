package json

import "github.com/colinrs/pkgx/kq/core"

type extractor struct {
}

func NewExtractor() core.Extractor {
	return nil
}

func (e *extractor) Unmarshal(*core.InputMessage) (interface{}, error) {
	return nil, nil
}