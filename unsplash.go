package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type photo struct {
	ID             string `json:"id"`
	Description    string `json:"description"`
	AltDescription string `json:"alt_description"`
	User           struct {
		Username string `json:"username"`
		Name     string `json:"name"`
	} `json:"user"`
	Links struct {
		HTML             string `json:"html"`
		DownloadLocation string `json:"download_location"`
	} `json:"links"`
	URLs struct {
		Full string `json:"full"`
	} `json:"urls"`
}

type searchPage struct {
	TotalPages int     `json:"total_pages"`
	Results    []photo `json:"results"`
}

type unsplash struct {
	accessKey string
	http      *http.Client
}

func newUnsplash(accessKey string) unsplash {
	return unsplash{accessKey: accessKey, http: &http.Client{Timeout: 60 * time.Second}}
}

func (u unsplash) search(query, orientation string, page int) (searchPage, error) {
	params := url.Values{
		"query":    {query},
		"per_page": {"30"},
		"page":     {strconv.Itoa(page)},
	}
	if orientation != "" {
		params.Set("orientation", orientation)
	}
	var result searchPage
	err := u.getJSON("https://api.unsplash.com/search/photos?"+params.Encode(), &result)
	return result, err
}

func (u unsplash) trackDownload(p photo) error {
	return u.getJSON(p.Links.DownloadLocation, &struct{}{})
}

func (u unsplash) getJSON(target string, out any) error {
	body, err := u.get(target, true)
	if err != nil {
		return err
	}
	defer body.Close()
	return json.NewDecoder(body).Decode(out)
}

func download(target, path string) error {
	body, err := unsplash{http: &http.Client{Timeout: 5 * time.Minute}}.get(target, false)
	if err != nil {
		return err
	}
	defer body.Close()
	tmp := path + ".part"
	file, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := io.Copy(file, body); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (u unsplash) get(target string, authenticated bool) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}
	if authenticated {
		req.Header.Set("Authorization", "Client-ID "+u.accessKey)
		req.Header.Set("Accept-Version", "v1")
	}
	resp, err := u.http.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("GET %s: %s", req.URL.Host+req.URL.Path, resp.Status)
	}
	return resp.Body, nil
}
