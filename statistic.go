package cloudinary

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ResourceStatistic struct {
	Usage        float64
	CreditsUsage float64 `json:"credits_usage"`
}

type CreditsStatistic struct {
	Usage       float64
	Limit       float64
	UsedPercent float64 `json:"used_percent"`
}

type Statistic struct {
	Transformations ResourceStatistic
	Objects         ResourceStatistic
	Bandwidth       ResourceStatistic
	Storage         ResourceStatistic

	Credits CreditsStatistic

	Requests         int64
	Resources        int64
	DerivedResources int64 `json:"derived_resources"`
}

const statisticPath = "/usage"

func (s *Service) Statistic() (Statistic, error) {
	statisticURL := fmt.Sprintf("%s%s", s.adminURI, statisticPath)

	req, err := http.NewRequest(http.MethodGet, statisticURL, nil)
	if err != nil {
		return Statistic{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Statistic{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Statistic{}, err
	}
	var result Statistic
	if err := json.Unmarshal(body, &result); err != nil {
		return Statistic{}, err
	}

	return result, nil
}
