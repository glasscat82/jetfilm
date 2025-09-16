package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// color: 1 red, 2 green, 3 yello, 4 blue, 5 purple, 6 blue
func p(color int, sep string, str ...any) {
	newStr := []any{}
	for index, v := range str {
		if index == 0 {
			newStr = append(newStr, v)
		} else {
			newStr = append(newStr, sep, v)
		}
	}

	suffixColor := "\033[3" + strconv.Itoa(color) + "m"
	fmt.Printf("%s%s%s", suffixColor, fmt.Sprint(newStr...), "\033[0m\n")
}

func writeJson(data any, filename string) error {
	// write to files
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Errorf("error marchal %s: %w", filename, err)
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error write file %s: %w", filename, err)
	}

	fmt.Printf("JSON data write to file ~ %s\n", filename)
	return nil
}
