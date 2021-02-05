package tools

import (
	"strconv"
	"testing"
)

func TestCheckErr(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	_, err := strconv.Atoi("x1")
	CheckErr(err)
}
