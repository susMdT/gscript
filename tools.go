//go:build tools
// +build tools

package tools

import (
	_ "github.com/susMdT/gscript/stdlib/crypto"
	_ "github.com/susMdT/gscript/stdlib/encoding"
	_ "github.com/susMdT/gscript/stdlib/exec"
	_ "github.com/susMdT/gscript/stdlib/file"
	_ "github.com/susMdT/gscript/stdlib/net"
	_ "github.com/susMdT/gscript/stdlib/os"
	_ "github.com/susMdT/gscript/stdlib/rand"
	_ "github.com/susMdT/gscript/stdlib/requests"
	_ "github.com/susMdT/gscript/stdlib/time"
)
