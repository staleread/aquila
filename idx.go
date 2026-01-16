package main

type idx uint16

func randIdx() idx {
	return idx(randByte()) | idx(randByte())<<8
}
