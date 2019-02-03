# Manage POSIX lock files with Golang

This provides a simple way to create a lock that will be automatically released
when the process dies, even if it is killed with SIGKILL.

**Currently, you cannot release the lock without exiting the process.**

## Usage

``` go
import "github.com/danielparks/lockfile"

func main() {
	lockfile.ObtainLock("/path/to/lockfile")
	// Do stuff.
	// The lock will be held until the process ends.
}
```

## Authorship and license

Written by [Daniel Parks](https://demon.horse). Licensed under the [MIT
License](LICENSE), so you can modify and distribute this code as long as you
include a copy of the [license](LICENSE).
