package unsplash

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/pterm/pterm"
)

// ErrRequestFailed indicates a general error in service request.
var ErrRequestFailed = errors.New("request failed")

// nolint: gosec
// const token = "4c483af1b27cf8d55fc29504bc48e3755e47eb7a3dd3a320e92b23fc4e5aa1b8"
const token = "xteUYimSedgjkBsS6mJ8hG79rZihf7eitsJy8mBJL6w"

// Source is source implmentation for unsplash image service.
type Source struct {
	response    []Image
	N           int
	Query       string
	Orientation string
}

// Init initiates source and return number of available images.
func (s *Source) Init() (int, error) {
	resp, err := resty.New().
		SetBaseURL("https://api.unsplash.com").
		SetHeader("Accept-Version", "v1").
		SetHeader("Authorization", fmt.Sprintf("Client-ID %s", token)).
		R().
		SetResult(&s.response).
		SetQueryParam("count", strconv.Itoa(s.N)).
		SetQueryParam("orientation", s.Orientation).
		SetQueryParam("query", s.Query).
		Get("/photos/random")
	if err != nil {
		return 0, fmt.Errorf("network failure: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return 0, ErrRequestFailed
	}

	return len(s.response), nil
}

// Name returns source name.
func (s *Source) Name() string {
	return "unsplash"
}

// Fetch fetches given index from source.
func (s *Source) Fetch(index int) (string, io.ReadCloser, error) {
	image := s.response[index]

	pterm.Info.Printf("Getting %s (%s)\n", image.ID, image.Description)

	resp, err := resty.New().R().SetDoNotParseResponse(true).Get(image.URLs.Full)
	if err != nil {
		return "", nil, fmt.Errorf("network failure: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", nil, ErrRequestFailed
	}

	pterm.Success.Printf("%s was gotten\n", image.ID)

	return fmt.Sprintf("%s.jpg", image.ID), resp.RawBody(), nil
}
