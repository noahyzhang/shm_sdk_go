package shm_sdk_go

var arrCRCTable = [...]uint32{
	0x0000, 0x1021, 0x2042, 0x3063, 0x4084, 0x50a5, 0x60c6, 0x70e7,
	0x8108, 0x9129, 0xa14a, 0xb16b, 0xc18c, 0xd1ad, 0xe1ce, 0xf1ef}

func CalcCRCVal(pBuf []uint8, length uint32) (uint32, error) {
	var uc uint8 = 0
	var crc uint32 = 0
	var i uint32 = 0
	for ; i < length; i++ {
		uc = uint8(crc/256) / 16
		crc <<= 4
		crc ^= arrCRCTable[uc^(pBuf[i]/16)]

		uc = uint8(crc/256) / 16
		crc <<= 4
		crc ^= arrCRCTable[uc^(pBuf[i]&0x0f)]
	}
	return crc, nil
}
