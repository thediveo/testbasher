# Test BASHer

[![PkgGoDev](https://pkg.go.dev/badge/github.com/thediveo/testbasher)](https://pkg.go.dev/github.com/thediveo/testbasher)
[![GitHub](https://img.shields.io/github/license/thediveo/testbasher)](https://img.shields.io/github/license/thediveo/testbasher)
![build and test](https://github.com/thediveo/testbasher/workflows/build%20and%20test/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/thediveo/testbasher)](https://goreportcard.com/report/github.com/thediveo/testbasher)
![Coverage](https://img.shields.io/badge/Coverage-91.2%25-brightgreen)

"Test BASHer" is a painfully simple bash script management and execution for
simple unit test script harnesses. It is intended for such cases where only
one or few short and simple test scripts are required per specific test and
it's best to keep them near the test case itself to keep scripts and test in
sync.

A test starts a (BASHer) script and then interacts with it: the script may
report information about things it set up dynamically, such as the IDs of
Linux kernel namespaces, et cetera, which the test needs to read in order to
test dynamic assumptions. And the script in turn may wait for the test to step
it through multiple phases in order to complete a specific test. An example is
[lxkns](https://github.com/thediveo/lxkns) where transient namespaces get
created, but the processes keeping them alive must not terminate before the
test has reached a certain phase in its course.

Please refer to the [![PkgGoDev](https://pkg.go.dev/badge/github.com/thediveo/testbasher)](https://pkg.go.dev/github.com/thediveo/testbasher) for details.

## Usage

The basic usage pattern is as follows:

- create a `b := Basher{}`, and don't forget to `defer b.Done()`.
- if required, add common BASH code using `b.Common("...script code...")` to
  be reused in your scripts.
- add one or more BASH scripts using `b.Script("name", "...script code...")`.
- start your entry point script with `c := b.Start("name")`, and `defer
  c.Close()`.
- read data output from your script: `c.Decode(&data)`.
  - Golang 1.14 and later: in case the expected data cannot be decoded,
    `c.Decode` panics with details including the exact (JSON) data read from
    the script which could not be decoded. No more stupid JSON "syntax errors
    at offset 666", but instead you'll see the JSON data read up to the point
    where things went south.
- in case of multiple phases, step forward by calling `c.Proceed()`.

And now for some code to further illustrate the above usage pattern list:

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

`testbasher` is Copyright 2020-23 Harald Albrecht, and licensed under the Apache
License, Version 2.0.
