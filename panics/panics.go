package panics

import (
	"fmt"

	"github.com/gomponents/gontainer-helpers/caller"
)

func WrapProvider(msg string, provider interface{}, params ...interface{}) interface{} {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		panic(fmt.Sprintf("%s: %s", msg, r))
	}()
	return caller.MustCallProvider(provider, params...)
}
