# Copier

[![Go Report Card](https://goreportcard.com/badge/github.com/knbr13/copier)](https://goreportcard.com/report/github.com/knbr13/copier)
[![Go Reference](https://pkg.go.dev/badge/github.com/knbr13/copier.svg)](https://pkg.go.dev/github.com/knbr13/copier)

A lightweight Go library for copying struct values with **shallow** and **deep copy** support. Handle pointers, nested structs, slices, maps, and common data types safely and efficiently.

## Features

- 🛠 **Shallow & Deep Copy**
  Choose between `ShallowCopyStruct` (shared references) or `DeepCopyStruct` (fully independent copies)
- 🧩 **Type-Safe**
  Automatically skips non-assignable fields and invalid types
- 🔍 **Nested Structure Support**
  Recursively copies structs, pointers, slices, and maps
- ⚡ **Efficient**
  Optimized field matching with precomputed maps
- 🛡️ **Error Handling**
  Clear error messages for invalid inputs

## Installation

```bash
go get github.com/knbr13/copier
```

## Quick Start

```go
package main

import (
	"fmt"
	"github.com/knbr13/copier"
)

type User struct {
	Name string
	Age  int
}

func main() {
	src := User{Name: "Alice", Age: 30}
	var dst User

	// Shallow copy
	if err := copier.ShallowCopyStruct(&dst, &src); err != nil {
		panic(err)
	}

	// Deep copy (same syntax)
	// copier.DeepCopyStruct(&dst, &src)

	fmt.Printf("Copied: %+v", dst) // Output: {Name:Alice Age:30}
}
```

## Usage

### Functions

```go
// Shallow copy - shares pointers/slices/maps
func ShallowCopyStruct(dst, src interface{}) error

// Deep copy - creates independent copies
func DeepCopyStruct(dst, src interface{}) error
```

### Basic Rules

1. **Destination** must be a **pointer to a struct**
2. **Source** can be a struct or pointer to struct
3. Field matching is **case-sensitive** and requires **exact name matches**
4. Unexported fields are **not copied**
5. Interface fields are **shallow copied**

## Examples

### 1. Pointer Handling
```go
num := 42
src := &struct{ Ptr *int }{Ptr: &num}
var dst struct{ Ptr *int }

copier.ShallowCopyStruct(&dst, src) // dst.Ptr == src.Ptr (shared)
copier.DeepCopyStruct(&dst, src)    // dst.Ptr != src.Ptr (new copy)
```

### 2. Nested Structures
```go
type Address struct { City string }
type User struct {
	Name    string
	Address *Address
}

src := &User{Name: "Bob", Address: &Address{City: "Paris"}}
var dst User

copier.DeepCopyStruct(&dst, src)
src.Address.City = "London" // dst.Address.City remains "Paris"
```

### 3. Slice/Map Handling
```go
src := &struct {
	IDs  []int
	Data map[string]bool
}{
	IDs:  []int{1, 2, 3},
	Data: map[string]bool{"admin": true},
}

var dst struct {
	IDs  []int
	Data map[string]bool
}

copier.DeepCopyStruct(&dst, src)
src.IDs[0] = 99           // dst.IDs remains [1,2,3]
src.Data["admin"] = false // dst.Data["admin"] remains true
```

## Error Handling

Common errors include:
```go
// Nil inputs
copier.ShallowCopyStruct(nil, &src) // Error: "destination cannot be nil"

// Type mismatches
src := struct{ A int }{A: 1}
dst := struct{ A string }{}
copier.DeepCopyStruct(&dst, &src) // Skipped: type mismatch

// Invalid destination
copier.DeepCopyStruct("not-a-pointer", &src) // Error: "pointer to a struct"
```

## Contributing

Contributions welcome!
1. Fork the repository
2. Create a feature branch
3. Submit a PR with tests

Report bugs via [GitHub Issues](https://github.com/knbr13/copier/issues).

## License

MIT License - See [LICENSE](https://github.com/knbr13/copier/blob/main/LICENSE) for details.
