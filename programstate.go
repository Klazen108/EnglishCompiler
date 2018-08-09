package main

import (
	"fmt"
)

//The ProgramState encapsulates the complete state of the program.
//Identifiers identify memory in this structure
type ProgramState struct {
	identifiers map[string]string
	types       map[string]DataType
	elseFlag    bool
}

//A StateChangeError indicates that the program state was unable to update
//itself based on the requested parameters
type StateChangeError struct {
	id      Identifier
	message string
}

func (state ProgramState) FlagElse(set bool) {
	state.elseFlag = set
}

func (state ProgramState) IsFlaggedElse() bool {
	return state.elseFlag
}

func (e StateChangeError) Error() string {
	return fmt.Sprintf("Error changing program state, Identifier %s Type %s: %s", e.id.name, e.id.dType.toString(), e.message)
}

func (state ProgramState) SetValue(id Identifier, value string, dType DataType) error {
	if vType, ok := state.types[id.name]; ok {
		if vType != dType {
			return StateChangeError{id: id, message: "Invalid datatype, expected " + vType.toString()}
		}
	} else {
		state.types[id.name] = dType
	}
	state.identifiers[id.name] = value

	return nil
}

//Return: (type, exists?)
func (s ProgramState) getType(i Identifier) (DataType, bool) {
	t, ok := s.types[i.name]
	return t, ok
}

func (s ProgramState) getValue(i Identifier) (string, bool) {
	t, ok := s.identifiers[i.name]
	return t, ok
}

func (s ProgramState) PrintState() {
	for k, v := range s.identifiers {
		fmt.Printf("(%6s)%s = %s\n", s.types[k].toString(), k, v)
	}
}
