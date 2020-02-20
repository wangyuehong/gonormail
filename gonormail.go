package gonormail

import (
	"fmt"
	"strings"
	"sync"
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
type NormalizeFuncs []NormalizeFunc

var defaultFuncs = NormalizeFuncs{strings.ToLower}
var gmailLocalFuncs = NormalizeFuncs{DeleteDots, CutPlusRight}
var defaultNormalizer = DefaultNormalizer()

type Normalizer struct {
	mux                sync.Mutex
	localFuncs         NormalizeFuncs
	domainFuncs        NormalizeFuncs
	localFuncsByDomain map[string]NormalizeFuncs
}

func DefaultNormalizer() *Normalizer {
	return NewNormalizer(defaultFuncs, defaultFuncs, map[string]NormalizeFuncs{
		DomainGmail:      gmailLocalFuncs,
		DomainGmailAlias: gmailLocalFuncs,
	})
}

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
	for _, fuc := range funcs {
		str = fuc(str)
	}
	return str
}

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

func (n *Normalizer) RegisterLocalFuncs(domain string, funcs ...NormalizeFunc) {
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
