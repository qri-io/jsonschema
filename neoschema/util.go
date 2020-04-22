package main

import (
	"fmt"
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