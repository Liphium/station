package magic_tests

import (
	"fmt"
	"testing"

	"github.com/Liphium/magic/mconfig"
)

// Do not call this function anything with Test, it will cause errors
func MagicDefault(t *testing.T, p *mconfig.Plan) {
	fmt.Println("Hello, I'm the greatest wizzard of all time!")
}
