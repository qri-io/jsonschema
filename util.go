package jsonschema

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

var showDebug = os.Getenv("JSON_SCHEMA_DEBUG") == "1"

// schemaDebug provides a logging interface
// which is off by defauly but can be activated
// for debuging purposes
func schemaDebug(message string, args ...interface{}) {
	if showDebug {
		if message[len(message)-1] != '\n' {
			message += "\n"
		}
		fmt.Printf(message, args...)
	}
}

// SafeResolveURL resolves a string url against the current context url
func SafeResolveURL(ctxURL, resURL string) (string, error) {
	cu, err := url.Parse(ctxURL)
	if err != nil {
		return "", err
	}
	u, err := url.Parse(resURL)
	if err != nil {
		return "", err
	}
	resolvedURL := cu.ResolveReference(u)
	if resolvedURL.Scheme == "file" && cu.Scheme != "file" {
		return "", fmt.Errorf("cannot access file resources from network context")
	}
	resolvedURLString := resolvedURL.String()
	return resolvedURLString, nil
}

// IsLocalSchemaID validates if a given id is a local id
func IsLocalSchemaID(id string) bool {
	splitID := strings.Split(id, "#")
	if len(splitID) > 1 && len(splitID[0]) > 0 && splitID[0][0] != '#' {
		return false
	}
	return id != "#" && !strings.HasPrefix(id, "#/") && strings.Contains(id, "#")
}
