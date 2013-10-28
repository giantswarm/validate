package validate

import (
	"testing"
)

func Test_ValidateLowAlphabet_1(t *testing.T) {
	s := "goiscool"
	if ValidateLowAlphabet([]byte(s)) != true {
		t.Error("Failed")
	} else {
		t.Log("Passed")
	}
}

func Test_ValidateLowAlphabet_2(t *testing.T) {
	s := "goiscoolandshit"
	if ValidateLowAlphabet([]byte(s)) != true {
		t.Error("Failed")
	} else {
		t.Log("Passed")
	}
}

func Test_ValidatePrintableRunes_1(t *testing.T) {
	s := "goiscool"
	if ValidatePrintableRunes([]byte(s)) != true {
		t.Error("Failed")
	} else {
		t.Log("Passed")
	}
}

func Test_ValidatePrintableRunes_2(t *testing.T) {
	s := "goiscoolandshit"
	if ValidatePrintableRunes([]byte(s)) != true {
		t.Error("Failed")
	} else {
		t.Log("Passed")
	}
}
