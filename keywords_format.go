package jsonschema

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
)

const (
	email    string = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	hostname string = `^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])(\.([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9]))*$`
)

var (
	emailPattern    = regexp.MustCompile(email)
	hostnamePattern = regexp.MustCompile(hostname)
)

// func FormatType(data interface{}) string {
// 	switch
// }
// Note: Date and time format names are derived from RFC 3339, section
// 5.6  [RFC3339].
// http://json-schema.org/latest/json-schema-validation.html#RFC3339

type format string

func newFormat() Validator {
	return new(format)
}

func (f format) Validate(data interface{}) error {
	if str, ok := data.(string); ok {
		switch f {
		case "date-time":
			return isValidDateTime(str)
		case "date":
			return isValidDate(str)
		case "email":
			return isValidEmail(str)
		case "hostname":
			return isValidHostname(str)
		case "idn-email":
			return isValidIdnEmail(str)
		case "idn-hostname":
			return isValidIdnHostname(str)
		case "ipv4":
			return isValidIPv4(str)
		case "ipv6":
			return isValidIPv6(str)
		case "iri-reference":
			return isValidIriRef(str)
		case "iri":
			return isValidIri(str)
		case "json-pointer":
			return isValidJSONPointer(str)
		case "regex":
			return isValidRegex(str)
		case "relative-json-pointer":
			return isValidRelJSONPointer(str)
		case "time":
			return isValidTime(str)
		case "uri-reference":
			return isValidURIRef(str)
		case "uri-template":
			return isValidURITemplate(str)
		case "uri":
			return isValidURI(str)
		}
	}
	return nil
}

func isValidDateTime(dateTime string) error {
	if _, err := time.Parse(time.RFC3339, dateTime); err != nil {
		return fmt.Errorf("date-time incorrectly formatted: %s", err.Error())
	}
	return nil
}

func isValidDate(date string) error {
	arbitraryTime := "T08:30:06.283185Z"
	dateTime := fmt.Sprintf("%s%s", date, arbitraryTime)
	return isValidDateTime(dateTime)
}

func isValidEmail(email string) error {
	if !emailPattern.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}
func isValidHostname(hostname string) error {
	if !hostnamePattern.MatchString(hostname) || len(hostname) > 255 {
		return fmt.Errorf("invalid hostname string")
	}
	return nil
}
func isValidIdnEmail(idnEmail string) error {
	return nil
}
func isValidIdnHostname(idnHostname string) error {
	return nil
}
func isValidIPv4(ipv4 string) error {
	parsedIP := net.ParseIP(ipv4)
	hasDots := strings.Contains(ipv4, ".")
	if !hasDots || parsedIP == nil {
		return fmt.Errorf("invalid IPv4 address")
	}
	return nil
}

func isValidIPv6(ipv6 string) error {
	parsedIP := net.ParseIP(ipv6)
	hasColons := strings.Contains(ipv6, ":")
	if !hasColons || parsedIP == nil {
		return fmt.Errorf("invalid IPv4 address")
	}
	return nil
}

func isValidIriRef(iriRef string) error {
	return nil
}
func isValidIri(iri string) error {
	return nil
}
func isValidJSONPointer(jsonPointer string) error {
	return nil
}
func isValidRegex(regex string) error {
	return nil
}
func isValidRelJSONPointer(relJsonPointer string) error {
	return nil
}
func isValidTime(time string) error {
	return nil
}
func isValidURIRef(uriRef string) error {
	return nil
}
func isValidURITemplate(uriTemplate string) error {
	return nil
}
func isValidURI(uri string) error {
	return nil
}
