package logging

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

// Println will print the current file & calling line including timestamp.
func Println(format string, v ...any) {
	_, file, line, _ := runtime.Caller(1)
	log.Printf("[%s:%d] %s\r\n", filepath.Base(file), line, fmt.Sprintf(format, v...))
}
