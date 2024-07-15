package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/matzapata/ipfs-compute/provider/internal/domain"
	"github.com/wabarc/ipfs-pinner/pkg/pinata"
)

type IpfsArtifactRepository struct {
	Gateway string
	Pinata  pinata.Pinata
	Cache   IFileCache
}

func NewIpfsArtifactRepository(cache IFileCache, gateway string, pinataApiKey string, pinataSecret string) *IpfsArtifactRepository {
	return &IpfsArtifactRepository{
		Gateway: gateway,
		Pinata:  pinata.Pinata{Apikey: pinataApiKey, Secret: pinataSecret},
		Cache:   cache,
	}
}

func (dt *IpfsArtifactRepository) GetZippedExecutable(cid string, maxSize uint) ([]byte, error) {
	if dt.Cache != nil && dt.Cache.Exists(cid) {
		return dt.Cache.Get(cid)
	}

	res, err := downloadFile(dt.Gateway, cid, int64(maxSize))
	if err != nil {
		return nil, err
	}

	if dt.Cache != nil {
		dt.Cache.Set(cid, res)
	}
	return res, nil
}

func (dt *IpfsArtifactRepository) GetSpecificationFile(cid string) ([]byte, error) {
	if dt.Cache != nil && dt.Cache.Exists(cid) {
		return dt.Cache.Get(cid)
	}

	const maxSize = 1024 * 1024 // 1MB
	res, err := downloadFile(dt.Gateway, cid, maxSize)
	if err != nil {
		return nil, err
	}

	if dt.Cache != nil {
		dt.Cache.Set(cid, res)
	}
	return res, nil
}

func (dt *IpfsArtifactRepository) PublishArtifact(artPath string) (cid string, err error) {
	return dt.Pinata.PinFile(artPath)
}

func (dt *IpfsArtifactRepository) PublishArtifactSpecification(spec *domain.ArtifactSpec) (cid string, err error) {
	specFile, err := os.CreateTemp("", "artspec-*.json")
	if err != nil {
		return "", err
	}
	defer os.Remove(specFile.Name())

	specBytes, err := json.Marshal(spec)
	if err != nil {
		return "", err
	}

	// save and store it
	err = os.WriteFile(specFile.Name(), specBytes, 0644)
	if err != nil {
		return "", err
	}

	return dt.Pinata.PinFile(specFile.Name())
}

func downloadFile(gateway string, cid string, maxSize int64) ([]byte, error) {
	url := fmt.Sprintf("%v/ipfs/%s", gateway, cid)

	// first, check the file size is within the limit
	size, err := getFileSize(url)
	if err != nil {
		return nil, err
	}
	if size > int64(maxSize) {
		return nil, fmt.Errorf("file size exceeds the limit of %d bytes", maxSize)
	}

	// All good, download the file
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	// read the response body into a buffer
	var buf bytes.Buffer
	limitedReader := &io.LimitedReader{R: resp.Body, N: maxSize}
	_, err = io.Copy(&buf, limitedReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return buf.Bytes(), nil
}

func getFileSize(url string) (int64, error) {
	// send a HEAD request
	resp, err := http.Head(url)
	if err != nil {
		return 0, fmt.Errorf("failed to send HEAD request: %w", err)
	}
	defer resp.Body.Close()

	// check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bad status: %s", resp.Status)
	}

	// get the Content-Length header
	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		return 0, fmt.Errorf("Content-Length header is missing")
	}

	// parse the Content-Length value
	size, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid Content-Length value: %w", err)
	}

	return size, nil
}
