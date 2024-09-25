package main

import (
    "bytes"
    "fmt"
    "os"
    "strings"
    "unicode/utf8"

    "golang.org/x/term"
)

func main() {
	const s = 
`
  /\     /\  
 /  \___/  \ 
(  o     o  )
(     ^     )
 \_________/ 
  /       \  
 /         \ 
/_/|     |\_\
`
    
    // Center multi-line string
    buf := CenterMultiLine(s)
    fmt.Println(buf)
}

// NCenter centers a single line string to the column width.
func NCenter(width int, s string) *bytes.Buffer {
    const half, space = 2, "\u0020"
    var b bytes.Buffer
    n := (width - utf8.RuneCountInString(s)) / half
    if n < 1 {
        fmt.Fprintf(&b, s)
        return &b
    }
    fmt.Fprintf(&b, "%s%s", strings.Repeat(space, n), s)
    return &b
}

// CenterMultiLine centers each line of the multi-line string.
func CenterMultiLine(s string) *bytes.Buffer {
    fd := int(os.Stdin.Fd())
    w, _, err := term.GetSize(fd)
    if err != nil {
        return NCenter(0, s)
    }

    var b bytes.Buffer
    lines := strings.Split(s, "\n")
    for _, line := range lines {
        // Center each line and write to buffer
        centeredLine := NCenter(w, line)
        b.Write(centeredLine.Bytes())
        b.WriteString("\n") // Add a newline after each centered line
    }
    return &b
}
