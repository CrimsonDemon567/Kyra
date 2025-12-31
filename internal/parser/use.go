package parser

import (
	"fmt"
)

// ---------------------------
// Utility: pretty-print module path
// ---------------------------

func (u *UseStmt) String() string {
	p := ""
	for i, s := range u.Path {
		if i > 0 {
			p += "/"
		}
		p += s
	}
	if u.IsStdlib {
		return fmt.Sprintf("use sdt/%s", p)
	}
	return fmt.Sprintf("use %s", p)
}
