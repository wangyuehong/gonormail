package gonormail

import (
	"fmt"
	"strings"
)

const (
	AT               = "@"
	DOT              = "."
	PLUS             = "+"
	EMPTY            = ""
	DomainGmail      = "gmail.com"
	DomainGmailAlias = "googlemail.com"
)

type NormalizeFunc func(string) string

var defaultFuncs = []NormalizeFunc{strings.ToLower}
var gmailLocalFuncs = []NormalizeFunc{DeleteDots, CutPlusRight}
var defaultNormalizer = DefaultNormalizer()

type Normalizer struct {
	localFuncs         []NormalizeFunc
	domainFuncs        []NormalizeFunc
	localFuncsByDomain map[string][]NormalizeFunc
}

func DefaultNormalizer() *Normalizer {
	return NewNormalizer(defaultFuncs, defaultFuncs, map[string][]NormalizeFunc{
		DomainGmail:      gmailLocalFuncs,
		DomainGmailAlias: gmailLocalFuncs,
	})
}

func NewNormalizer(localFuncs, domainFuncs []NormalizeFunc, funcMap map[string][]NormalizeFunc) *Normalizer {
	return &Normalizer{localFuncs: localFuncs, domainFuncs: domainFuncs, localFuncsByDomain: funcMap}
}

func (n *Normalizer) Normalize(email string) string {
	splitted := strings.Split(email, AT)
	if len(splitted) != 2 {
		return email
	}

	localPart, domainPart := splitted[0], splitted[1]
	for _, f := range n.domainFuncs {
		domainPart = f(domainPart)
	}

	for _, f := range n.localFuncs {
		localPart = f(localPart)
	}

	if n.localFuncsByDomain != nil {
		if funcs, ok := n.localFuncsByDomain[domainPart]; ok {
			for _, f := range funcs {
				localPart = f(localPart)
			}
		}
	}

	return fmt.Sprintf("%s%s%s", localPart, AT, domainPart)
}

func Normalize(email string) string {
	return defaultNormalizer.Normalize(email)
}

func DeleteDots(localPart string) string {
	return strings.ReplaceAll(localPart, DOT, EMPTY)
}

func CutPlusRight(localPart string) string {
	return strings.Split(localPart, PLUS)[0]
}
