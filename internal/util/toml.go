// Copyright 2026 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

// These code is copied from https://github.com/pelletier/go-toml/blob/v2/marshaler.go
const literalQuote = '\''

func EncodeTomlKey(k string) []byte {
	needsQuotation := false
	cannotUseLiteral := false

	var b []byte
	if len(k) == 0 {
		return append(b, "''"...)
	}

	for _, c := range k {
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '_' {
			continue
		}

		if c == literalQuote {
			cannotUseLiteral = true
		}

		needsQuotation = true
	}

	if needsQuotation && needsQuoting(k) {
		cannotUseLiteral = true
	}

	switch {
	case cannotUseLiteral:
		return encodeQuotedString(b, k)
	case needsQuotation:
		return encodeLiteralString(b, k)
	default:
		return encodeUnquotedKey(b, k)
	}

}

func encodeLiteralString(b []byte, v string) []byte {
	b = append(b, literalQuote)
	b = append(b, v...)
	b = append(b, literalQuote)

	return b
}

func encodeUnquotedKey(b []byte, v string) []byte {
	return append(b, v...)
}

func encodeQuotedString(b []byte, v string) []byte {
	stringQuote := `"`

	b = append(b, stringQuote...)

	const (
		hextable = "0123456789ABCDEF"
		// U+0000 to U+0008, U+000A to U+001F, U+007F
		nul = 0x0
		bs  = 0x8
		lf  = 0xa
		us  = 0x1f
		del = 0x7f
	)

	for _, r := range []byte(v) {
		switch r {
		case '\\':
			b = append(b, `\`...)
		case '"':
			b = append(b, `\"`...)
		case '\b':
			b = append(b, `\b`...)
		case '\f':
			b = append(b, `\f`...)
		case '\n':
			b = append(b, `\n`...)
		case '\r':
			b = append(b, `\r`...)
		case '\t':
			b = append(b, `\t`...)
		default:
			switch {
			case r >= nul && r <= bs, r >= lf && r <= us, r == del:
				b = append(b, `\u00`...)
				b = append(b, hextable[r>>4])
				b = append(b, hextable[r&0x0f])
			default:
				b = append(b, r)
			}
		}
	}

	b = append(b, stringQuote...)

	return b
}

func needsQuoting(v string) bool {
	for _, b := range []byte(v) {
		if b == '\'' || b == '\r' || b == '\n' || invalidAscii(b) {
			return true
		}
	}
	return false
}

var invalidAsciiTable = [256]bool{
	0x00: true,
	0x01: true,
	0x02: true,
	0x03: true,
	0x04: true,
	0x05: true,
	0x06: true,
	0x07: true,
	0x08: true,
	// 0x09 TAB
	// 0x0A LF
	0x0B: true,
	0x0C: true,
	// 0x0D CR
	0x0E: true,
	0x0F: true,
	0x10: true,
	0x11: true,
	0x12: true,
	0x13: true,
	0x14: true,
	0x15: true,
	0x16: true,
	0x17: true,
	0x18: true,
	0x19: true,
	0x1A: true,
	0x1B: true,
	0x1C: true,
	0x1D: true,
	0x1E: true,
	0x1F: true,
	// 0x20 - 0x7E Printable ASCII characters
	0x7F: true,
}

func invalidAscii(b byte) bool {
	return invalidAsciiTable[b]
}
