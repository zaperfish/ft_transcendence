package defaults

import (
    // Std
	"runtime"

    // External
	"github.com/alexedwards/argon2id"
)

var ArgonParams = &argon2id.Params{
	Memory:      128 * 1024,
	Iterations:  4,
	Parallelism: uint8(runtime.NumCPU()),
	SaltLength:  16,
	KeyLength:   32,
}
