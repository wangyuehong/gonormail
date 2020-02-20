package gonormail

import (
	"fmt"
	"strings"
	"sync"
)

const (
	AT    = "@"
	DOT   = "."
	PLUS  = "+"
	EMPTY = ""

	domainGmail      = "gmail.com"
	domainGmailAlias = "googlemail.com"
)

type NormalizeFunc func(string) string

type NormalizeFuncs []NormalizeFunc

var (
	defaultFuncs         = NormalizeFuncs{strings.ToLower}
	gmailLocalFuncs      = NormalizeFuncs{DeleteDots, CutPlusRight}
	defaultGSuiteDomains = []string{domainGmail, domainGmailAlias}
	defaultNormalizer    = DefaultNormalizer()
)

// Normalizer struct that holding normalization funcs.
type Normalizer struct {
	mux                sync.Mutex
	localFuncs         NormalizeFuncs
	domainFuncs        NormalizeFuncs
	localFuncsByDomain map[string]NormalizeFuncs
}

// DefaultNormalizer create default Normalizer.
// domainFuncs: normalize domain part to lower case.
// localFuncs: normalize local part to lower case.
// localFuncsByDomain: local part of gmail domain will be normalized by dot(".") deletion and plus("+") cutting.
func DefaultNormalizer() *Normalizer {
	localFuncsByDomain := make(map[string]NormalizeFuncs, len(defaultGSuiteDomains))
	for _, domain := range defaultGSuiteDomains {
		localFuncsByDomain[domain] = gmailLocalFuncs
	}
	return NewNormalizer(defaultFuncs, defaultFuncs, localFuncsByDomain)
}

// NewNormalizer create new Normalizer by given domainFuncs, localFuncs and localFuncsByDomain.
func NewNormalizer(domainFuncs, localFuncs NormalizeFuncs, localFuncsByDomain map[string]NormalizeFuncs) *Normalizer {
	normalizedMap := make(map[string]NormalizeFuncs, len(localFuncsByDomain))
	for domain, lfuncs := range localFuncsByDomain {
		ndomain := normalize(domainFuncs, domain)
		if _, ok := normalizedMap[ndomain]; ok {
			panic("duplicated normalized domain")
		}
		normalizedMap[ndomain] = lfuncs
	}
	return &Normalizer{localFuncs: localFuncs, domainFuncs: domainFuncs, localFuncsByDomain: normalizedMap}
}

func normalize(funcs NormalizeFuncs, str string) string {
	for _, f := range funcs {
		if f != nil {
			str = f(str)
		}
	}
	return str
}

// Normalize normalize given email parameter by registered functions.
// local part and domain part of the email will by normalized by the localFuncs and domainFuncs,
// then registered normalization functions by domain part will by applied to local part.
// email parameter should by validated before calling this method.
func (n *Normalizer) Normalize(email string) string {
	splitted := strings.Split(email, AT)
	if len(splitted) != 2 {
		return email
	}

	domainPart := normalize(n.localFuncs, splitted[1])
	localPart := normalize(n.domainFuncs, splitted[0])

	if n.localFuncsByDomain != nil {
		if funcs, ok := n.localFuncsByDomain[domainPart]; ok {
			localPart = normalize(funcs, localPart)
		}
	}

	return fmt.Sprintf("%s%s%s", localPart, AT, domainPart)
}

// RegisterLocalFuncs register normalize functions for local part by domain.
// if the domain has been registered already. the functions given will be appended to the end of functions.
func (n *Normalizer) RegisterLocalFuncs(domain string, funcs ...NormalizeFunc) *Normalizer {
	n.mux.Lock()
	defer n.mux.Unlock()

	if n.localFuncsByDomain == nil {
		n.localFuncsByDomain = map[string]NormalizeFuncs{}
	}

	domain = normalize(n.domainFuncs, domain)
	if _, ok := n.localFuncsByDomain[domain]; ok {
		n.localFuncsByDomain[domain] = append(n.localFuncsByDomain[domain], funcs...)
	} else {
		n.localFuncsByDomain[domain] = funcs
	}
	return n
}

// Normalize normalize given email by default Normalizer
func Normalize(email string) string {
	return defaultNormalizer.Normalize(email)
}

// DeleteDots return string with all dot(".") deleted from the given local part.
func DeleteDots(localPart string) string {
	return strings.ReplaceAll(localPart, DOT, EMPTY)
}

// CutPlusRight cut the first plus("+") and the right part of then given local part.
func CutPlusRight(localPart string) string {
	return strings.Split(localPart, PLUS)[0]
}
