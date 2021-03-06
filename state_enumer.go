// Code generated by "enumer -type=State -json"; DO NOT EDIT

package main

import (
	"encoding/json"
	"fmt"
)

const _StateName = "NEWENDORSEDFINANCED"

var _StateIndex = [...]uint8{0, 3, 11, 19}

func (i State) String() string {
	i -= 1
	if i < 0 || i >= State(len(_StateIndex)-1) {
		return fmt.Sprintf("State(%d)", i+1)
	}
	return _StateName[_StateIndex[i]:_StateIndex[i+1]]
}

var _StateValues = []State{1, 2, 3}

var _StateNameToValueMap = map[string]State{
	_StateName[0:3]:   1,
	_StateName[3:11]:  2,
	_StateName[11:19]: 3}

// StateString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func StateString(s string) (State, error) {
	if val, ok := _StateNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to State values", s)
}

// StateValues returns all values of the enum
func StateValues() []State {
	return _StateValues
}

// IsAState returns "true" if the value is listed in the enum definition. "false" otherwise
func (i State) IsAState() bool {
	for _, v := range _StateValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for State
func (i State) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for State
func (i *State) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("State should be a string, got %s", data)
	}

	var err error
	*i, err = StateString(s)
	return err
}
