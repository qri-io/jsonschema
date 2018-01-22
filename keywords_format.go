package jsonschema

import (
	// "encoding/json"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	email          string = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	hostname       string = `^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])(\.([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9]))*$`
	unescapedTilda        = `\~[^01]`
	endingTilda           = `\~$`
)

var (
	emailPattern           = regexp.MustCompile(email)
	hostnamePattern        = regexp.MustCompile(hostname)
	unescaptedTildaPattern = regexp.MustCompile(unescapedTilda)
	endingTildaPattern     = regexp.MustCompile(endingTilda)
)

// for json pointers

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
		default:
			// TODO: should we return an error saying that we don't know that
			// format? or should we keep it as is (ignore, return nil)
			return nil
		}
	}
	return nil
}

// A string instance is valid against "date-time" if it is a valid
// representation according to the "date-time" production derived
// from RFC 3339, section 5.6 [RFC3339]
// https://tools.ietf.org/html/rfc3339#section-5.6
func isValidDateTime(dateTime string) error {
	if _, err := time.Parse(time.RFC3339, dateTime); err != nil {
		return fmt.Errorf("date-time incorrectly formatted: %s", err.Error())
	}
	return nil
}

// A string instance is valid against "date" if it is a valid
// representation according to the "full-date" production derived
// from RFC 3339, section 5.6 [RFC3339]
// https://tools.ietf.org/html/rfc3339#section-5.6
func isValidDate(date string) error {
	arbitraryTime := "T08:30:06.283185Z"
	dateTime := fmt.Sprintf("%s%s", date, arbitraryTime)
	return isValidDateTime(dateTime)
}

