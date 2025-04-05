package opkl

type tokenOption struct {
	token string
}

// apply implements option.
func (t *tokenOption) apply(opkl *opklReader) {
	opkl.opToken = ""
}

// Configures the 1Password token used when reading a secret
func WithToken(token string) option {
	return &tokenOption{
		token: token,
	}
}
