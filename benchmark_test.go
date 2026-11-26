package jsonschema

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// BenchmarkEnumString looks up the middle value of an enum of strings.
func BenchmarkEnumString(b *testing.B) {
	runBenchmark(b,
		func(sampleSize int) (string, interface{}) {
			data := make([]interface{}, sampleSize)
			for i := 0; i < sampleSize; i++ {
				data[i] = fmt.Sprintf("value-%v", i)
			}
			d, err := json.Marshal(data)
			if err != nil {
				b.Errorf("unable to marshal data: %v", err)
				return "", nil
			}
			return `{
				"enum": ` + string(d) + `
			}`, fmt.Sprintf("value-%v", sampleSize/2)
		},
	)
}

// BenchmarkOneOfDiscriminator validates an object against a oneOf whose
// branches are told apart by a required property with a single-value enum.
// The object matches the middle branch.
func BenchmarkOneOfDiscriminator(b *testing.B) {
	runBenchmark(b,
		func(sampleSize int) (string, interface{}) {
			data := map[string]interface{}{
				"kind":  fmt.Sprintf("kind-%v", sampleSize/2),
				"value": "foobar",
			}
			var branches strings.Builder

			for i := 0; i < sampleSize; i++ {
				branches.WriteString(fmt.Sprintf(`{
					"type": "object",
					"properties": {
						"kind": { "enum": ["kind-%v"] },
						"value": { "type": "string" }
					},
					"required": ["kind", "value"],
					"additionalProperties": false
				}`, i))

				if i != sampleSize-1 {
					branches.WriteString(",")
				}
			}

			return `{
				"oneOf": [ ` + branches.String() + ` ]
			}`, data
		},
	)
}

// requestSchema is shaped like the schema of an API request body: nested
// objects and arrays that are closed with "additionalProperties": false, and
// polymorphic objects expressed as a oneOf whose branches are told apart by a
// required property with a single-value enum.
const requestSchema = `{
	"$schema": "https://json-schema.org/draft/2019-09/schema",
	"title": "Order",
	"type": "object",
	"additionalProperties": false,
	"required": ["id", "createdAt", "version", "status", "currency", "customer", "items", "payment"],
	"properties": {
		"id": { "type": "string", "format": "uuid" },
		"createdAt": { "type": "string", "format": "date-time" },
		"version": { "const": 2 },
		"status": { "type": "string", "enum": ["draft", "pending", "paid", "shipped", "delivered", "cancelled"] },
		"currency": { "type": "string", "pattern": "^[A-Z]{3}$" },
		"note": { "type": "string", "maxLength": 500 },
		"customer": {
			"type": "object",
			"additionalProperties": false,
			"required": ["id", "email", "name"],
			"properties": {
				"id": { "type": "integer", "minimum": 1 },
				"email": { "type": "string", "format": "email", "maxLength": 254 },
				"name": { "type": "string", "minLength": 1, "maxLength": 100 },
				"phone": { "type": "string", "pattern": "^\\+[0-9]{6,15}$" },
				"address": {
					"type": "object",
					"additionalProperties": false,
					"required": ["street", "city", "postalCode", "country"],
					"properties": {
						"street": { "type": "string", "minLength": 1, "maxLength": 200 },
						"city": { "type": "string", "minLength": 1, "maxLength": 100 },
						"postalCode": { "type": "string", "pattern": "^[0-9A-Z -]{3,10}$" },
						"country": { "type": "string", "pattern": "^[A-Z]{2}$" }
					}
				}
			}
		},
		"items": {
			"type": "array",
			"minItems": 1,
			"maxItems": 100,
			"items": {
				"type": "object",
				"additionalProperties": false,
				"required": ["sku", "name", "quantity", "unitPrice"],
				"properties": {
					"sku": { "type": "string", "pattern": "^[A-Z]{3}-[0-9]{5}$" },
					"name": { "type": "string", "minLength": 1, "maxLength": 200 },
					"quantity": { "type": "integer", "minimum": 1, "maximum": 1000 },
					"unitPrice": { "type": "number", "minimum": 0 },
					"category": { "enum": ["books", "electronics", "clothing", "food", "toys", "other"] },
					"tags": {
						"type": "array",
						"maxItems": 10,
						"items": { "type": "string", "minLength": 1, "maxLength": 32 }
					}
				}
			}
		},
		"payment": {
			"oneOf": [
				{
					"type": "object",
					"additionalProperties": false,
					"required": ["method", "cardLast4", "cardBrand"],
					"properties": {
						"method": { "enum": ["card"] },
						"cardLast4": { "type": "string", "pattern": "^[0-9]{4}$" },
						"cardBrand": { "enum": ["visa", "mastercard", "amex"] }
					}
				},
				{
					"type": "object",
					"additionalProperties": false,
					"required": ["method", "iban"],
					"properties": {
						"method": { "enum": ["bank_transfer"] },
						"iban": { "type": "string", "minLength": 15, "maxLength": 34 },
						"reference": { "type": "string", "maxLength": 140 }
					}
				},
				{
					"type": "object",
					"additionalProperties": false,
					"required": ["method", "dueDate"],
					"properties": {
						"method": { "enum": ["invoice"] },
						"dueDate": { "type": "string", "format": "date" }
					}
				}
			]
		},
		"adjustments": {
			"type": "array",
			"maxItems": 20,
			"items": {
				"oneOf": [
					{
						"type": "object",
						"additionalProperties": false,
						"required": ["kind", "code", "percent"],
						"properties": {
							"kind": { "enum": ["discount"] },
							"code": { "type": "string", "pattern": "^[A-Z0-9]{4,12}$" },
							"percent": { "type": "number", "minimum": 0, "maximum": 100 }
						}
					},
					{
						"type": "object",
						"additionalProperties": false,
						"required": ["kind", "carrier", "amount"],
						"properties": {
							"kind": { "enum": ["shipping"] },
							"carrier": { "enum": ["post", "courier", "pickup"] },
							"amount": { "type": "number", "minimum": 0 }
						}
					},
					{
						"type": "object",
						"additionalProperties": false,
						"required": ["kind", "rate", "region"],
						"properties": {
							"kind": { "enum": ["tax"] },
							"rate": { "type": "number", "minimum": 0, "maximum": 1 },
							"region": { "type": "string", "pattern": "^[A-Z]{2}(-[A-Z0-9]{1,3})?$" }
						}
					}
				]
			}
		}
	}
}`

