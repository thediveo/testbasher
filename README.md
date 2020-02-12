# Test BASHer

[![GoDoc](https://godoc.org/github.com/TheDiveO/testbasher?status.svg)](http://godoc.org/github.com/TheDiveO/testbasher)
[![GitHub](https://img.shields.io/github/license/thediveo/testbasher)](https://img.shields.io/github/license/thediveo/testbasher)
![build and test](https://github.com/TheDiveO/testbasher/workflows/build%20and%20test/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/thediveo/testbasher)](https://goreportcard.com/report/github.com/thediveo/testbasher)

"Test BASHer" is a painfully simple bash script management and execution for
simple unit test script harnesses. It is intended for such cases where only
one or few short and simple test scripts are required per specific test and
it's best to keep them near the test case itself to keep scripts and test in
sync.

A test starts a script and then interacts with it: the script may report
information about things it set up dynamically, such as the IDs of Linux
kernel namespaces, et cetera. And the script may wait for the test to step it
through multiple phases (or steps) in order to complete the specific test.

An example is [lxkns](https://github.com/TheDiveO/lxkns) where transient
namespaces get created, but the processes keeping them alive must not
terminate before the test has reached a certain phase in its course.

Please refer to the [![GoDoc](https://godoc.org/github.com/TheDiveO/testbasher?status.svg)](http://godoc.org/github.com/TheDiveO/testbasher) for details.

## Usage

```go

func example() {
    scripts := Basher{}
    defer scripts.Done()
    // Define a first script named "newuserns", which we will later start as
    // out entry point. This script creates a new Linux kernel user namespace,
    // where the unshare(2) command then executes a second script (which we'll
    // define next).
    //
    // Since the defined script are stored in temporary files, there is no way
    // to know beforehand where these files will be stored and named, and thus
    // we don't know how to directly call them. Instead, we use the $userns
    // environment variable which will point to the correct temporary
    // filepath.
    scripts.Script("newuserns", `
unshare -Ufr $userns # call the script named "userns" via $userns.
`)
    // This second script named "userns" returns information about the newly
    // created and then waits to proceed into termination, thereby destroying
    // the user namespace.
    scripts.Script("userns", `
echo "\"$(readlink /proc/self/ns/user)\"" # ...turns it into a JSON string.
read # wait for test to proceed()
`)
    // Start the first defined script named "newuserns"...
    cmd := scripts.Start("newuserns")
    defer cmd.Close()
    // Read in the (JSON) information sent by the running script.
    var userns string
    cmd.Decode(&userns)
    fmt.Println("temporary user namespace:", userns)
    // Tell the script to finish; this can be omitted if the last step,
    // because calling Close() on a Basher will automatically issue a final
    // Proceed().
    cmd.Proceed()
}
```

## Copyright and License

`testbasher` is Copyright 2020 Harald Albrecht, and licensed under the Apache
License, Version 2.0.
