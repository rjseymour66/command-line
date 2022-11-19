# Notes

## Go commands

Format code:
```go
$ go -w <file>.go
```
Build the binary:
```go
$ go build // uses module name for binary name
$ go build -o <binary-name>
```
Run tests:
Build the binary:
```go
$ go test -v
```

### Cross-compilation

Build static go binaries for operating systems that are different than the one that you are building it on. Because you build a static binary, the target machine does not need any additional libraries or tools to run the binary.

For example, use the `GOOS` environment variable with the `build` command to compile for a Windows machine:

```go
$ GOOS=window go build
// For a list of accepted GOOS values, see https://go.dev/src/go/build/syslist.go
```

## Strings

Initialize a buffer with string contents using the bytes.NewBufferString("string") func. This simulates an input (like STDIN):
```go
b := bytes.NewBufferString("string")
```

## Interfaces

```go
io.Reader // any go type that you can read data from
io.Writer // any go type that you can write to
```



### go.mod and go.sum

Go modules group related packages into a single unit to be versioned together. Because they track an application's dependencies, they ensure that users build the application with the same dependencies as the original developer.

Go modules allow you to write go programs outside of the $GOPATH directory, as in previous releases. 

Go sum records the checksum for each module in the application to ensure that each build uses the correct version.

## App structure

## Reading data

### Scanner for lines and words

The Scanner type accepts an `io.Reader` and reads data that is delimited by spaces or new lines. By default, it reads lines, but you can configure it to read words:

```go
scanner := bufio.NewScanner(r)
// scan words
scanner.Split(bufio.ScanWords)
```
Use the `.Scan()` function in a loop to read tokens:
```go
for scanner.Scan() {
    // do something 
}
```

To find the number of bytes in each scanned token:
```go
// scan words
scanner.Split(bufio.ScanWords)

byteLength := 0
for scanner.Scan() {
    byteLength += len(scanner.Bytes())    
}
```
## Flags

`flag.<FunctionName>` lets you define CLI flags. For example, to create a flag that performs an action if it exists, you can use `flag.Bool`.

The following flag function definition returns the value of a `bool` variable that stores the value of the flag. After you create a flag, you have to call the `Parse()` function to parse the arguments provided to the command line:

```go
lines := flag.Bool("l", false, "Count the number of lines")
      // flag.Bool(flagName, default value, usage info)
flag.Parse()
```
Now, you have a variable `lines` that stores the address of a `bool` set to `false`. To use the value in this variable that 'points' to an address, you have to derefence it with the `*` symbol. If you don't dereference, you will use the address of the variable, not the value stored at the address:

```go
func cliFunc(r io.Reader, useLines bool) {}
cliFunc(os.Stdin, *lines)
```


# Tests

