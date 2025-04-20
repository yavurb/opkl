package opkl

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strings"
	"testing"
)

type clientMock struct {
	wantSecret    []byte
	wantError     error
	wantReference string
}

func (c *clientMock) ReadSecret(ctx context.Context, ref string) (string, error) {
	if c.wantReference != "" {
		if c.wantReference != ref {
			xd := strings.Compare(c.wantReference, ref)
			fmt.Println(xd)
			fmt.Println(c.wantReference)
			fmt.Println(ref)
			return "", fmt.Errorf("Reference %s not found.", ref)
		}
	}

	return string(c.wantSecret), c.wantError
}

func Test_opklReader_Read(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		resourceUrl string
		want        []byte
		wantErr     bool
		err         error
		wantRef     string
	}{
		{
			name:        "should read a secret",
			resourceUrl: "op://test-vault/test-item/secret",
			want:        []byte("My Secret"),
			wantErr:     false,
		},
		{
			name:        "should read a secret from a base64 encoded reference",
			resourceUrl: "op:Ly90ZXN0LXZhdWx0L3Rlc3QgaXRlbS9zZWNyZXQK",
			want:        []byte("My Secret"),
			wantErr:     false,
			wantRef:     "op://test-vault/test item/secret",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &clientMock{
				wantSecret:    tt.want,
				wantReference: tt.wantRef,
				wantError:     tt.err,
			}

			reader, err := New(WithClient(client))
			if err != nil {
				t.Errorf("expected no error creating reader, got: %s", err)
			}

			u, err := url.Parse(tt.resourceUrl)
			if err != nil {
				t.Errorf("expected no error creating test URL, got: %s", err)
			}

			got, gotErr := reader.Read(*u)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Read() failed: %v", gotErr)
				}
				return
			}

			if tt.wantErr {
				t.Fatal("Read() succeeded unexpectedly")
			}

			if !bytes.Equal(got, tt.want) {
				t.Errorf("Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
