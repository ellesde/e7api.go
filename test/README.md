# e7api.go tests

This directory contains additional test suites beyond the unit tests already in 
[../e7](../e7). Whereas the unit tests run very quickly since they don't make
any network calls and are run by Travis on every commit. The tests in this
directory are only run manually.

The test packages are:

integration
-----------
This will exercise the entire e7api.go library against the live EpicSevenDB API.
These tests will verify that the library is properly coded against the actual
behavior of the API and will, hopefully, fail upon any incompatible chagne in
the API.

Because these tests are running using live data, there is a much higher
probability of false positives in test failures due to network issues, test
data having been changed, etc.

Run tests using:

    go test -v -tags=integration ./integration