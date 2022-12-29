package parsing

func ParseIconString(s string) string {
	var result string

	for bi := 0; bi < len(s); bi++ {
		b := s[bi]

		if b == '\\' {
			// escape code
			c := s[bi+1]
			switch c {
			case 'a':
				// CSI
				result += "\033["
				bi++
				continue
			case 'u':
				// Parse unicode character
				encoding, elen := ParseUnicode(s[bi+2:])
				result += encoding
				bi += elen + 1
				continue
			}
		}

		result += string(b)
	}

	return result
}

func ParseUnicode(s string) (string, int) {
	var result []byte

	var parsingLength int = 0
	var codeSegments []int
	for bi := 0; bi < len(s); bi++ {
		b := s[bi]

		if b >= '0' && b <= '9' {
			codeSegments = append(codeSegments, int(b-'0'))
		} else if b >= 'a' && b <= 'f' {
			codeSegments = append(codeSegments, int(b-'a'+10))
		} else if b >= 'A' && b <= 'F' {
			codeSegments = append(codeSegments, int(b-'A'+10))
		} else {
			break
		}

		parsingLength++
	}

	var code int32
	for si, seg := range codeSegments {
		multiplier := len(codeSegments) - si - 1
		addition := seg
		for mi := 0; mi < multiplier; mi++ {
			addition *= 16
		}
		code += int32(addition)
	}

	var b1, b2, b3, b4 byte = 0, 0x80, 0x80, 0x80
	if code >= 0 && code <= 0x7f {
		result = append(result, byte(code)&0x7f)
	} else if code >= 0x80 && code <= 0x7ff {
		b1 = 0xc0
		b1 |= byte(code >> 6)
		b2 |= byte(code & 0x3f)
		result = append(result, b1, b2)
	} else if code >= 0x800 && code <= 0xffff {
		b1 = 0xe0
		b1 |= byte(code >> 12)
		b2 |= byte((code & 0xfc0) >> 6)
		b3 |= byte(code & 0x3f)
		result = append(result, b1, b2, b3)
	} else {
		b1 = 0xf0
		b1 |= byte(code >> 18)
		b2 |= byte((code & 0x3f000) >> 12)
		b3 |= byte((code & 0xfc0) >> 6)
		b4 |= byte(code & 0x3f)
		result = append(result, b1, b2, b3, b4)
	}

	return string(result), parsingLength
}
