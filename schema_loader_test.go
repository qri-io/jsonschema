package jsonschema_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/qri-io/jsonschema"
)

func createTestServer() *httptest.Server {
	validSchema := `{
		"type": "string"
	}`
	invalidSchema := "invalid_schema"
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/valid_schema.json":
				fmt.Fprintln(w, validSchema)
			default:
				fmt.Fprintln(w, invalidSchema)
			}
		}))
}

func TestFetchSchema(t *testing.T) {

	ts := createTestServer()
	defer ts.Close()

	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("failed to get working dir: %q", err)
	}

	cases := []struct {
		uri          string
		expectsError bool
		message      string
	}{
		{fmt.Sprintf("%s/valid_schema.json", ts.URL), false, ""},
		{fmt.Sprintf("%s/invalid_schema.json", ts.URL), true, "invalid character"},
		{fmt.Sprintf("file://%s/testdata/draft-07_schema.json", wd), false, ""},
		{fmt.Sprintf("file://%s/testdata/missing_file.json", wd), true, "no such file or directory"},
		{"unknownscheme://resource.json#definitions/property", true, "unknownscheme is not supported for uri"},
	}

	for i, c := range cases {
		rs := &jsonschema.Schema{}
		err := jsonschema.FetchSchema(context.Background(), c.uri, rs)

		if !c.expectsError && err == nil {
			continue
		}

		if c.expectsError {
			if err == nil {
				t.Errorf("case %d expected an error", i)
				continue
			}

			if !strings.Contains(err.Error(), c.message) {
				t.Errorf("case %d expected error to include %q actual: %q", i, c.message, err.Error())
				continue
			}

		} else if err != nil {
			t.Errorf("case %d unexpected error: %s", i, err)
			continue
		}

	}

}

func TestCustomSchemaLoader(t *testing.T) {

	lr := jsonschema.GetSchemaLoaderRegistry()
	lr.Register("special", func(ctx context.Context, uri *url.URL, schema *jsonschema.Schema) error {

		path := uri.Host + uri.Path
		body := fmt.Sprintf(`{ "type": "string", "description": "example description for '%s'"}`, path)
		if schema == nil {
			schema = &jsonschema.Schema{}
		}
		return json.Unmarshal([]byte(body), schema)

	})

	resourceURI := "special://schema_name"
	rs := &jsonschema.Schema{}
	err := jsonschema.FetchSchema(context.Background(), resourceURI, rs)

	if err != nil {
		t.Errorf("failed to load schema: %s", err)
		return
	}

	if rs.TopLevelType() != "string" {
		t.Errorf("expected schema top level type to be %q, actual: %q", "string", rs.TopLevelType())
	}

	expectedDesc := "example description for 'schema_name'"
	actualDesc := string(*rs.JSONProp("description").(*jsonschema.Description))
	if actualDesc != expectedDesc {
		t.Errorf("expected 'description' to be %q, actual: %v", expectedDesc, actualDesc)
	}

}
