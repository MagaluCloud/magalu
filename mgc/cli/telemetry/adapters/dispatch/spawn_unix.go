//go:build !windows

package dispatch

import "syscall"

// Uma sessão nova desliga o processo de envio do terminal, então o Ctrl+C e o
// fechamento do terminal não chegam a ele
func sysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setsid: true}
}
