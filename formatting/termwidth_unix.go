package formatting

import (
  "golang.org/x/crypto/ssh/terminal"
  "os"
)

func TermWidth() int {
  stdout := int(os.Stdout.Fd())
  if !terminal.IsTerminal(stdout) {
    return 0
  }
  w, _, err := terminal.GetSize(stdout)
  if err != nil || w <= 0 {
    return 0
  }
  return w
}
