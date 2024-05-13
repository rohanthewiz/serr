## SErr package
The Structured Error package can enrich an error with attributes
at each stack frame  while the error is being bubbled up to the caller
without concern for logging at each frame. The SErr can then be logged
with a structured logger like `github.com/rohanthewiz/logger` or
printed with it's own string functions.

### Usage
(See the included tests for more examples)

```go
package main

import (
	"errors"

	"github.com/rohanthewiz/logger
	"github.com/rohanthewiz/serr
)

func main() {
	// Given an error
	err := errors.New("some error has occurred")

	// We can wrap the error with a message
	errWrapped := serr.Wrap(err, "Error occurred when trying to do things")

	// We can printout errors and attributes in a nice format
	fmt.Println(serr.StringFromErr(errWrapped))
	// ==> some error has occurred => location->logtest/main.go:16; function->main.main; msg->Error occurred when trying to do things

	// A structured error aware logger like github.com/rohanthewiz/logger can output all attributes
	logger.LogErr(errWrapped, "An Error occurred")
	// => ERRO[0000] some error has occurred	error="some error has occurred" fields.msg="An Error occurred - Error occurred when trying to do things"
	// function=main.main location="logtest/main.go:16"

	// We can wrap an error with some attributes
	err3 := serr.Wrap(err, "cats", "okay", "dogs", "hmm")

	logger.LogErr(err3, "Animals are cool")
	// => ERRO[0000] some error has occurred  cats=okay dogs=hmm error="some error has occurred" fields.msg="Animals are cool"
	// function=main.main location="logtest/main.go:27"
}
```
