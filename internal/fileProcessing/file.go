package fileprocessing

import "os"

func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
