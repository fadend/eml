// Code to parse .eml files such as one gets from Gmail's "Download Original".
// The main focus for this package is to extract attachments from such files.
package eml

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

var ErrMissingBoundary = errors.New("missing boundary")

func extractBoundary(header string) string {
	header = strings.TrimSpace(header)
	if rest, matched := strings.CutPrefix(header, "boundary=\""); matched {
		if boundary, matched := strings.CutSuffix(rest, "\""); matched {
			return boundary
		}
	}
	return ""
}

func findBoundary(s *bufio.Scanner) (string, error) {
	foundPrefix := false
	for s.Scan() {
		line := s.Text()
		// Reached end of header without finding boundary.
		if line == "" {
			return "", ErrMissingBoundary
		}
		if suffix, matched := strings.CutPrefix(line, "Content-Type: multipart/mixed; "); matched {
			foundPrefix = true
			boundary := extractBoundary(suffix)
			if len(boundary) > 0 {
				return boundary, nil
			}
		} else if foundPrefix {
			boundary := extractBoundary(line)
			if len(boundary) > 0 {
				return boundary, nil
			}
			return "", ErrMissingBoundary
		}
	}
	return "", ErrMissingBoundary
}

// ExtractFileNameToAttachment returns a filename to byte content map for the attachments found.
func ExtractFileNameToAttachment(r io.Reader) (map[string][]byte, error) {
	scanner := bufio.NewScanner(r)
	boundary, err := findBoundary(scanner)
	if err != nil {
		return nil, err
	}
	boundary = fmt.Sprintf("--%s", boundary)
	hitBlank := false
	fileName := ""
	withinContent := false
	result := make(map[string][]byte)
	var builder strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		if remainder, match := strings.CutPrefix(line, boundary); match && strings.Trim(remainder, "-") == "" {
			hitBlank = false
			withinContent = true
			if fileName != "" {
				content, err := base64.StdEncoding.DecodeString(builder.String())
				if err != nil {
					return nil, err
				}
				result[fileName] = content
			}
			builder.Reset()
			fileName = ""
		} else if !withinContent {
			continue
		} else if line == "" {
			hitBlank = true
		} else if hitBlank && fileName != "" {
			builder.WriteString(line)
		} else if fileNameSuffix, match := strings.CutPrefix(line, "Content-Disposition: attachment; filename=\""); match {
			fileName = strings.Trim(fileNameSuffix, "\"")
		}
	}
	return result, nil
}
