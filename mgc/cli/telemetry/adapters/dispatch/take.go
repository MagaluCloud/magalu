package dispatch

import "os"

// Take lê o payload deixado pela CLI e apaga o arquivo em qualquer caso, porque ele
// foi escrito só para esta execução
func Take(path string) (Payload, error) {
	data, err := os.ReadFile(path)
	_ = os.Remove(path)
	if err != nil {
		return Payload{}, err
	}
	return Decode(data)
}
