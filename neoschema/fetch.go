package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"net/url"
)

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
