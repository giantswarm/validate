// Package validate/web provides functions for validaing web data types such as
// URLs, domains, email addresses, etc.
package web

import (
	"bytes"
	"errors"
	"github.com/giantswarm/validate"
	"github.com/inhies/go-tld"
	"unicode"
	"unicode/utf8"
)

var (
	IANA = tld.IANA  // URL of the default TLD list.
	TLDs = &tld.TLDs // Pointer to the tld packages URL slice.
)

// A domain value to be validated
type Domain struct {
	domain  []byte
	message string
	checks  map[string]interface{}
}

// Return the domain
func (d *Domain) String() string {
	return string(d.domain)
}

// Create a new domain value to be validated
func NewDomain(domain string) *Domain {
	d := Domain{
		domain: []byte(domain),
		checks: make(map[string]interface{}),
	}
	return &d
}

// Sets the validation failure message.
func (d *Domain) SetMessage(msg string) validate.Method {
	d.message = msg
	return d
}

// Return the failed validation message.
func (d *Domain) Message() string {
	return d.message
}

// Specifies the minimum number of sub-domains to allow.
func (d *Domain) MinSubdomains(m int) *Domain {
	d.checks["minsubs"] = m
	return d
}

// Specifies the maximum number of sub-domains to allow.
func (d *Domain) MaxSubdomains(m int) *Domain {
	d.checks["maxsubs"] = m
	return d
}

// Specifies the maximum number of sub-domains to allow.
func (d *Domain) MaxLength(m int) *Domain {
	d.checks["maxlength"] = m
	return d
}

var (
	// A-Z, a-z, 0-9, and hyphen are the only valid characters for domains.
	domainTable = &unicode.RangeTable{
		R16: []unicode.Range16{
			{'-', '-', 1},
			{'0', '9', 1},
			{'A', 'Z', 1},
			{'a', 'z', 1},
		},
		LatinOffset: 4,
	}
)
var (
	// Entire domain or a single label has an invalid length.
	ErrDomainLength = &validate.ValidatorError{
		ErrLevel: validate.ErrInvalid,
		Message:  errors.New("Invalid length"),
	}

	// Invalid characters or punctuation.
	ErrFormatting = &validate.ValidatorError{
		ErrLevel: validate.ErrInvalid,
		Message:  errors.New("Invalid formatting"),
	}

	// Unknown error, probably really bad.
	ErrUnknown = &validate.ValidatorError{
		ErrLevel: 2,
		Message:  errors.New("Unknown error"),
	}
)

// Checks for a valid domain name. Checks lengths, characters, and looks for a
// valid TLD (according to IANA).
func (d *Domain) Validate(v validate.Validator) validate.Error {
	//func IsDomain(p []byte) (res validate.Result) {
	// Domain rules:
	// - 255 character total length max
	// - 63 character label max
	// - 127 sub-domains
	// - Characters a-z, A-Z, 0-9, and -
	// - Labels may not start or end with -
	// - TLD may not be all numeric

	// Check for max length.
	// NOTE: Invalid unicode will count as a 1 byte rune, but we'll catch that
	// later.

	p := d.domain
	// If a max length was specified, use it
	if d.checks["maxlength"] != nil &&
		utf8.RuneCount(p) > d.checks["maxlength"].(int) {
		return ErrDomainLength

	} else if utf8.RuneCount(p) > 255 {
		return ErrDomainLength
	}

	// First we split by label
	domain := bytes.Split(p, []byte("."))
	// 127 sub-domains max (not including TLD)
	if len(domain) > 128 {
		return ErrDomainLength
	}

	if d.checks["minsubs"] != nil && len(domain) < d.checks["minsubs"].(int)+1 {
		return ErrDomainLength
	}
	if d.checks["maxsubs"] != nil && len(domain) > d.checks["maxsubs"].(int)+1 {
		return ErrDomainLength
	}
	// Check each domain for valid characters
	for _, subDomain := range domain {
		length := len(subDomain)
		// Check for a domain with two periods next to eachother.
		if length < 1 {
			return ErrFormatting
		}

		// Check 63 character max.
		if length > 63 {
			return ErrDomainLength
		}

		// Check that label doesn't start or end with hyphen.
		r, size := utf8.DecodeRune(subDomain)
		if r == utf8.RuneError && size == 1 {
			// Invalid rune
			return validate.ErrInvalidUTF8
		}

		if r == '-' {
			return ErrFormatting
		}

		r, size = utf8.DecodeLastRune(subDomain)
		if r == utf8.RuneError && size == 1 {
			// Invalid rune
			return validate.ErrInvalidUTF8
		}

		if r == '-' {
			return ErrFormatting
		}

		// Now we check each rune individually to make sure its valid unicode
		// and an acceptable character.
		for i := 0; i < length; {
			if subDomain[i] < utf8.RuneSelf {
				// Check if it's a valid domain character
				if !unicode.Is(domainTable, rune(subDomain[i])) {
					return ErrFormatting
				}
				i++
			} else {
				r, size := utf8.DecodeRune(subDomain[i:])
				if size == 1 {
					// All valid runes of size 1 (those
					// below RuneSelf) were handled above.
					// This must be a RuneError.
					return validate.ErrInvalidUTF8
				}
				// Check if it's a valid domain character
				if !unicode.Is(domainTable, r) {
					return ErrFormatting
				}
				i += size
			}
		}
	}

	// We have all valid unicode characters, now make sure the TLD is real.
	// TODO(inhies): Add check for an all numeric TLD.
	domainTLD := domain[len(domain)-1]
	if tld.Valid(domainTLD) {
		return nil
	}

	// Not sure how we got here, but lets return false just in case.
	return ErrUnknown
}

// Update the included list of TLDs from the given URL.
// Uses github.com/inhies/go-tld
func UpdateTLDs(url string) (err error) {
	return tld.Update(url)
}
