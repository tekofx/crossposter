package logger

import (
	"fmt"
	"os"
	"time"
)

func Log(msgs ...any) {
	output := fmt.Sprintf("%s %s", time.Now().Format("2006-01-02 15:04:05"), msgs)
	fmt.Println(output)
}

func Fatal(msgs ...any) {
	Log(append([]any{" ☠️ | Fatal:"}, msgs...)...)
	os.Exit(1)
}

func Error(msgs ...any) {
	Log(append([]any{" ❌ | Error:"}, msgs...)...)
}

func Warning(msgs ...any) {
	Log(append([]any{" ⚠️  | Warning:"}, msgs...)...)
}
