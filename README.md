# CBLogForwarder

## Install
`go get -u github.com/nicouch/CBLogForwarder`

## Usage
CBlogForwarder takes multiple couchbase log files as input, and "output" a single filtered log stream.
It parse every logs incoming newline and output to stdout only selected indices (spit char is specified in configuration for each file)

Configuration file must be named `CBLogForwarder.conf` and be located where the program is called.

## Configuration Example
```
{
    "files": [
        {
            "file": "fakeFile1.log",
            "splitOn": " ",
            "outputIndices": [3, 4, 0, 1, 2]
        },
        {
            "file": "fakeFile2.log",
            "splitOn": " ",
            "outputIndices": [0, 3, 2]
        }
    ]
}
```
