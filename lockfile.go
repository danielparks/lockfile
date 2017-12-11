package lockfile

import (
	"bufio"
	"golang.org/x/sys/unix"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func ObtainLock(path string) {
	err := os.MkdirAll(filepath.Dir(path), os.FileMode(0755))
	if err != nil {
		log.Panic(err)
	}

	// This has to be 0600 so that a non-priviledged user can't deny access.
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.FileMode(0600))
	if err != nil {
		log.Panic(err)
	}

	// unix.Flock is not supported on Solaris. This call does not block.
	lock := unix.Flock_t{
		Type:   unix.F_WRLCK,
		Whence: 0, // There's no constant for this
		Start:  0,
		Len:    0, // to end of file
	}

	err = unix.FcntlFlock(uintptr(file.Fd()), unix.F_SETLK, &lock)
	if err == unix.EAGAIN {
		// Another process has the lock
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			log.Fatal("Another process is already running. Could not read PID from lock file: " + err.Error())
		}
		log.Fatal("Another process is already running with pid " + scanner.Text())
	} else if err != nil {
		log.Panic(err)
	}

	_, err = file.Seek(io.SeekStart, 0)
	if err != nil {
		log.Panic(err)
	}

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(strconv.Itoa(os.Getpid()) + `

In order to prevent a race condition, this file is not deleted when the process
ends. This uses fcntl (POSIX) locking in order to ensure that the lock is
cleanly released even if the process is terminated.
`)
	if err != nil {
		log.Panic(err)
	}

	err = writer.Flush()
	if err != nil {
		log.Panic(err)
	}
}
