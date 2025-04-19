package opkl

// Configures the 1Password token used when reading a secret
func WithToken(token string) option {
	return func(opkl *opklReader) error {
		opkl.opToken = token

		return nil
	}
}
