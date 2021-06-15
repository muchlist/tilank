package enum

// State - Custom type to hold value for violation state
type State int

// Declare related constants for each direction starting with index -1
const (
	StUndefined   State = iota - 1 // EnumIndex = -1
	StDraft                        // EnumIndex = 0
	StNeedApprove                  // EnumIndex = 1
	StApproved                     // EnumIndex = 2
	StEmailSend                    // EnumIndex = 3
)

// String - Creating common behavior - give the type a String function
func (s State) String() string {
	return [...]string{"Undefined", "Draft", "NeedApprove", "Approved", "EmailSend"}[s+1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex function
func (s State) EnumIndex() int {
	return int(s)
}

// IntToState - return Undefined if not valid
func IntToState(value int) State {
	if value < StUndefined.EnumIndex() || value > StEmailSend.EnumIndex() {
		return StUndefined
	}
	return State(value)
}
