// Copyright 2020 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy
// of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package testbasher

import "fmt"

// Defines two scripts, where the first script will be the entry point, and
// the second script is called (indirectly) from the first script and run
// within a new Linux kernel user namespace. The second script will send
// information about the newly created user namespace to our example, which
// might be a test case then checking the results. Finally, it tells the
// auxiliary script to enter the next phase, which is simply terminating,
// thereby automatically destroying the newly created user namespace.
//
// Please note that some distributions mistakenly see the Linux kernel ability
// to create new user namespaces without special capabilities to be insecure
// (but cannot demonstrate how), so they patch upstream kernels to block this
// functionality. On such distributions this example will fail unless run as
// root ... "thank you for observing the security precautions".
func ExampleBasher() {
	// A zero Basher is already usable, so need for NewBasher().
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
