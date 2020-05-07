package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func SchemaDebug(message string) {
	debug := false
	debugEnvVar := os.Getenv("JSON_SCHEMA_DEBUG")
	if debugEnvVar == "1" {
		debug = true
	}
	if debug {
		fmt.Printf("%s\n", message)
	}
}

func SafeResolveUrl(ctxUrl, resUrl string) (string, error) {
	cu, err := url.Parse(ctxUrl)
	if err != nil {
		return "", err
	}
	u, err := url.Parse(resUrl)
	if err != nil {
		return "", err
	}
	resolvedUrl := cu.ResolveReference(u)
	if resolvedUrl.Scheme == "file" && cu.Scheme != "file" {
		return "", fmt.Errorf("cannot access file resources from network context")
	}
	resolvedUrlString := resolvedUrl.String()
	return resolvedUrlString, nil
}

func IsLocalSchemaId(id string) bool {
	splitId := strings.Split(id, "#")
	if len(splitId) > 1 && len(splitId[0]) > 0 && splitId[0][0] != '#' {
		return false
	}
	return id != "#" && !strings.HasPrefix(id, "#/") && strings.Contains(id, "#")
}

func FetchSchema(ctx *context.Context, uri string, schema *Schema) error {
	SchemaDebug(fmt.Sprintf("[FetchSchema] Fetching: %s", uri))
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	if u.Scheme == "http" || u.Scheme == "https" {
		var req *http.Request
		if ctx != nil {
			req, _ = http.NewRequestWithContext(*ctx, "GET", u.String(), nil)
		}  else {
			req, _ = http.NewRequest("GET", u.String(), nil)	
		}
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		if schema == nil {
			schema = &Schema{}
		}
		return json.Unmarshal(body, schema)
	} else {
		return fmt.Errorf("URI scheme %s is not supported for uri: %s", u.Scheme, uri)
	}
}
