package common

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

const (
	FILENAME_CSV  = "students-export.csv"
	FILENAME_JSON = "students-export.json"

	CSV_FILE_EXTENSION  = ".csv"
	JSON_FILE_EXTENSION = ".json"
	XLSX_FILE_EXTENSION = ".xlsx"
	XLS_FILE_EXTENSION  = ".xls"
)

const (
	FILE_TYPE_CSV     = "csv"
	FILE_TYPE_JSON    = "json"
	FILE_TYPE_UNKNOWN = "unknown"
)

// DetermineFileTypeFromContent analyzes the content to determine file type
func DetermineFileTypeFromContent(reader io.ReadSeeker) string {
	// Save original position
	currentPosition, err := reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return FILE_TYPE_UNKNOWN
	}
	defer reader.Seek(currentPosition, io.SeekStart) // Restore position afterward

	// Read first 1024 bytes to analyze
	buf := make([]byte, 1024)
	n, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		return FILE_TYPE_UNKNOWN
	}

	sample := buf[:n]

	// Check if it's JSON by looking for JSON structure
	trimmedSample := bytes.TrimSpace(sample)
	if len(trimmedSample) > 0 {
		firstChar := trimmedSample[0]
		if (firstChar == '{' && bytes.Contains(trimmedSample, []byte{':'})) ||
			(firstChar == '[' && bytes.Contains(trimmedSample, []byte{'{'})) {
			return FILE_TYPE_JSON
		}
	}

	// Check if it looks like CSV by looking for comma-separated values and newlines
	// Count commas and newlines to ensure it's consistent with CSV format
	commaCount := bytes.Count(sample, []byte{','})
	newlineCount := bytes.Count(sample, []byte{'\n'})

	// Only consider it CSV if there are reasonable comma counts per line (at least one comma)
	// and more than one line
	if newlineCount > 0 && commaCount > 0 && commaCount/newlineCount >= 1 {
		// Check if first line looks like a header (no numeric values)
		lines := bytes.Split(sample, []byte{'\n'})
		if len(lines) > 0 {
			firstLine := lines[0]
			// Check if first line has commas and doesn't look like numeric data
			if bytes.Contains(firstLine, []byte{','}) &&
				!bytes.ContainsAny(firstLine, "0123456789") {
				return FILE_TYPE_CSV
			}
		}

		// Even without clear header, if it has commas and lines, likely CSV
		return FILE_TYPE_CSV
	}

	// Default
	return FILE_TYPE_UNKNOWN
}

// DetermineFileTypeFromMetadata examines URL and headers for file type clues
func DetermineFileTypeFromMetadata(fileURL string, headers http.Header) string {
	// Try to determine from URL
	lowerURL := strings.ToLower(fileURL)
	if strings.HasSuffix(lowerURL, ".csv") {
		return FILE_TYPE_CSV
	} else if strings.HasSuffix(lowerURL, ".json") {
		return FILE_TYPE_JSON
	}

	// Look for content disposition header which might have filename
	contentDisposition := headers.Get("Content-Disposition")
	if contentDisposition != "" {
		if strings.Contains(strings.ToLower(contentDisposition), ".csv") {
			return FILE_TYPE_CSV
		} else if strings.Contains(strings.ToLower(contentDisposition), ".json") {
			return FILE_TYPE_JSON
		}
	}

	// Check content type
	contentType := headers.Get("Content-Type")
	if strings.Contains(contentType, "csv") || strings.Contains(contentType, "text/comma-separated-values") {
		return FILE_TYPE_CSV
	} else if strings.Contains(contentType, "json") || strings.Contains(contentType, "application/json") {
		return FILE_TYPE_JSON
	} else if strings.Contains(contentType, "text/plain") {
		// Many CSVs are served as text/plain
		return FILE_TYPE_CSV
	}

	// Default
	return FILE_TYPE_UNKNOWN
}
