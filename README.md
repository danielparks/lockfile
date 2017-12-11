# Manage POSIX lock files with Go

This provides a simple way to create a lock that will be automatically released
when the process dies, even if it is killed with SIGKILL.
