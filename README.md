## Installation

`$ go get -u github.com/daniilty/go-exported-local-variables/...`

## Usage

`$ go-exported-local-variables -- path/to/file.go`

**Or**
`$ go-exported-local-variables ./...`
If you wish to run it on all your project's files

### Sample of a bad code

```golang
func Incorrect() string {
        A := "Oh it seems not so good" \\ sample.go:10:2: local variable A should not be exported

        return A
}
```