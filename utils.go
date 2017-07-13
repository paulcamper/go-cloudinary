package cloudinary

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/pkg/errors"
)

// TODO(oleh): add custom http client, as default can hang out on timeout forever.
func (s *Service) makeRequest(url string, body *APIBody) (*uploadResponse, error) {
	req, err := http.NewRequest(http.MethodPost, url, body.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create a new http request")
	}
	req.Header.Set("Content-Type", body.ContentType)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to do a request")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errInfo struct {
			Error struct {
				Message string
			}
		}
		if err := json.NewDecoder(resp.Body).Decode(&errInfo); err != nil {
			return nil, errors.Wrap(err, "Failed to decode with json")
		}
		return nil, errors.New(errInfo.Error.Message)
	}
	var upInfo uploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&upInfo); err != nil {
		return nil, errors.Wrap(err, "Failed to decode with json")
	}
	return &upInfo, nil
}

func (s *Service) signature(params params) (string, error) {
	hash := sha1.New()
	_, err := hash.Write([]byte(params.stringForSignature() + s.apiSecret))
	if err != nil {
		return "", errors.Wrap(err, "Failed to write to hash writer")
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

type APIBody struct {
	Bytes       *bytes.Buffer
	ContentType string
}

func apiBody(params params, r io.Reader) (*APIBody, error) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	for pName, pValue := range params {
		ws, err := w.CreateFormField(pName)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create form field")
		}
		_, err = ws.Write([]byte(pValue))
		if err != nil {
			return nil, errors.Wrap(err, "Failed to write field to writer")
		}
	}

	publicID, found := params.publicID()
	if found && r != nil {
		fw, err := w.CreateFormFile("file", publicID)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create form file")
		}
		tmp, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to read file from reader")
		}
		_, err = fw.Write(tmp)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to write file data to writer")
		}

	}

	err := w.Close()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to close multipart writer")
	}
	return &APIBody{
		Bytes:       buf,
		ContentType: w.FormDataContentType(),
	}, nil
}

func (s *Service) apiURL(baseURL string, resourceType ResourceType, action Action) string {
	return fmt.Sprintf("%s/%s/%s/%s", baseURL, s.cloudName, resourceType, action)
}
