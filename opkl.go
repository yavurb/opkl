package opkl

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/apple/pkl-go/pkl"
)

type opklReader struct {
	opToken  string
	opClient Client
}

// HasHierarchicalUris implements pkl.ResourceReader.
func (o *opklReader) HasHierarchicalUris() bool {
	return false
}

// IsGlobbable implements pkl.ResourceReader.
func (o *opklReader) IsGlobbable() bool {
	return false
}

// ListElements implements pkl.ResourceReader.
func (o *opklReader) ListElements(url url.URL) ([]pkl.PathElement, error) {
	return nil, nil
}

// Read implements pkl.ResourceReader.
func (o *opklReader) Read(resourceUrl url.URL) ([]byte, error) {
	strReference := resourceUrl.String()

	if b, err := base64.StdEncoding.DecodeString(strReference); err == nil {
		strReference = string(b)
	}

	opReference, err := url.QueryUnescape(strReference)
	if err != nil {
		return nil, fmt.Errorf("error unescaping secret reference %s: %w", opReference, err)
	}

	// secret, err := o.opClient.Secrets().Resolve(context.Background(), opReference)
	secret, err := o.opClient.ReadSecret(context.Background(), opReference)
	if err != nil {
		return nil, fmt.Errorf("error resolving secret reference %s: %w", opReference, err)
	}

	return []byte(secret), nil
}

// Scheme implements pkl.ResourceReader.
func (o *opklReader) Scheme() string {
	return "op"
}

func New(options ...option) (pkl.ResourceReader, error) {
	opkl := &opklReader{}

	for _, option := range options {
		err := option(opkl)
		if err != nil {
			return nil, fmt.Errorf("error applying option: %w", err)
		}
	}

	return opkl, nil
}
