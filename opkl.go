package opkl

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/1password/onepassword-sdk-go"
	"github.com/apple/pkl-go/pkl"
)

type opklReader struct {
	opToken  string
	opClient *onepassword.Client
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

	secret, err := o.opClient.Secrets().Resolve(context.Background(), opReference)
	if err != nil {
		return nil, fmt.Errorf("error resolving secret reference %s: %w", opReference, err)
	}

	return []byte(secret), nil
}

// Scheme implements pkl.ResourceReader.
func (o *opklReader) Scheme() string {
	return "op"
}

type option interface {
	apply(opkl *opklReader)
}

func New(options ...option) (pkl.ResourceReader, error) {
	opkl := &opklReader{}

	for _, option := range options {
		option.apply(opkl)
	}

	if opkl.opToken == "" {
		token, ok := os.LookupEnv("OP_SERVICE_ACCOUNT_TOKEN")
		if !ok {
			return nil, errors.New("1Password service account token not defined. Set the env OP_SERVICE_ACCOUNT_TOKEN")
		}

		opkl.opToken = token
	}

	client, err := onepassword.NewClient(
		context.Background(),
		onepassword.WithServiceAccountToken(opkl.opToken),
		onepassword.WithIntegrationInfo(
			onepassword.DefaultIntegrationName,
			onepassword.DefaultIntegrationVersion,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error initializing 1password client: %w", err)
	}

	opkl.opClient = client

	return opkl, nil
}
