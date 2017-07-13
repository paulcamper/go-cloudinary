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
func (s *Service) makeRequest(url string, body *APIBody) ([]byte, error) {
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
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed read bytes from resp body")
	}
	return respBytes, nil
}

func uploadInfoFromBytes(bodyBytes []byte) (*uploadResponse, error) {
	var upInfo uploadResponse
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&upInfo); err != nil {
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

func apiBody(params params, file interface{}) (*APIBody, error) {
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

	switch f := file.(type) {
	case io.Reader:
		publicID, found := params.publicID()
		if found {
			fw, err := w.CreateFormFile("file", publicID)
			if err != nil {
				return nil, errors.Wrap(err, "Failed to create form file")
			}
			tmp, err := ioutil.ReadAll(f)
			if err != nil {
				return nil, errors.Wrap(err, "Failed to read file from reader")
			}
			_, err = fw.Write(tmp)
			if err != nil {
				return nil, errors.Wrap(err, "Failed to write file data to writer")
			}

		}
	case string:
		// TODO: make dependency on link or filepath. Currently, link only applies
		ws, err := w.CreateFormField("file")
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create form field")
		}
		_, err = ws.Write([]byte(f))
		if err != nil {
			return nil, errors.Wrap(err, "Failed to write file string")
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

type ToParamser interface {
	ToParams() params
}

func (s *Service) paramsForAPICall(p ToParamser) (params, error) {
	params := s.basicParams()
	params.join(p.ToParams())

	signature, err := s.signature(params)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create a signature")
	}
	params.set("signature", signature)
	return params, nil
}