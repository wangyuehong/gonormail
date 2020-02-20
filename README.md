# gonormail
gonormail is a Go library to normalize email or build a email normalizer with default support of gmail.

## normalization
 - normalize the local and domain parts of email to lower case.
 - normalize the local part for gmail
   - delete dots(`.`) from local part.
   - delete sub-addressing starting with(`+`).

## Usage

simple normalization. email should be validated before normalization.
```golang
gonormail.Normalize("Not A Email")              // Not A Email
gonormail.Normalize("Not@A@Email")              // Not@A@Email
gonormail.Normalize("A.B.c@Gmail.com")          // abc@gmail.com
gonormail.Normalize("a.b.c@gmail.com")          // abc@gmail.com
gonormail.Normalize("a.b.c+001@gmail.com")      // abc@gmail.com
gonormail.Normalize("a.b.c+001@googlemail.com") // abc@googlemail.com
gonormail.Normalize("a.b.c+001@whatever.com")   // a.b.c+001@whatever.com
```

customized normalization.
```golang
norm := gonormail.DefaultNormalizer().
  Register("live.com", gonormail.DeleteDots, gonormail.CutPlusRight).
  Register("hotmail.com", gonormail.CutPlusRight).
  Register("whatever.com", func(s string) string { return s + "+s" })

norm.Normalize("A.B.c+001@Gmail.com")      // abc@gmail.com
norm.Normalize("A.b.c+002@googlemail.com") // abc@googlemail.com
norm.Normalize("A.B.c+003@Live.Com")       // abc@live.com
norm.Normalize("A.B.c+004@Hotmail.Com")    // a.b.c@hotmail.com
norm.Normalize("hello@Whatever.Com")       // hello+s@whatever.com
```
