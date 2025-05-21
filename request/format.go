package request

import "fmt"

func GetSSE(path string) string {
	return fmt.Sprintf("@get('%s')", path)
}
