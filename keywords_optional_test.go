package jsonschema

import (
	"context"
	"testing"
)

func TestEmailFormat(t *testing.T) {
	ctx := context.Background()
	email := Must(`{"format":"email"}`)
	idnEmail := Must(`{"format":"idn-email"}`)

	cases := []struct {
		address  string
		email    bool
		idnEmail bool
	}{
		{`joe.bloggs@example.com`, true, true},
		{`JOE@EXAMPLE.COM`, true, true},
		{`user+tag@sub.example.co.uk`, true, true},
		{`te~st@example.com`, true, true},
		{`~test@example.com`, true, true},
		{`test~@example.com`, true, true},
		{`te.s.t@example.com`, true, true},
		{`a@b`, true, true},
		{`"joe bloggs"@example.com`, true, true},
		{`"joe..bloggs"@example.com`, true, true},
		{`"joe@bloggs"@example.com`, true, true},
		{`"joe\"bloggs"@example.com`, true, true},
		{`joe.bloggs@[127.0.0.1]`, true, true},
		{`joe.bloggs@[IPv6:::1]`, true, true},
		{`joe.bloggs@[ipv6:2001:db8::1]`, true, true},

		// characters beyond ASCII need idn-email
		{`실례@실례.테스트`, false, true},
		{`用户@例子.广告`, false, true},
		{`joe@exämple.com`, false, true},
		{`jöe@example.com`, false, true},
		{`"jöe bloggs"@example.com`, false, true},

		// not an address at all
		{``, false, false},
		{`2962`, false, false},
		{`ab.com`, false, false},
		{`@b.com`, false, false},
		{`a@`, false, false},
		{`a@b@c.com`, false, false},

		// a mailbox of a message header is more than an address
		{`Evil <a@b.com>`, false, false},
		{`<a@b.com>`, false, false},
		{`"<script>alert(1)</script>" <a@b.com>`, false, false},
		{`=?utf-8?q?x?= <a@b.com>`, false, false},
		{`a@b.com (comment <script>)`, false, false},
		{`(comment) a@b.com`, false, false},
		{`group: a@b.com;`, false, false},
		{`a@b.com, c@d.com`, false, false},

		// white space and line breaks
		{` a@b.com`, false, false},
		{`a@b.com `, false, false},
		{`a@ b.com`, false, false},
		{`a @b.com`, false, false},
		{"a@b.com\n", false, false},
		{"a@b.com\r\nBcc: x@y.z", false, false},

		// local part
		{`.test@example.com`, false, false},
		{`test.@example.com`, false, false},
		{`te..st@example.com`, false, false},
		{`a"b@c.com`, false, false},
		{`""@example.com`, false, false},
		{`"joe"bloggs"@example.com`, false, false},

		// domain
		{`joe.bloggs@invalid=domain.com`, false, false},
		{`a@-b.com`, false, false},
		{`a@b-.com`, false, false},
		{`a@b_c.com`, false, false},
		{`a@b..com`, false, false},
		{`a@.com`, false, false},
		{`a@b.com.`, false, false},
		{`joe.bloggs@[127.0.0.300]`, false, false},
		{`a@[]`, false, false},
		{`a@[not an ip]`, false, false},
		{`a@[::1]`, false, false},
		{`a@[IPv6:127.0.0.1]`, false, false},
		{`a@[127.0.0.1`, false, false},
	}

	for _, c := range cases {
		if got := email.Validate(ctx, c.address).IsValid(); got != c.email {
			t.Errorf("email %q: expected valid: %t, got: %t", c.address, c.email, got)
		}
		if got := idnEmail.Validate(ctx, c.address).IsValid(); got != c.idnEmail {
			t.Errorf("idn-email %q: expected valid: %t, got: %t", c.address, c.idnEmail, got)
		}
	}
}
