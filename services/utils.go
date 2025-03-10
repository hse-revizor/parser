package services

import (
	"fmt"
	"strings"
)

func ParseRowParam(rowParam string) (startRow, endRow int, err error) {
	if rowParam == "" {
		return 0, 0, nil // No row parameter, return the entire content
	}

	if strings.Contains(rowParam, "-") {
		_, err := fmt.Sscanf(rowParam, "%d-%d", &startRow, &endRow)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid row range: %s", rowParam)
		}
		if startRow < 1 || endRow < 1 || startRow > endRow {
			return 0, 0, fmt.Errorf("invalid row range: %s (must be in format start-end, where start <= end)", rowParam)
		}

		return startRow, endRow, nil
	}

	startRow = 0
	endRow = 0
	_, err = fmt.Sscanf(rowParam, "%d", &startRow)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid row number: %s", rowParam)
	}
	if startRow < 1 {
		return 0, 0, fmt.Errorf("row number must be greater than 0")
	}

	return startRow, startRow, nil
}

func ExtractRow(content string, row int) (string, error) {
	lines := strings.Split(content, "\n")
	if row > len(lines) {
		return "", fmt.Errorf("row number %d is out of range (max %d)", row, len(lines))
	}

	return lines[row-1], nil
}

func ExtractRowRange(content string, startRow, endRow int) (string, error) {
	lines := strings.Split(content, "\n")
	if startRow > len(lines) || endRow > len(lines) {
		return "", fmt.Errorf("row range %d-%d is out of range (max %d)", startRow, endRow, len(lines))
	}

	if startRow > endRow {
		return "", fmt.Errorf("invalid row range: start (%d) must be less than or equal to end (%d)", startRow, endRow)
	}

	selectedLines := lines[startRow-1 : endRow]
	return strings.Join(selectedLines, "\n"), nil
}
