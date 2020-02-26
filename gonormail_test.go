package gonormail

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		email string
		want  string
	}{
		{email: "Not A Email", want: "Not A Email"},
		{email: "Not@A@Email", want: "Not@A@Email"},
		{email: "A.B.c@Gmail.com", want: "abc@gmail.com"},
		{email: "a.B..c@Gmail.com", want: "abc@gmail.com"},
		{email: "a.b.c+001@googlemail.com", want: "abc@gmail.com"},
		{email: "a.b.c+001@whatever.com", want: "a.b.c+001@whatever.com"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if got := Normalize(tt.email); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmailNormalizer_Normalize(t *testing.T) {
	tests := []struct {
		emailNormalizer *EmailNormalizer
		email           string
		want            string
	}{
		{
			emailNormalizer: (&EmailNormalizer{}).AddFunc(nil),
			email:           "abc@email.com",
			want:            "abc@email.com",
		},
		{
			emailNormalizer: (&EmailNormalizer{}).AddNormalizer(NormalizeFunc(nil)),
			email:           "abc@email.com",
			want:            "abc@email.com",
		},
		{
			emailNormalizer: (&EmailNormalizer{}).AddFunc(func(e *EmailAddress) { e.Local += "l" },
				func(e *EmailAddress) { e.Domain += "d" }),
			email: "abc@email.com",
			want:  "abcl@email.comd",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if got := tt.emailNormalizer.Normalize(tt.email); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EmailNormalizer.Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveLocalDots_Normalize(t *testing.T) {
	tests := []struct {
		domains []string
		email   string
		want    string
	}{
		{
			domains: []string{"email.com"},
			email:   "a.b.c@email.com",
			want:    "abc@email.com",
		},
		{
			domains: []string{"email.com"},
			email:   "a..b..c..@email.com",
			want:    "abc@email.com",
		},
		{
			domains: []string{"email.com"},
			email:   "a.b.c@cmail.com",
			want:    "a.b.c@cmail.com",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			ea := NewEmailAddress(tt.email)
			NewRemoveLocalDots(tt.domains...).Normalize(ea)
			if got := ea.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveLocalDots.Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveSubAddressing_Normalize(t *testing.T) {
	tests := []struct {
		tags  map[string]string
		email string
		want  string
	}{
		{
			tags:  map[string]string{"email.com": "+"},
			email: "a@email.com",
			want:  "a@email.com",
		},
		{
			tags:  map[string]string{"email.com": "+"},
			email: "a+b+c@email.com",
			want:  "a@email.com",
		},
		{
			tags:  map[string]string{"email.com": "-"},
			email: "a--b-c@email.com",
			want:  "a@email.com",
		},
		{
			tags:  map[string]string{},
			email: "a+b@email.com",
			want:  "a+b@email.com",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			ea := NewEmailAddress(tt.email)
			NewRemoveSubAddressing(tt.tags).Normalize(ea)
			if got := ea.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveSubAddressing.Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainAlias_Normalize(t *testing.T) {
	tests := []struct {
		aliases map[string]string
		email   string
		want    string
	}{
		{
			aliases: map[string]string{
				"examplemail.com": "email.com",
			},
			email: "a@examplemail.com",
			want:  "a@email.com",
		},
		{
			aliases: map[string]string{},
			email:   "a@examplemail.com",
			want:    "a@examplemail.com",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			ea := NewEmailAddress(tt.email)
			NewDomainAlias(tt.aliases).Normalize(ea)
			if got := ea.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DomainAlias.Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToLowerCase(t *testing.T) {
	tests := []struct {
		email string
		want  string
	}{
		{
			email: "Abc@Email.Com",
			want:  "abc@email.com",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			ea := NewEmailAddress(tt.email)
			ToLowerCase(ea)
			if got := ea.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToLowerCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
