# Style Guide

This doc describes rules that are not enforceable in linter, and require
conscious effort to uphold during writing and reviewing code.

## General

[Golang Best Practices](
https://google.github.io/styleguide/go/best-practices.html)

## Type and method naming

Type names are noun, methods are verb. Type describes concepts, so they should
be naming something, thus nouns. Method performs actions, so they should do
something, thus verbs.

Perfer
```
type Deployment struct { ... }
func (d *Deployment) Deploy() { .... }
```
Over
```
type Deploy struct { ... }
func (d *Deploy) Worker()
```

## Error message

Always provide contexts for error message. Use
`github.com/tricorder/src/utils/errors.Wrap` to wrap error with contextual
information.

```
import "github.com/tricorder/src/utils/errors"

func DoA() {
  err := DoB()
  if err != nil {
    return errors.Wrap("doing A", "do B", err)
  }
}
```

### Rationale
This consistent style gives the context of an error message. And provide clues
of understanding the error in one place.

## README.md

The title should be named after the parent directory
```
cat src/agent/README.md
# Agent
Agent is doing someting incredible ...
```

### Rationale

Keeps things consistent

