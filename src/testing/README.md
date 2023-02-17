# Testing

Testing utilities

## Test comments
The comment of a test should explain:
* The test's **subject under test** (or `SUT`). A SUT can be a particular
  function, or more broadly a particular logical functionality involving
  multiple APIs.
* The test's expected results from the SUT. Describe the expectation in terms of
  logical results. Use examples to help readers understand.

## Test APIs

Use `testify/assert` and `testify/require`.

* `require.*` when the check failed, instantly stop the test, and the rest of
  test's code is skipped.
* `assert.*` when failed, only record error, and the rest of test's code is
  still executed.

Choose require and assert accordingly. The sample code looks like

```golang
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestXXX(t *testing.T) {
  assert := assert.New(t)
  require := require.New(t)
  ...
  assert.Nil(err)  // record failure and continues
  require.Nil(err) // instancely stops test, and rest of code is skipped
}
```
