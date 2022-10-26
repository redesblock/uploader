package util

import (
	"fmt"
	"testing"
)

func TestReferenceUsabe(t *testing.T) {
	NodeUsabe("58.34.1.130")
	fmt.Println(VoucherUsabe("58.34.1.130", "e92110b77f959065768e24a44c5ab04de4f6bc20f0010fbba726ee4b31291797"))
	//ReferenceUsabe("https://gateway.mopweb3.com:13443", "206b444e7bbbfbf8d565823cf63094574ab55715b5d4b8b25b2f327fafe0bc8b")
}
