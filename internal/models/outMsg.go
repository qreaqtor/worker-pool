package models

type OutMsg[T any] struct {
	Id     int
	Result T
	Data   string
	Err    error
}
