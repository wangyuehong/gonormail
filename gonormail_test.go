package gonormail

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultNormalizer(t *testing.T) {
	want := &Normalizer{
		localFuncs:  defaultFuncs,
		domainFuncs: defaultFuncs,
		localFuncsByDomain: map[string]NormalizeFuncs{
			domainGmail:      gmailLocalFuncs,
			domainGmailAlias: gmailLocalFuncs,
		},
	}
	assert.Equal(t, want, DefaultNormalizer(), "DefaultNormalizer")
}

func TestNormalizer_RegisterLocalFuncs(t *testing.T) {
	type fields struct {
		localFuncs         NormalizeFuncs
		domainFuncs        NormalizeFuncs
		localFuncsByDomain map[string]NormalizeFuncs
	}
	type args struct {
		domain string
		funcs  []NormalizeFunc
	}
	tests := []struct {
		fields fields
		argss  []args
		email  string
		want   string
	}{
		{
			fields: fields{
				localFuncs:         nil,
				domainFuncs:        nil,
				localFuncsByDomain: nil,
			},
			argss: []args{
				{
					domain: "",
					funcs:  nil,
				},
				{
					domain: "",
					funcs:  nil,
				},
			},
			email: "abc@email.com",
			want:  "abc@email.com",
		},
		{
			fields: fields{
				localFuncs:         defaultFuncs,
				domainFuncs:        defaultFuncs,
				localFuncsByDomain: nil,
			},
			argss: []args{
				{
					domain: "email.COM",
					funcs: NormalizeFuncs{
						func(s string) string { return s + "+" },
						func(s string) string { return s + "m" },
					},
				},
				{
					domain: "EMAIL.com",
					funcs: NormalizeFuncs{
						nil,
						func(s string) string { return s + "n" },
					},
				},
			},
			email: "ABC@EMAIL.COM",
			want:  "abc+mn@email.com",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			n := NewNormalizer(tt.fields.domainFuncs, tt.fields.localFuncs, tt.fields.localFuncsByDomain)
			for _, args := range tt.argss {
				n = n.RegisterLocalFuncs(args.domain, args.funcs...)
			}
			assert.Equal(t, tt.want, n.Normalize(tt.email), "RegisterLocalFuncs.Normalize()")
		})
	}
}

func TestNormalizer_Normalize(t *testing.T) {
	type fields struct {
		localFuncs         NormalizeFuncs
		domainFuncs        NormalizeFuncs
		localFuncsByDomain map[string]NormalizeFuncs
	}
	tests := []struct {
		fields fields
		email  string
		want   string
	}{
		{
			fields: fields{
				localFuncs:         nil,
				domainFuncs:        nil,
				localFuncsByDomain: nil,
			},
			email: "abc@email.com",
			want:  "abc@email.com",
		},
		{
			fields: fields{
				localFuncs:         nil,
				domainFuncs:        nil,
				localFuncsByDomain: map[string]NormalizeFuncs{"email.com": nil},
			},
			email: "abc@email.com",
			want:  "abc@email.com",
		},
		{
			fields: fields{
				localFuncs:  NormalizeFuncs{strings.ToUpper},
				domainFuncs: NormalizeFuncs{strings.ToUpper},
				localFuncsByDomain: map[string]NormalizeFuncs{
					"EMAIL.COM": NormalizeFuncs{func(s string) string { return s + "+s" }},
				},
			},
			email: "abc@email.com",
			want:  "ABC+s@EMAIL.COM",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			n := NewNormalizer(tt.fields.domainFuncs, tt.fields.localFuncs, tt.fields.localFuncsByDomain)
			assert.Equal(t, tt.want, n.Normalize(tt.email), "Normalizer.Normalize()")
		})
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		email string
		want  string
	}{
		{email: "not.A.email", want: "not.A.email"},
		{email: "not@@Email", want: "not@@Email"},
		{email: "abcd@email.com", want: "abcd@email.com"},
		{email: "Abcd@Email.com", want: "abcd@email.com"},
		{email: "A.B.C.D+001@Gmail.com", want: "abcd@gmail.com"},
		{email: "A.B.C..D+001@googlemail.com", want: "abcd@googlemail.com"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			assert.Equal(t, tt.want, Normalize(tt.email), "Normalize")
		})
	}
}

func TestDeleteDots(t *testing.T) {
	tests := []struct {
		localPart string
		want      string
	}{
		{localPart: "", want: ""},
		{localPart: ".", want: ""},
		{localPart: "a.b", want: "ab"},
		{localPart: "a.b.c", want: "abc"},
		{localPart: ".a.b.c.", want: "abc"},
		{localPart: "a..b...c", want: "abc"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			assert.Equal(t, tt.want, DeleteDots(tt.localPart), "DeleteDots")
		})
	}
}

func TestCutPlusRight(t *testing.T) {
	tests := []struct {
		localPart string
		want      string
	}{
		{localPart: "", want: ""},
		{localPart: "+", want: ""},
		{localPart: "a+b", want: "a"},
		{localPart: "a+b+c", want: "a"},
		{localPart: "+c", want: ""},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			assert.Equal(t, tt.want, CutPlusRight(tt.localPart), "CutPlusRight")
		})
	}
}

func Test_normalize(t *testing.T) {
	tests := []struct {
		funcs NormalizeFuncs
		str   string
		want  string
	}{
		{
			funcs: nil,
			str:   "a",
			want:  "a",
		},
		{
			funcs: NormalizeFuncs{nil},
			str:   "a",
			want:  "a",
		},
		{
			funcs: NormalizeFuncs{
				func(s string) string { return s + "j" },
				func(s string) string { return s + "k" },
			},
			str:  "a",
			want: "ajk",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			assert.Equal(t, tt.want, normalize(tt.funcs, tt.str), "normalize")
		})
	}
}
