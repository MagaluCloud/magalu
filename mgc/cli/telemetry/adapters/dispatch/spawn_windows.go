//go:build windows

package dispatch

import "syscall"

// detachedProcess é o DETACHED_PROCESS da API do Windows, que não está no pacote syscall
// usado para desvincular o processo do terminal
const detachedProcess = 0x00000008

// O processo de envio nasce sem console e fora do grupo da CLI, então fechar o
// console ou dar Ctrl+C não o encerra
func sysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | detachedProcess,
		HideWindow:    true,
	}
}
