package service

import "errors"

var (
	NotFound            = errors.New("not found")
	IntervalServerError = errors.New("interval server error")
)
