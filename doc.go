/*
Package testbasher provides painfully simple BASH script management and
execution for small unit test script harnesses. It is intended for such cases
where only one or few short and simple auxiliary scripts are required per
specific test and it's best to keep them near the test case itself to keep
scripts and test in sync.

This package defines only these two elements: Basher and TestCommand.

# Basher

Basher handles auxiliary test harness scripts as integral parts of your unit
tests. Running a Basher script (or often set of scripts) is done via
TestCommand.

# TestCommand

TestCommand simplifies handling and interaction with (test) commands and
scripts. It has a simplified reporting and interaction interface tailored
towards test harness scripts. As its name already suggests, TestCommand is
good for use in some types of tests, but it is not a general purpose tool.
*/
package testbasher
