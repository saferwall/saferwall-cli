// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package webapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/saferwall/saferwall-cli/internal/entity"
)

func (s Service) newfileUploadRequest(fieldname, filename string, params map[string]string) (*http.Request, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	for key, val := range params {
		err := writer.WriteField(key, val)
		if err != nil {
			return nil, err
		}
	}

	part, err := writer.CreateFormFile(fieldname, filepath.Base(filename))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.filesURL, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

// FileExists determines file existence.
// TODO: use HEAD instead.
func (s Service) FileExists(sha256 string) (bool, error) {

	url := s.filesURL + sha256
	resp, err := http.Head(url)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	defer resp.Body.Close()
	return true, nil
}

// ListFiles list all the files in DB.
func (s Service) ListFiles(authToken string, page int) (*Pages, error) {

	var pages Pages
	url := fmt.Sprintf("%s?per_page=%d&page=%d&fields=sha256", s.filesURL, 1000, page)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Cookie", "JWTCookie="+authToken)

	// Perform the http post request.
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	// Read the response.
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var jsonBody map[string]interface{}
		err = json.Unmarshal(body.Bytes(), &jsonBody)
		if err != nil {
			return nil, err
		}
		msg := jsonBody["message"].(string)
		return nil, errors.New(msg)
	}

	err = json.Unmarshal(body.Bytes(), &pages)
	if err != nil {
		return nil, err
	}

	return &pages, nil

}

func (s Service) Scan(filepath string, authToken, preferredOS string, skipDetonation bool, timeout int) (string, error) {
	params := map[string]string{
		"skip_detonation": strconv.FormatBool(skipDetonation),
		"os":              preferredOS,
		"timeout":         strconv.Itoa(timeout),
	}

	// Create a new file upload request.
	request, err := s.newfileUploadRequest("file", filepath, params)
	if err != nil {
		return "", err
	}

	// Add our auth token.
	request.Header.Set("Cookie", "JWTCookie="+authToken)

	// Perform the http post request.
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	// Read the response.
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	return body.String(), nil
}

func (s Service) Rescan(sha256, authToken, preferredOS string, skipDetonation bool, timeout int) error {

	url := s.filesURL + sha256 + "/rescan"

	requestBody, err := json.Marshal(map[string]interface{}{
		"skip_detonation": skipDetonation,
		"scan_cfg": map[string]interface{}{
			"os":      preferredOS,
			"timeout": timeout,
		},
	})
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(requestBody)
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", "JWTCookie="+authToken)

	// Perform the http post request.
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	// Read the response.
	body = &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	resp.Body.Close()
	fmt.Println(body)
	return nil
}

// GetFile retrieves the file report given a sha256.
func (s Service) GetFile(sha256 string, file *entity.File) error {

	url := s.filesURL + sha256
	client := &http.Client{}
	client.Timeout = time.Second * 10

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(d, &file)
}

func (s Service) Download(sha256, authToken string) (*bytes.Buffer, error) {

	url := s.filesURL + sha256 + "/download"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", "JWTCookie="+authToken)

	// Perform the http post request.
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	// Read the response.
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return body, nil
}

func (s Service) Delete(sha256, authToken string) error {

	url := s.filesURL + sha256
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cookie", "JWTCookie="+authToken)

	// Perform the http post request.
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	// Read the response.
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}
