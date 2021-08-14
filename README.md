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
    // Given an error
    err := errors.New("Some error has occurred")
   
    // We can wrap the error with a message
    err2 := serr.Wrap(err, "Error occurred trying to do things")
    // A structured error aware logger can output all attributes
    rlog.LogErr(err2, "An Error occurred")
        // => ERROR[0000] "An Error occurred msg="Error occurred trying to do things" error="Some error has occurred"
    
    // We can wrap an error with some attributes  
    err3 := serr.Wrap(err, "cat", "aight", "dogs", "I dunno")
    rlog.LogErr(err3, "Animals are cool")
       // => ERRO[0000] Animals are cool cat=aight dogs="I dunno" error="Some error has occurred"
}
```
