package utils

type State uint8

const (
	StronglyNotTaken State = iota // 00
	WeaklyNotTaken                // 01
	WeaklyTaken                   // 10
	StronglyTaken                 // 11
)

var StateMap = map[string]State{
	"StronglyNotTaken": StronglyNotTaken,
	"WeaklyNotTaken":   WeaklyNotTaken,
	"WeaklyTaken":      WeaklyTaken,
	"StronglyTaken":    StronglyTaken,
}
