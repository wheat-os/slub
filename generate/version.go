package generate

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	Version = "v1.0.4-alpha"
)

type gitTag []struct {
	Name       string `json:"name"`
	ZipballURL string `json:"zipball_url"`
	TarballURL string `json:"tarball_url"`
	Commit     struct {
		Sha string `json:"sha"`
		URL string `json:"url"`
	} `json:"commit"`
	NodeID string `json:"node_id"`
}

func GetSlubbyVersion() (gitTag, error) {
	url := "https://api.github.com/repos/wheat-os/slubby/tags"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get tag err status: %d", resp.StatusCode)
	}

	tags := make(gitTag, 0)
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(&tags); err != nil {
		return nil, err
	}

	return tags, nil
}
