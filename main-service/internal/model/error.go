package model

import "errors"

var ErrObjectNotFound = errors.New("object not found")
var ErrObjectAlreadyExists = errors.New("object already exists")
var ErrInvalidInput = errors.New("invalid input")
