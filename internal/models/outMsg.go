package models

// Сообщение с результатами работы воркера над входящими данными
type OutMsg[In, Out any] struct {
	Id          int
	Result      Out
	IncomigData In
	Err         error
}
