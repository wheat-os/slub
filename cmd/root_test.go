package cmd

import (
	"fmt"
	"regexp"
	"testing"
)

func TestCmd(t *testing.T) {
	fmt.Println(regexp.MatchString(`^([a-z]|[A-Z])\w+$`, "a1"))
}