// requestDocument builds a request body for requestSchema with the given
// number of line items. Twelve line items give about 3 KB of JSON.
func requestDocument(lineItems int) string {
	categories := []string{"books", "electronics", "clothing", "food", "toys", "other"}

	var items strings.Builder
	for i := 0; i < lineItems; i++ {
		if i > 0 {
			items.WriteString(",")
		}
		fmt.Fprintf(&items, `{
			"sku": "ABC-%05d",
			"name": "Line item number %d with a reasonably long product name",
			"quantity": %d,
			"unitPrice": %d.95,
			"category": %q,
			"tags": ["tag-%d", "seasonal", "gift"]
		}`, i, i, i%7+1, i*3+4, categories[i%len(categories)], i)
	}

	return `{
		"id": "3f2b8c1e-7a4d-4e0b-9c56-1d2e3f4a5b6c",
		"createdAt": "2024-05-17T09:30:00Z",
		"version": 2,
		"status": "paid",
		"currency": "EUR",
		"note": "Please leave the parcel with the neighbours if nobody is at home.",
		"customer": {
			"id": 48213,
			"email": "ada.lovelace@example.com",
			"name": "Ada Lovelace",
			"phone": "+358401234567",
			"address": {
				"street": "Mannerheimintie 12 B 34",
				"city": "Helsinki",
				"postalCode": "00100",
				"country": "FI"
			}
		},
		"items": [` + items.String() + `],
		"payment": {
			"method": "invoice",
			"dueDate": "2024-06-16"
		},
		"adjustments": [
			{ "kind": "discount", "code": "SPRING24", "percent": 10 },
			{ "kind": "shipping", "carrier": "courier", "amount": 6.9 },
			{ "kind": "tax", "rate": 0.24, "region": "FI" }
		]
	}`
}

// requestBenchmarkData returns the parsed request schema, a decoded request
// body that is valid against it, and a copy of that body with a handful of
// errors spread over the keywords the schema uses.
func requestBenchmarkData(b *testing.B) (schema *Schema, valid, invalid interface{}) {
	ctx := context.Background()
	schema = &Schema{}
	if err := json.Unmarshal([]byte(requestSchema), schema); err != nil {
		b.Fatalf("error parsing schema: %s", err.Error())
	}

	decode := func(doc string, expectValid bool) (data interface{}) {
		if err := json.Unmarshal([]byte(doc), &data); err != nil {
			b.Fatalf("error parsing document: %s", err.Error())
		}
		if state := schema.Validate(ctx, data); state.IsValid() != expectValid {
			b.Fatalf("expected valid document: %v, got errors: %v", expectValid, *state.Errs)
		}
		return data
	}

	doc := requestDocument(12)
	valid = decode(doc, true)
	invalid = decode(strings.NewReplacer(
		`"status": "paid"`, `"status": "refunded"`,
		`"currency": "EUR"`, `"currency": "euro"`,
		`"country": "FI"`, `"country": "FI", "county": "Uusimaa"`,
		`"sku": "ABC-00003"`, `"sku": 3`,
		`"method": "invoice"`, `"method": "cash"`,
	).Replace(doc), false)
	return schema, valid, invalid
}

// BenchmarkValidateRequest measures Schema.Validate with a schema that is
// parsed once, on a request body of about 3 KB.
func BenchmarkValidateRequest(b *testing.B) {
	ctx := context.Background()
	schema, valid, invalid := requestBenchmarkData(b)

	b.Run("valid", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			schema.Validate(ctx, valid)
		}
	})
	b.Run("invalid", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			schema.Validate(ctx, invalid)
		}
	})
}

// BenchmarkValidateRequestParallel validates the valid request body from many
// goroutines that share one schema, the way a server validates requests.
func BenchmarkValidateRequestParallel(b *testing.B) {
	ctx := context.Background()
	schema, valid, _ := requestBenchmarkData(b)

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if !schema.Validate(ctx, valid).IsValid() {
				b.Error("expected the document to be valid")
				return
			}
		}
	})
}
