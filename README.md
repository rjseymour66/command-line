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

## Methods

### Value recievers

Use a value receiver when the method:
- mutates the receiver
- is too large to reasonably pass in memory

Inside the method body, dereference the receiver using the `*` operator to access and mutate the value of the receiver. Otherwise, you are operating on the address value:

```go
func (r *Receiver) MethodName(param type) {
    *r = // do something else
}
```
> **Best practice**: The method set of a single type should use the same receiver type. If the method does not mutate the receiver, you can assign the pointer receiver to a value at the start of the method.


## Errors

`fmt.Errorf` creates a custom formatted error:
```go
return fmt.Errorf("Error: %s is not a valid string", s)
```

Test if the action returns a specific error. For example, the following snippet returns `nil` if the file does not exist; otherwise, it returns the error:
```go
file, err := os.ReadFile(filename)
if err != nil {
    if errors.Is(err, os.ErrNotExist) {
        return nil
    }
    return err
}
```

## Data structures and formats

### Slices

Add to a slice with append:
```go
*sliceName = append(*sliceName, valToAppend)
```
### JSON

### Marshalling

> **IMPORTANT**: Always pass pointers to `json.Marshall` and `json.Unmarshall`.

**Marshalling** transforms a memory representation of an object into the JSON data format for storage or transmission.

### Unmarshalling
**Unmarshalling** transforms a JSON object into a memory representation that is executable.

To unmarshall a JSON object into memory, pass the data and a pointer to the data structure that you want to store the data in:
```go
type person struct {
    name    string
    age     int
}

var jsonData := `[
{"name": "Steve", "age": "21"},
{"name": "Bob", "age": "68"}
]`

var unmarshalled []person

json.Unmarshall(data, &unmarshalled)
```

### go.mod and go.sum

Go modules group related packages into a single unit to be versioned together. Because they track an application's dependencies, they ensure that users build the application with the same dependencies as the original developer.

Go modules allow you to write go programs outside of the $GOPATH directory, as in previous releases. 

Go sum records the checksum for each module in the application to ensure that each build uses the correct version.

## Structure



## Reading data

### Writing from a file

Read data from a file with the `os` package. `ReadFile` reads the contents of a file and returns a `nil` error:
```go
os.ReadFile(filename)
```

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
## Writing data

### Writing to a file

Write data to a file with the `os` package. `WriteFile` writes to an existing file or creates one, if necessary:
```go
os.WriteFile(filename, dataToWrite, perms)
```
> **Linux permissions**: Set Linux file permissions for the file owner, group, and all other users (`owner-group-others`). The permission options are read, write, and execute. Use octal values to set permssions:  
  read: 4  
  write: 2  
  execute: 1  

When you assign permissions in an programming language, you have to tell the program that you are using the octal base. You do this by beginning the number with a `0`. So, `0644` permissions means that the file owner has read and write permissions, and the group and all other users have read permissions.

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
## Time

Get the current time:
```go
current = time.Now()
```
Get the zero value for time.Time with an empty struct:
```go
zeroVal = time.Time{}
```
# Tests

