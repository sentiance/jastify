package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fgm/jastify/converter"
)

// sanitizeResourceName converts a file path to a valid Terraform resource name
func sanitizeResourceName(name string) string {
	// Remove file extension and get base name
	base := strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))

	// Replace invalid characters with underscores
	reg := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	sanitized := reg.ReplaceAllString(base, "_")

	// Ensure it starts with a letter or underscore
	if len(sanitized) > 0 && !regexp.MustCompile(`^[a-zA-Z_]`).MatchString(sanitized) {
		sanitized = "_" + sanitized
	}

	// Ensure it's not empty
	if sanitized == "" {
		sanitized = "resource"
	}

	return sanitized
}

func main() {
	var (
		jsonData     []byte
		err          error
		resourceName string
	)

	switch len(os.Args) {
	case 1:
		// No args: read from stdin
		jsonData, err = io.ReadAll(os.Stdin)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
			os.Exit(1)
		}
	case 2:
		// One arg: it's a file, read from it.
		jsonData, err = os.ReadFile(os.Args[1])
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Error reading file:", err)
			os.Exit(1)
		}
		resourceName = sanitizeResourceName(os.Args[1])
	default:
		_, _ = fmt.Fprintln(os.Stderr, "Usage: \nconvert < somefile.json\nor\nconvert somefile.json")
		os.Exit(1)
	}

	var parsedJson converter.Jmap
	if err := json.Unmarshal(jsonData, &parsedJson); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error parsing JSON:", err)
		os.Exit(1)
	}

	var tf string

	if _, exists := parsedJson["name"]; exists {
		if resourceName == "" {
			resourceName = "monitor_1"
		}
		tf = converter.Must(converter.GenerateMonitorTerraformCode(resourceName, parsedJson))
	} else {
		if resourceName == "" {
			resourceName = "dashboard_1"
		}
		tf = converter.Must(converter.GenerateDashboardTerraformCode(resourceName, parsedJson))
	}

	fmt.Println(tf)
}
