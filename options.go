package opkl

import (
	"context"
	"fmt"

	"github.com/1password/onepassword-sdk-go"
)

// Configures the 1Password token used when reading a secret
func WithToken(token string) option {
	return func(opkl *opklReader) error {
		opkl.opToken = token

		return nil
	}
}

func WithClient(client Client) option {
	return func(opkl *opklReader) error {
		opkl.opClient = client

		return nil
	}
}

type opklClient struct {
	opClient *onepassword.Client
}

func (client *opklClient) ReadSecret(ctx context.Context, ref string) (string, error) {
	secret, err := client.opClient.Secrets().Resolve(ctx, ref)
	if err != nil {
		return "", err
	}

	return secret, nil
}

func WithClientDefault(token string) option {
	return func(opkl *opklReader) error {
		client, err := onepassword.NewClient(
			context.Background(),
			onepassword.WithServiceAccountToken(token),
			onepassword.WithIntegrationInfo(
				onepassword.DefaultIntegrationName,
				onepassword.DefaultIntegrationVersion,
			),
		)
		if err != nil {
			return fmt.Errorf("error initializing 1password client: %w", err)
		}

		opkl.opClient = &opklClient{
			opClient: client,
		}

		return nil
	}
}
