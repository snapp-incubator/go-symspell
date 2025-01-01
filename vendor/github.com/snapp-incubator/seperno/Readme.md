Here’s the README content in Markdown format:

# Persian Text Normalizer - Seperno 
### *Search Persian Normalizer*

Seperno is a powerful and customizable Persian text normalizer. It provides various tools to clean, preprocess, and
normalize Persian text by handling spaces, URLs, punctuation, and more. This project is part of *
*[Snapp Incubator](https://github.com/snapp-incubator/seperno)** and aims to simplify text processing in Persian
applications.

---

## Features

- **Convert Half-Space to Space**: Converts Persian half-spaces (`\u200c`) into regular spaces.
- **Remove URLs**: Cleans text by removing URLs.
- **Combine Multiple Spaces**: Reduces multiple consecutive spaces into a single space.
- **Remove Outer Spaces**: Trims unnecessary spaces from the start and end of the text.
- **Remove End-of-Line Characters**: Removes specific characters like `.` or `؟` at the end of a sentence.
- **Normalize Punctuation**: Replaces punctuation marks with spaces or their normalized equivalents.
- **Customizable**: Use modular options to tailor the normalization process.

---

## Installation

Install the package using `go get`:

```bash
go get github.com/snapp-incubator/seperno
```

## Usage

## Basic Example

```go

package main

import (
	"fmt"
	"github.com/snapp-incubator/seperno"
)

func main() {
	normalizer := seperno.NewNormalize(
		seperno.WithConvertHalfSpaceToSpace(),
		seperno.WithURLRemover(),
		seperno.WithSpaceCombiner(),
	)

	text := "اسنپ‌کپ  تست   https://example.com"
	normalized := normalizer.BasicNormalizer(text)

	fmt.Println(normalized) // Output: "اسنپ کپ تست"
}
```

## Advanced Examples

#### Convert Half-Space to Space

```go
package main

import (
	"fmt"
	"github.com/snapp-incubator/seperno"
)

func main() {
	normalizer := seperno.NewNormalize(seperno.WithConvertHalfSpaceToSpace())
	text := "آسمان‌آبی"
	fmt.Println(normalizer.BasicNormalizer(text)) // Output: "اسمان ابی"
}
```

#### Remove URLs

```go
package main

import (
	"fmt"
	"github.com/snapp-incubator/seperno"
)

func main() {
	normalizer := seperno.NewNormalize(seperno.WithURLRemover())
	text := "تست https://example.com"
	fmt.Println(normalizer.BasicNormalizer(text)) // Output: "تست "
}
```

#### Combine Multiple Spaces

```go
package main

import (
	"fmt"
	"github.com/snapp-incubator/seperno"
)

func main() {
	normalizer := seperno.NewNormalize(seperno.WithSpaceCombiner())
	text := "تست تست"
	fmt.Println(normalizer.BasicNormalizer(text)) // Output: "تست تست"
}
```

#### Remove Outer Spaces

```go
package main

import (
	"fmt"
	"github.com/snapp-incubator/seperno"
)

func main() {
	normalizer := seperno.NewNormalize(seperno.WithOuterSpaceRemover())
	text := "   تست   "
	fmt.Println(normalizer.BasicNormalizer(text)) // Output: "تست"
}
```

#### Remove End-of-Line Characters

```go
package main

import (
	"fmt"
	"github.com/snapp-incubator/seperno"
)

func main() {
	normalizer := seperno.NewNormalize(seperno.WithEndsWithEndOfLineChar())
	text := "تست."
	fmt.Println(normalizer.BasicNormalizer(text)) // Output: "تست"
}
```

#### Normalize Punctuation

```go
package main

import (
	"fmt"
	"github.com/snapp-incubator/seperno"
)

func main() {
	normalizer := seperno.NewNormalize(
		seperno.WithNormalizePunctuations(),
		seperno.WithOuterSpaceRemover(),
	)
	text := "سلام,خوبی؟چه خبرا."
	fmt.Println(normalizer.BasicNormalizer(text)) // Output: "سلام خوبی چه خبرا"
}
```

### Run Tests

To validate functionality, run the included test suite:

```bash
go test ./...
```
