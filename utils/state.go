package utils

type State uint8

const (
	StronglyNotTaken State = iota // 00
	WeaklyNotTaken                // 01
	WeaklyTaken                   // 10
	StronglyTaken                 // 11
)
