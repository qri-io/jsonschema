package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func SplitUrl(url string) (string, string) {
	urlSlice := strings.SplitN(url, "#", 2)
	ref := ""
	fragment := ""
	if len(urlSlice) > 0 {
		ref = urlSlice[0]
	}
	if len(urlSlice) > 1 {
		fragment = urlSlice[1]
	}
	return ref, fragment
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
	return id != "#" && !strings.Contains(id, "#/") && strings.Contains(id, "#")
}

func FetchSchema(uri string, schema *Schema) error {
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	if u.Scheme == "http" || u.Scheme == "https" {
		res, err := http.Get(u.String())
		if err != nil {
			return err
		}
		return json.NewDecoder(res.Body).Decode(schema)
	} else {
		return fmt.Errorf("URI scheme %s is not supported", u.Scheme)
	}
}
