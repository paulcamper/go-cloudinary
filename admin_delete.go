package cloudinary

import (
	"fmt"
	"net/http"
	"net/url"
)

const (
	// maxDeleteItemsPerRequest shows how many public_ids can be passed at a time to
	// a bulk delete method. Cloudinary doc says that the maximum is 100, but if public_ids
	// are long, you'll get 414(URI is too long) from the API.
	maxDeleteItemsPerRequest = 20

	pathUploadedImages = "/resources/image/upload"
	pathUploadedRaws   = "/resources/raw/upload"
)

func (s *Service) bulkDeleteUploadedWithInvalidation(resourceType ResourceType, publicIDs []string) error {
	if len(publicIDs) == 0 {
		return nil
	}

	path := pathUploadedImages
	if resourceType == RawType {
		path = pathUploadedRaws
	}

	var (
		currentDelimiter  = 0
		previousDelimiter = 0
		finish            = false
	)
	for !finish {
		// we cannot give all public ids as the input for the query, thus have
		// to make requests with subsets of public_ids.
		// The code take a new bunch of ids and guarantee that there is no out of bounds access.
		currentDelimiter += maxDeleteItemsPerRequest
		if currentDelimiter > len(publicIDs) {
			currentDelimiter = len(publicIDs)
			finish = true // come to the end of the publicIDs list, shouldn't loop again.
		}
		idsForQuery := publicIDs[previousDelimiter:currentDelimiter]
		previousDelimiter = currentDelimiter
		qs := url.Values{
			"invalidate":   []string{"true"},
			"public_ids[]": idsForQuery, // "[]" here is by design.
		}

		for {
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s%s?%s", s.adminURI, path, qs.Encode()), nil)
			if err != nil {
				return err
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			m, err := handleHttpResponse(resp)
			if err != nil {
				fmt.Println("idsForQuery", idsForQuery)
				return err
			}
			if e, ok := m["next_cursor"]; ok {
				qs.Set("next_cursor", e.(string))
			} else {
				break
			}
		}
	}

	return nil
}

func (s *Service) DeleteUploadedImages(publicIDs []string) error {
	return s.bulkDeleteUploadedWithInvalidation(ImageType, publicIDs)
}

func (s *Service) DeleteUploadedRawFiles(publicIDs []string) error {
	return s.bulkDeleteUploadedWithInvalidation(RawType, publicIDs)
}
