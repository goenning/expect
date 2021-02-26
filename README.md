# github.com/goenning/expect

An opinionated and minimalist assert module for Go

📦 Zero Dependencies

## What it does

Writing assertions for unit test using Go standard can be quite cumbersome and repetitive. `expect` is a tiny library that simplifies and adds more clarify to assertions.

Before:

```go
if err != nil {
    t.Fatalf("Err is not nil")
}
if number != 4 {
    t.Fatalf("Number is not 4")
}
```

After:
```go
Expect(err).IsNil()
Expect(result).Equals(4)
```

## How to use

```go
import (
    . "github.com/goenning/expect"
)

func TestCanAddNumbers(t *testing.T) {
    RegisterT(t)

    result, err := DoSomething()

    Expect(err).IsNil()
    Expect(result).Equals(4)
}
```

All supported assertion operations:

- Equals(expected)
- NotEquals(other)
- ContainsString(substr)
- IsTrue()
- IsFalse()
- IsEmpty()
- IsNotEmpty()
- IsNotNil()
- IsNil()
- HasLen(length)
- Panics() 
- EventuallyEquals(expected)
- WithinTime(time, duration)