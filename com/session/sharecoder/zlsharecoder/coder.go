package zlsharecoder


func FixedEncode(b []byte, n uint64, bits int) int {
	if len(b) < (bits >> 3) {
		return 0
	}
	switch bits {
	case 8:
		b[0] = byte(n)
		return 1
	case 16:
		b[0] = byte(n >> 8)
		b[1] = byte(n)
		return 2
	case 32:
		b[0] = byte(n >> 24)
		b[1] = byte(n >> 16)
		b[2] = byte(n >> 8)
		b[3] = byte(n)
		return 4
	default:
		b[0] = byte(n >> 56)
		b[1] = byte(n >> 48)
		b[2] = byte(n >> 40)
		b[3] = byte(n >> 32)
		b[4] = byte(n >> 24)
		b[5] = byte(n >> 16)
		b[6] = byte(n >> 8)
		b[7] = byte(n)
		return 8
	}
}

func FixedDecode(b []byte, bits int) (n uint64, l int) {
	if len(b) < (bits >> 3) {
		l = 0
		return
	}
	switch bits {
	case 64:
		n = (uint64(b[0]) << 56) | (uint64(b[1]) << 48) | (uint64(b[2]) << 40) | (uint64(b[3]) << 32) |
			(uint64(b[4]) << 24) | (uint64(b[5]) << 16) | (uint64(b[6]) << 8) | (uint64(b[7]))
		l = 8
		return
	case 32:
		n = uint64((uint32(b[0]) << 24) | (uint32(b[1]) << 16) | (uint32(b[2]) << 8) | (uint32(b[3])))
		l = 4
		return
	case 16:
		n = uint64((uint16(b[0]) << 8) | (uint16(b[1])))
		l = 2
		return
	case 8:
		n = uint64(b[0])
	}
	return
}

func ShareEncode(b []byte, n uint64, bits int) int {
	if bits != 32 {
		return FixedEncode(b, n, bits)
	}
	if n < 0x80 {
		if len(b) < 1 {
			return 0
		}
		b[0] = byte(n)
		return 1
	}
	if n < 0x4000 {
		if len(b) < 2 {
			return 0
		}
		b[0] = byte(n>>8) | 0x80
		b[1] = byte(n)
		return 2
	}
	if n < 0x20000000 {
		if len(b) < 4 {
			return 0
		}
		b[0] = byte(n>>24) | 0xc0
		b[1] = byte(n >> 16)
		b[2] = byte(n >> 8)
		b[3] = byte(n)
		return 4
	}
	if len(b) < 5 {
		return 0
	}
	b[0] = 0xe0
	b[1] = byte(n >> 24)
	b[2] = byte(n >> 16)
	b[3] = byte(n >> 8)
	b[4] = byte(n)
	return 5
}

func ShareDecode(b []byte, bits int) (n uint64, l int) {
	if bits != 32 {
		return FixedDecode(b, bits)
	}
	if len(b) <= 0 {
		l = 0
		return
	}
	switch b[0] & 0xe0 {
	case 0xe0:
		if len(b) < 5 {
			l = 0
			return
		}
		n = uint64(b[4]) | uint64(b[3])<<8 | uint64(b[2])<<16 | uint64(b[1])<<24
		l = 5
		return
	case 0xc0:
		if len(b) < 4 {
			l = 0
			return
		}
		n = 0x1FFFFFFF & (uint64(b[3]) | uint64(b[2])<<8 | uint64(b[1])<<16 | uint64(b[0])<<24)
		l = 4
		return
	case 0xa0, 0x80:
		if len(b) < 2 {
			l = 0
			return
		}
		n = 0x3FFF & (uint64(b[1]) | uint64(b[0])<<8)
		l = 2
		return
	default:
	}
	n = uint64(b[0])
	l = 1
	return
}
