package main

import (
	"sync"
	"testing"
)

func TestIsCombo(t *testing.T) {
	c := &curCmds{
		mu:   &sync.Mutex{},
		cmds: []string{},
	}

	if isCombo(c, SELECT) {
		t.Fatalf("Expected single SELECT invocation to return false")
	}

	c = &curCmds{
		mu:   &sync.Mutex{},
		cmds: []string{UP, DOWN},
	}

	if isCombo(c, DOWN) {
		t.Fatalf("Only SELECT in the last postion should trigger the command, DOWN triggered instead")
	}

	c = &curCmds{
		mu:   &sync.Mutex{},
		cmds: []string{UP, UP, DOWN, DOWN},
	}

	if !isCombo(c, SELECT) {
		t.Fatalf("Correct invocation should have returned true")
	}
}
