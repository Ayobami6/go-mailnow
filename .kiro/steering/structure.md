---
inclusion: always
---

# Project Structure

## Directory Layout

```
.
├── mailnow.go          # Main package file with core types
├── go.mod              # Go module definition
├── types/              # Type definitions and interfaces
│   └── types.go
└── tests/              # Test files
    └── tests.go
```

## Conventions

### Package Organization

- Root package (`mailnow`) contains the main client struct and public API
- `types/` package for shared type definitions, interfaces, and data structures
- `tests/` package for test utilities and test suites

### Code Style

- Follow standard Go conventions and idioms
- Use pointer receivers for struct methods that modify state
- API keys and sensitive data should be handled as pointers for optional/nullable behavior
- Package names are lowercase, single-word

### Naming

- Main client struct: `Mailnow`
- Use descriptive, idiomatic Go names
- Exported types start with capital letters
- Private fields use lowercase

### Testing

- Test files go in the `tests/` package
- Use table-driven tests where appropriate
- Follow Go testing best practices
