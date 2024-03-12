## SErr package
The Structured Error package (serr) wraps errors for use with a structured logger like `github.com/rohanthewiz/rlog`.
This allows errors to be conveniently decorated with attributes and bubbled up from deep within a library without concern for actual logging within the library

### Usage
(See the included tests for more examples)

```go
package main
import (
    "github.com/rohanthewiz/rlog
    "github.com/rohanthewiz/serr
)

func ExerciseLogging() {
    // Given an error - this is the root error
    err := errors.New("Some error has occurred")
   
    // We can wrap the error with a message
    err2 := serr.Wrap(err, "Error occurred when trying to do things")
	
	// We can printout errors and attributes in a nice format
	serr.StringFromErr(err2)
	// ==> Some error has occurred => msg->Error occurred when trying to do things
	
    // A structured error aware logger like github.com/rohanthewiz/rlog can output all attributes
    rlog.LogErr(err2, "An Error occurred")
        // => ERROR[0000] "An Error occurred msg="Error occurred trying to do things" error="Some error has occurred"
    
    // We can wrap an error with some attributes  
    err3 := serr.Wrap(err, "cats", "okay", "dogs", "I dunno")
    rlog.LogErr(err3, "Animals are cool")
       // => ERRO[0000] Animals are cool cats=okay dogs="I dunno" error="Some error has occurred"
}
```
