package utils

// State represents the state of the branch predictor
type State uint8

// State constants, representing the state of the branch predictor
const (
	StronglyNotTaken State = iota // 00
	WeaklyNotTaken                // 01
	WeaklyTaken                   // 10
	StronglyTaken                 // 11
)

// StateMap maps the string representation of the state to the actual state
var StateMap = map[string]State{
	"StronglyNotTaken": StronglyNotTaken,
	"WeaklyNotTaken":   WeaklyNotTaken,
	"WeaklyTaken":      WeaklyTaken,
	"StronglyTaken":    StronglyTaken,
}
