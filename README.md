# gonormail
gonormail is a Go library to normalize email or build a email normalizer with the default support of gmail.

## Usage

Normalization by default normalizer supported gmail. email should be validated before normalization.
```golang
gonormail.Normalize("Not A Email")              // Not A Email
gonormail.Normalize("Not@A@Email")              // Not@A@Email
gonormail.Normalize("A.B.c@Gmail.com")          // abc@gmail.com
gonormail.Normalize("a.B..c@gmail.com")         // abc@gmail.com
gonormail.Normalize("a.b.c+001@gmail.com")      // abc@gmail.com
gonormail.Normalize("a.b.c+001@googlemail.com") // abc@gmail.com
gonormail.Normalize("a.b.c+001@whatever.com")   // a.b.c+001@whatever.com
```

Customized normalization by extending the default EmailNormalizer.
```golang
norm := gonormail.DefaultEmailNormalizer().
  AddNormalizer(gonormail.NewDomainAlias(map[string]string{"examplemail.com": "email.com"})).
  AddNormalizer(gonormail.NewRemoveLocalDots("email.com")).
  AddNormalizer(gonormail.NewRemoveSubAddressing(map[string]string{"email.com": "-"}))

norm.Normalize("A.B.c+001@Gmail.com")       // abc@email.com
norm.Normalize("A.b.c+002@googlemail.com")  // abc@email.com
norm.Normalize("A.B.c-003@Examplemail.Com") // abc@email.com
norm.Normalize("A.B.c-004@Email.Com")       // abc@email.com
```

Create a new EmailNormalizer by fully customization.
```golang
norm := gonormail.NewEmailNormalizer().
  AddFunc(gonormail.ToLowerCase).
  AddNormalizer(gonormail.NewDomainAlias(map[string]string{"googlemail.com": "gmail.com", "examplemail.com": "email.com"})).
  AddNormalizer(gonormail.NewRemoveLocalDots("email.com", "gmail.com")).
  AddNormalizer(gonormail.NewRemoveSubAddressing(map[string]string{"email.com": "-", "gmail.com": "+"}))

norm.Normalize("A.B.c+001@Gmail.com")       // abc@gmail.com
norm.Normalize("A.b.c+002@googlemail.com")  // abc@gmail.com
norm.Normalize("A.B.c-003@Examplemail.Com") // abc@email.com
norm.Normalize("A.B.c-004@Email.Com")       // abc@email.com
```

Create your own Normalizer.
```golang
norm := gonormail.NewEmailNormalizer().
  AddFunc(func(e *gonormail.EmailAddress) {
    e.Local = strings.ToUpper(e.Local)
    e.Domain = strings.ToUpper(e.Domain)
})
norm.Normalize("abc@email.com") // ABC@EMAIL.COM
```