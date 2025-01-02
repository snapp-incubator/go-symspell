# SymSpell Package
[![GoDoc](https://godoc.org/github.com/snapp-incubator/symspell?status.svg)](https://pkg.go.dev/github.com/snapp-incubator/symspell)
[![Go Report Card](https://goreportcard.com/badge/github.com/snapp-incubator/symspell)](https://goreportcard.com/report/github.com/snapp-incubator/symspell)

## Overview

The `symspell` package provides a Golang implementation of the SymSpell algorithm, a fast and memory-efficient algorithm
for spelling correction, word segmentation, and fuzzy string matching. It supports both unigrams and bigrams for
advanced contextual correction.

## Features

- Fast lookup for single-word corrections
- Compound word corrections
- Customizable edit distance and prefix length
- Support for unigram and bigram dictionaries
- Configurable thresholds for performance tuning

## Installation

Install the package using `go get`:

```sh
go get github.com/snapp-incubator/symspell
```

## Usage

- Import the Package

- import "github.com/snapp-incubator/symspell"

- Initialize SymSpell

- Simple Lookup

##### Lookup

###### Load a unigram dictionary:

```go
package main

import "github.com/snapp-incubator/symspell"

func main() {
    symSpell := symspell.NewSymSpellWithLoadDictionary("path/to/vocab.txt", 0, 1,
        symspell.WithCountThreshold(10),
        symspell.WithMaxDictionaryEditDistance(3),
        symspell.WithPrefixLength(5),
    )
}
```

##### Compound Lookup

###### Load both unigram and bigram dictionaries:

```go
package main

func main()  {
    symSpell := symspell.NewSymSpellWithLoadBigramDictionary("path/to/vocab.txt", "path/to/vocab_bigram.txt", 0, 1,
        symspell.WithCountThreshold(1),
        symspell.WithMaxDictionaryEditDistance(3),
        symspell.WithPrefixLength(7),
    )
}
```


### Perform Lookup

#### Single Word Lookup
```go
suggestions, err := symSpell.Lookup("حیابان", symspell.Top, 3)
if err != nil {
    log.Fatal(err)
}
fmt.Println(suggestions[0].Term) // Output: خیابان
```

Compound Word Lookup
```go
suggestion := symSpell.LookupCompound("حیابان ملاصدزا", 3)
fmt.Println(suggestion.Term) // Output: خیابان ملاصدرا
```

## Examples

#### Unit Tests

The repository includes comprehensive unit tests. Run the tests with:
```shell
go test ./...
```

Example test cases include single-word corrections, compound word corrections, and edge cases.

### Configuration Options
- WithMaxDictionaryEditDistance: Sets the maximum edit distance for corrections.
- WithPrefixLength: Sets the prefix length for index optimization.
- WithCountThreshold: Filters dictionary entries with low frequency.

Dictionaries

The dictionaries should be formatted as plain text files:
- Unigram file: Each line should contain a term and its frequency, separated by a space.(or could be custom seperator)
- Bigram file: Each line should contain two terms and their frequency, separated by a space.

#### Example:

Unigram (vocab.txt):
```text
خیابان 1000
میدان 800
```

Bigram (vocab_bigram.txt):
```text
خیابان کارگر 500
میدان آزادی 300
```

### Performance

SymSpell is optimized for speed and memory efficiency. For large vocabularies, tune maxEditDistance, prefixLength, and
countThreshold to balance performance and accuracy.
