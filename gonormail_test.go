package gonormail

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultNormalizer(t *testing.T) {
	tests := []struct {
		name string
		want *Normalizer
	}{
		{
			want: &Normalizer{
				localFuncs:  defaultFuncs,
				domainFuncs: defaultFuncs,
				localFuncsByDomain: map[string]NormalizeFuncs{
					DomainGmail:      gmailLocalFuncs,
					DomainGmailAlias: gmailLocalFuncs,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, DefaultNormalizer(), "DefaultNormalizer")
		})
	}
}

func TestNewNormalizer(t *testing.T) {
	domainFuncs := NormalizeFuncs{strings.ToLower}
	localFuncs := NormalizeFuncs{strings.ToUpper}
	funcMap := map[string]NormalizeFuncs{"whatever.com": NormalizeFuncs{strings.ToTitle}}

	type args struct {
		localFuncs  NormalizeFuncs
		domainFuncs NormalizeFuncs
		funcMap     map[string]NormalizeFuncs
	}
	tests := []struct {
		name string
		args args
		want *Normalizer
	}{
		{
			name: "empty func",
			args: args{
				localFuncs:  NormalizeFuncs{},
				domainFuncs: NormalizeFuncs{},
				funcMap:     map[string]NormalizeFuncs{},
			},
			want: &Normalizer{localFuncs: NormalizeFuncs{}, domainFuncs: NormalizeFuncs{},
				localFuncsByDomain: map[string]NormalizeFuncs{}},
		},
		{
			name: "same func",
			args: args{
				localFuncs:  localFuncs,
				domainFuncs: domainFuncs,
				funcMap:     funcMap,
			},
			want: &Normalizer{localFuncs: localFuncs, domainFuncs: domainFuncs, localFuncsByDomain: funcMap},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewNormalizer(tt.args.localFuncs, tt.args.domainFuncs, tt.args.funcMap), "NewNormalizer")
		})
	}
}

func TestNormalizer_Normalize(t *testing.T) {
	type fields struct {
		localFuncs         NormalizeFuncs
		domainFuncs        NormalizeFuncs
		localFuncsByDomain map[string]NormalizeFuncs
	}
	type args struct {
		email string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNormalizer(tt.fields.localFuncs, tt.fields.domainFuncs, tt.fields.localFuncsByDomain)
			assert.Equal(t, tt.want, n.Normalize(tt.args.email), "Normalizer.Normalize()")
		})
	}
}

func TestNormalize(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{args: args{email: "not.A.email"}, want: "not.A.email"},
		{args: args{email: "not@@Email"}, want: "not@@Email"},
		{args: args{email: "abcd@email.com"}, want: "abcd@email.com"},
		{args: args{email: "Abcd@Email.com"}, want: "abcd@email.com"},
		{args: args{email: "A.B.C.D+001@Gmail.com"}, want: "abcd@gmail.com"},
		{args: args{email: "A.B.C..D+001@googlemail.com"}, want: "abcd@googlemail.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Normalize(tt.args.email), "Normalize")
		})
	}
}

func TestDeleteDots(t *testing.T) {
	type args struct {
		localPart string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{args: args{localPart: ""}, want: ""},
		{args: args{localPart: "a.b"}, want: "ab"},
		{args: args{localPart: "a.b.c"}, want: "abc"},
		{args: args{localPart: ".a.b.c."}, want: "abc"},
		{args: args{localPart: "a..b...c"}, want: "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, DeleteDots(tt.args.localPart), "DeleteDots")
		})
	}
}

func TestCutPlusRight(t *testing.T) {
	type args struct {
		localPart string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{args: args{localPart: ""}, want: ""},
		{args: args{localPart: "a+b"}, want: "a"},
		{args: args{localPart: "a+b+c"}, want: "a"},
		{args: args{localPart: "+c"}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, CutPlusRight(tt.args.localPart), "CutPlusRight")
		})
	}
}

func TestNormalizer_RegisterLocalFuncs(t *testing.T) {
	type args struct {
		normalizedDomain string
		funcs            NormalizeFuncs
	}
	tests := []struct {
		name            string
		normalizer      *Normalizer
		args            args
		expectedFuncMap map[string]NormalizeFuncs
	}{
		{
			normalizer: &Normalizer{},
			args: args{
				normalizedDomain: "hotmail.com",
				funcs:            nil,
			},
			expectedFuncMap: map[string]NormalizeFuncs{
				"hotmail.com": nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.normalizer.RegisterLocalFuncs(tt.args.normalizedDomain, tt.args.funcs...)
			assert.Equal(t, tt.expectedFuncMap, tt.normalizer.localFuncsByDomain, "RegisterLocalFuncs")
		})
	}
}
