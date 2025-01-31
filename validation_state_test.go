package jsonschema_test

import (
	"context"
	"encoding/json"
	"testing"

	jptr "github.com/qri-io/jsonpointer"
	"github.com/qri-io/jsonschema"
)

const GeoLockKeyword = "geolock"
const GeoLockADKey = "geolockrefs"

type GeoLock struct {
	Scope string
	Zone  string
}

type GeoLockRef struct {
	DataPath string
	Scope    string
	Zone     string
	Value    string
}

func NewGeoLock() jsonschema.Keyword {
	return &GeoLock{}
}

func (f *GeoLock) Validate(propPath string, data interface{}, errs *[]jsonschema.KeyError) {}

func (f *GeoLock) Register(uri string, registry *jsonschema.SchemaRegistry) {}

func (f *GeoLock) Resolve(pointer jptr.Pointer, uri string) *jsonschema.Schema {
	return nil
}

func (f *GeoLock) ValidateKeyword(ctx context.Context, currentState *jsonschema.ValidationState, data interface{}) {
	glKeyword := currentState.Local.JSONProp("geolock").(*GeoLock)

	var refs []GeoLockRef
	if tmp, found := currentState.GetAdditionalValidationData(GeoLockADKey); found {
		refs = tmp.([]GeoLockRef)
	} else {
		refs = []GeoLockRef{}
	}

	ref := GeoLockRef{
		DataPath: currentState.InstanceLocation.String(),
		Scope:    glKeyword.Scope,
		Zone:     glKeyword.Zone,
		Value:    data.(string),
	}

	refs = append(refs, ref)

	currentState.SetAdditionalValidationData(GeoLockADKey, refs)
}

func TestAdditionalValidationData(t *testing.T) {

	ctx := context.Background()
	var schemaData = []byte(`{
		"$defs": {
			"person": {
				"type": "object",
				"properties": {
					"firstName": {
						"type": "string"
					},
					"lastName": {
						"type": "string"
					},
					"address_city": {
						"type": "string",
						"geolock": {
							"zone": "Europe",
							"scope": "city"
						}
					}
				}
			}
		},
		"type": "array",
		"items": {
			"$ref": "#/$defs/person"
		}
	}`)

	jsonschema.RegisterKeyword(GeoLockKeyword, NewGeoLock)
	jsonschema.LoadDraft2019_09()

	rs := &jsonschema.Schema{}
	if err := rs.UnmarshalJSON(schemaData); err != nil {
		t.Errorf("failure to unmarshal schema: %s", err.Error())
		t.FailNow()
		return
	}

	expectedCount := 3
	var valid = []byte(`[
		{ "firstName" : "Antonio", "lastName" : "Alves", "address_city": "pt#lisbon" },
		{ "firstName" : "Peter", "lastName" : "Parque", "address_city": "de#berlin" },
		{ "firstName" : "Juan", "lastName" : "Juniper", "address_city": "es#madrid" }
	]`)
	var doc interface{}
	if err := json.Unmarshal(valid, &doc); err != nil {
		t.Errorf("failure to unmarshal test data: %s", err.Error())
		t.FailNow()
		return
	}
	vs := rs.Validate(ctx, doc)
	if vs.Errs != nil && len(*vs.Errs) > 0 {
		t.Errorf("failure to validate test data: %v", vs.Errs)
		t.FailNow()
		return
	}

	tmp, refExists := vs.GetAdditionalValidationData(GeoLockADKey)
	if !refExists {
		t.Errorf("Expected to find key %q validation state additional data.", GeoLockADKey)
		t.FailNow()
		return
	}
	foundLocks := tmp.([]GeoLockRef)

	if len(foundLocks) != expectedCount {
		t.Errorf("Unexpected number of refs found. expected=%d, found=%d", expectedCount, len(foundLocks))
		t.FailNow()
		return
	}
}