// A string instance is valid against "email" if it is a valid
// representation as defined by RFC 5322, section 3.4.1 [RFC5322].
// https://tools.ietf.org/html/rfc5322#section-3.4.1
func isValidEmail(email string) error {
	if !emailPattern.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// A string instance is valid against "hostname" if it is a valid
// representation as defined by RFC 1034, section 3.1 [RFC1034],
// including host names produced using the Punycode algorithm
// specified in RFC 5891, section 4.4 [RFC5891].
// https://tools.ietf.org/html/rfc1034#section-3.1
// https://tools.ietf.org/html/rfc5891#section-4.4
func isValidHostname(hostname string) error {
	if !hostnamePattern.MatchString(hostname) || len(hostname) > 255 {
		return fmt.Errorf("invalid hostname string")
	}
	return nil
}

// A string instance is valid against "idn-email" if it is a valid
// representation as defined by RFC 6531 [RFC6531]
// https://tools.ietf.org/html/rfc6531
func isValidIdnEmail(idnEmail string) error {
	return nil
}

// A string instance is valid against "hostname" if it is a valid
// representation as defined by either RFC 1034 as for hostname, or
// an internationalized hostname as defined by RFC 5890, section
// 2.3.2.3 [RFC5890].
// https://tools.ietf.org/html/rfc1034
// http://tools.ietf.org/html/rfc5890#section-2.3.2.3
func isValidIdnHostname(idnHostname string) error {
	return nil
}

// A string instance is valid against "ipv4" if it is a valid
// representation of an IPv4 address according to the "dotted-quad"
// ABNF syntax as defined in RFC 2673, section 3.2 [RFC2673].
// https://tools.ietf.org/html/rfc2673#section-3.2
func isValidIPv4(ipv4 string) error {
	parsedIP := net.ParseIP(ipv4)
	hasDots := strings.Contains(ipv4, ".")
	if !hasDots || parsedIP == nil {
		return fmt.Errorf("invalid IPv4 address")
	}
	return nil
}

// A string instance is valid against "ipv6" if it is a valid
// representation of an IPv6 address as defined in RFC 4291, section
// 2.2 [RFC4291].
// https://tools.ietf.org/html/rfc4291#section-2.2
func isValidIPv6(ipv6 string) error {
	parsedIP := net.ParseIP(ipv6)
	hasColons := strings.Contains(ipv6, ":")
	if !hasColons || parsedIP == nil {
		return fmt.Errorf("invalid IPv4 address")
	}
	return nil
}

// A string instance is a valid against "iri-reference" if it is a
// valid IRI Reference (either an IRI or a relative-reference),
// according to [RFC3987].
// https://tools.ietf.org/html/rfc3987
func isValidIriRef(iriRef string) error {
	return nil
}

// A string instance is a valid against "iri" if it is a valid IRI,
// according to [RFC3987].
// https://tools.ietf.org/html/rfc3987
func isValidIri(iri string) error {
	return nil
}

// A string instance is a valid against "json-pointer" if it is a
// valid JSON string representation of a JSON Pointer, according to
// RFC 6901, section 5 [RFC6901].
// https://tools.ietf.org/html/rfc6901#section-5
func isValidJSONPointer(jsonPointer string) error {
	// if !validateEscapeChars(jsonPointer) {
	// 	return fmt.Errorf("json pointer includes unescaped characters")
	// }
	// if _, err := jsonpointer.Parse(jsonPointer); err != nil {
	// 	return fmt.Errorf("invalid json pointer: %s", err.Error())
	// }
	if len(jsonPointer) == 0 {
		return nil
	}
	if jsonPointer[0] != '/' {
		return fmt.Errorf("non-empty references must begin with a '/' character")
	}
	str := jsonPointer[1:]
	if unescaptedTildaPattern.MatchString(str) {
		return fmt.Errorf("unescaped tilda error")
	}
	if endingTildaPattern.MatchString(str) {
		return fmt.Errorf("unescaped tilda error")
	}
	return nil
}

// A string instance is a valid against "regex" if it is a valid
// regular expression according to the ECMA 262 [ecma262] regular
// expression dialect. Implementations that validate formats MUST
// accept at least the subset of ECMA 262 defined in the Regular
// Expressions [regexInterop] section of this specification, and
// SHOULD accept all valid ECMA 262 expressions.
// http://www.ecma-international.org/publications/files/ECMA-ST/Ecma-262.pdf
// http://json-schema.org/latest/json-schema-validation.html#regexInterop
// https://tools.ietf.org/html/rfc7159
func isValidRegex(regex string) error {
	return nil
}

// A string instance is a valid against "relative-json-pointer" if it
// is a valid Relative JSON Pointer [relative-json-pointer].
// https://tools.ietf.org/html/draft-handrews-relative-json-pointer-00
func isValidRelJSONPointer(relJSONPointer string) error {
	parts := strings.Split(relJSONPointer, "/")
	if len(parts) == 1 {
		parts = strings.Split(relJSONPointer, "#")
	}
	if i, err := strconv.Atoi(parts[0]); err != nil || i < 0 {
		return fmt.Errorf("RJP must begin with positive integer")
	}
	//skip over first part
	str := relJSONPointer[len(parts[0]):]
	if len(str) > 0 && str[0] == '#' {
		return nil
	}
	return isValidJSONPointer(str)
}

// A string instance is valid against "time" if it is a valid
// representation according to the "full-time" production derived
// from RFC 3339, section 5.6 [RFC3339]
// https://tools.ietf.org/html/rfc3339#section-5.6
func isValidTime(time string) error {
	arbitraryDate := "1963-06-19"
	dateTime := fmt.Sprintf("%sT%s", arbitraryDate, time)
	return isValidDateTime(dateTime)
	return nil
}

// A string instance is a valid against "uri-reference" if it is a
// valid URI Reference (either a URI or a relative-reference),
// according to [RFC3986].
// https://tools.ietf.org/html/rfc3986
func isValidURIRef(uriRef string) error {
	return nil
}

// A string instance is a valid against "uri-template" if it is a
// valid URI Template (of any level), according to [RFC6570]. Note
// that URI Templates may be used for IRIs; there is no separate IRI
// Template specification.
// https://tools.ietf.org/html/rfc6570
func isValidURITemplate(uriTemplate string) error {
	return nil
}

// A string instance is a valid against "uri" if it is a valid URI,
// according to [RFC3986].
// https://tools.ietf.org/html/rfc3986
func isValidURI(uri string) error {
	if _, err := url.Parse(uri); err != nil {
		return fmt.Errorf("uri incorrectly formatted: %s", err.Error())
	}
	return nil
}
