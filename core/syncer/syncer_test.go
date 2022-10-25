package syncer

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestTarFile(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	_, index, err := tarFile("../../core/model/DB.go", "", true)
	fmt.Println(index, err)
}

func TestUpload(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	fmt.Println(New(true).Upload("http://58.34.1.130:1633", "e92110b77f959065768e24a44c5ab04de4f6bc20f0010fbba726ee4b31291797", "syncer.go", ""))
}
