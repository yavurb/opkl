package opkl

import "github.com/1password/onepassword-sdk-go"

// Configures the 1Password token used when reading a secret
func WithToken(token string) option {
	return func(opkl *opklReader) error {
		opkl.opToken = token

		return nil
	}
}

func WithClient(client *onepassword.Client) option {
	return func(opkl *opklReader) error {
		opkl.opClient = client

		return nil
	}
}
