package typespec

import "os"

func WriteToFile(path string, data string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = f.WriteString(data)
	return err
}
