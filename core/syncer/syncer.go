package syncer

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/redesblock/uploader/core/util"
	log "github.com/sirupsen/logrus"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Syncer interface {
	Upload(host string, voucherID string, path string, indexExt string) (string, error)
}

func New(ignoreHidden bool) Syncer {
	return &syncer{
		ignoreHidden: ignoreHidden,
	}
}

const (
	PinHeader            = "Swarm-Pin"
	TagHeader            = "Swarm-Tag"
	EncryptHeader        = "Swarm-Encrypt"
	IndexDocumentHeader  = "Swarm-Index-Document"
	ErrorDocumentHeader  = "Swarm-Error-Document"
	FeedIndexHeader      = "Swarm-Feed-Index"
	FeedIndexNextHeader  = "Swarm-Feed-Index-Next"
	CollectionHeader     = "Swarm-Collection"
	PostageBatchIdHeader = "Swarm-Postage-Batch-Id"
	DeferredUploadHeader = "Swarm-Deferred-Upload"

	ContentTypeHeader = "Content-Type"
	MultiPartFormData = "multipart/form-data"
	ContentTypeTar    = "application/x-tar"
)

type syncer struct {
	ignoreHidden bool // ignore hidden files or not.
}

func (s *syncer) Upload(node string, voucherID string, path string, indexExt string) (string, error) {
	log.WithField("path", path).Debugf("uploading ...")

	t := time.Now()
	buf, filename, err := tarFile(path, indexExt, s.ignoreHidden)
	if err != nil {
		return "", fmt.Errorf("tar file error %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://"+node+":1633"+"/hop", buf)
	if err != nil {
		return "", fmt.Errorf("new request error %v", err)
	}
	req.Header.Add(DeferredUploadHeader, "true")
	req.Header.Add(PostageBatchIdHeader, voucherID)
	req.Header.Add(CollectionHeader, "true")
	req.Header.Add(ContentTypeHeader, ContentTypeTar)
	if len(filename) > 0 {
		req.Header.Add(IndexDocumentHeader, filename)
	}
	req.Header.Add(ErrorDocumentHeader, "")
	req.Header.Add(EncryptHeader, "false")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request error %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("resp status error %s, want %s", resp.Status, http.StatusText(http.StatusCreated))
	}
	var ret map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return "", fmt.Errorf("resp body error %v", err)
	}
	hash := ret["reference"]

	log.WithField("path", path).WithField("elapsed", time.Now().Sub(t)).Debug("upload success")
	return hash, nil
}

func tarFile(path string, indexExt string, ignoreHidden bool) (*bytes.Buffer, string, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, "", fmt.Errorf("load file info error %v", err)
	}
	basePath := path
	if !fi.IsDir() {
		basePath = filepath.Dir(path)
	}

	var buf bytes.Buffer
	var indexFile string
	writer := tar.NewWriter(&buf)
	if err := filepath.Walk(path, func(walkPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk error %v", err)
		}
		isHidden, err := util.IsHiddenFile(walkPath)
		if err != nil {
			return fmt.Errorf("isHidden error %v", err)
		}
		if ignoreHidden && isHidden {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if info.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(walkPath)
		if err != nil {
			return fmt.Errorf("read data error %v", err)
		}
		fileName, err := filepath.Rel(basePath, walkPath)
		if err != nil {
			return fmt.Errorf("relative path error %v", err)
		}
		if len(indexExt) > 0 && filepath.Dir(fileName) == "." && strings.HasSuffix(fileName, indexExt) {
			indexFile = fileName
		}
		hdr := &tar.Header{
			Name: fileName,
			Mode: 0600,
			Size: info.Size(),
		}
		if err := writer.WriteHeader(hdr); err != nil {
			return fmt.Errorf("write header error %v", err)
		}
		if _, err := writer.Write(data); err != nil {
			return fmt.Errorf("write data error %v", err)
		}
		log.WithField("path", path).WithField("file", fileName).Debugf("add to tar file")
		return nil
	}); err != nil {
		return nil, "", err
	}
	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("tar writer close error %v", err)
	}
	if len(indexExt) > 0 && len(indexFile) == 0 {
		return nil, "", fmt.Errorf("not found index")
	}
	return &buf, indexFile, nil
}
