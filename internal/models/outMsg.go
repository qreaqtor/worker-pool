package models

type OutMsg[In, Out any] struct {
	Id     int
	Result Out
	Data   In
	Err    error
}
