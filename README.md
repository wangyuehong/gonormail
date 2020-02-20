# gonormail
normalize email with Go

## Usage

default normalize rule.
```golang
gonormail.Normalize("A.B.c@Gmail.com")          // abc@gmail.com
gonormail.Normalize("a.b.c@gmail.com")          // abc@gmail.com
gonormail.Normalize("a.b.c+001@gmail.com")      // abc@gmail.com
gonormail.Normalize("a.b.c+001@googlemail.com") // abc@googlemail.com
```

customized normalization.
```golang
norm := gonormail.DefaultNormalizer().
  RegisterLocalFuncs("live.com", gonormail.DeleteDots, gonormail.CutPlusRight).
  RegisterLocalFuncs("hotmail.com", gonormail.CutPlusRight).
  RegisterLocalFuncs("whatever.com", func(s string) string { return s + "+s" })

norm.Normalize("A.B.c+001@Gmail.com")      // abc@gmail.com
norm.Normalize("A.b.c+002@googlemail.com") // abc@googlemail.com
norm.Normalize("A.B.c+003@Live.Com")       // abc@live.com
norm.Normalize("A.B.c+004@Hotmail.Com")    // a.b.c@hotmail.com
norm.Normalize("hello@Whatever.Com")       // hello+s@whatever.com
```
