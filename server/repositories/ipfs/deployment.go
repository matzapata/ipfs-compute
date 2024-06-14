package repositories

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type IpfsDeploymentsRepository struct {
}

func NewIpfsDeploymentsRepository() *IpfsDeploymentsRepository {
	return &IpfsDeploymentsRepository{}
}

func (dr *IpfsDeploymentsRepository) GetZippedDeployment(cid string, maxSize uint) ([]byte, error) {
	ipfsUrl := fmt.Sprintf("https://gateway.pinata.cloud/ipfs/%s", cid)

	// first, check the file size is within the limit
	size, err := getFileSize(ipfsUrl)
	if err != nil {
		return nil, err
	}
	if size > int64(maxSize) {
		return nil, fmt.Errorf("file size exceeds the limit of %d bytes", maxSize)
	}

	// download the file
	data, err := downloadFile(ipfsUrl, size)
	if err != nil {
		return nil, err
	}

	return data, nil
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

func downloadFile(url string, maxSize int64) ([]byte, error) {
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
