package auth

// Profile is the slice of a customer record the login flow needs. It is owned
// by auth so the OTP handlers can resolve an identity, enforce the
// email-required rule, and echo the canonical email back to the client without
// importing the customers module. The customers module fills it via SetLookup.
type Profile struct {
	FRID          string
	District      string
	Name          string
	Email         string // canonical delivery address ("" when none on file)
	Phone         string
	AccountNumber string
}

// Lookup resolves a login identity (email, phone, or account number) to a
// profile. ok=false means unknown identity.
type Lookup func(identity string) (*Profile, bool)

var lookupFn Lookup

// SetLookup wires the identity resolver (typically the customers module).
func SetLookup(f Lookup) { lookupFn = f }

// resolve returns the profile for an identity, or nil when unknown/unwired.
func resolve(identity string) *Profile {
	if lookupFn == nil {
		return nil
	}
	p, ok := lookupFn(identity)
	if !ok {
		return nil
	}
	return p
}
