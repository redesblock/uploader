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
	fmt.Println(New(true).Upload("183.131.181.164", "dd812517f2ecfe75d7b08e908a857c8703477949770fbb76f2244d0d90cb4a12", "syncer.go", ""))
}
