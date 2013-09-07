// random-string generates a configurable random string, for passwords.
package main

import (
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"
)

// Flags.
var (
	PasswordLength = flag.Int("length", 16, "Length of password to generate.")
	LowerAlphaNum  = flag.Bool("loweralphanum", false,
		"Only generate passwords with lowercase alphanumeric characters. "+
			"Mutually exclusive with --alphanum.")
	AlphaNum = flag.Bool("alphanum", false,
		"Only generate passwords with alphanumeric characters. "+
			"Mutually exclusive with --loweralphanum.")
	Quiet = flag.Bool("quiet", false, "Suppress unnecessary output.")
)

// Possible characters to be used in passwords.
var (
	LowerAlphaNumChars = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x',
		'y', 'z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	UpperAlphaChars = []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K',
		'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y',
		'Z'}
	SymbolChars = []byte{'~', '!', '@', '#', '%', '^', '&', '*', '(', ')', '-',
		'_', '=', '+', '[', ']', '{', '}', '|', ';', ':', '<', '>', ',', '.',
		'/', '?'}
)

func main() {
	flag.Parse()
	if *LowerAlphaNum && *AlphaNum {
		log.Fatal("--loweralphanum and --alphanum are mutually exclusive!")
	}

	chars := []byte(LowerAlphaNumChars)
	if !*LowerAlphaNum {
		chars = append(chars, UpperAlphaChars...)
	}
	if !*AlphaNum && !*LowerAlphaNum {
		chars = append(chars, SymbolChars...)
	}
	numChars := len(chars)
	numCharsBig := big.NewInt(int64(numChars)) // used for crypto/rand
	bitsPerChar := math.Log2(float64(numChars))
	totalBits := bitsPerChar * float64(*PasswordLength)

	if !*Quiet {
		fmt.Printf("%d characters long.\n", *PasswordLength)
		fmt.Printf("Choosing from %d characters. %2.3f bits of entropy per character.\n",
			numChars, bitsPerChar)
		fmt.Printf("%2.3f total bits of entropy.\n\n", totalBits)
		fmt.Println("Password:")
	}

	var passBuffer bytes.Buffer
	for i := 0; i < *PasswordLength; i++ {
		randBig, err := rand.Int(rand.Reader, numCharsBig)
		if err != nil {
			log.Fatalf("%s", err)
		}
		randChar := chars[randBig.Int64()]
		passBuffer.WriteByte(randChar)
	}

	fmt.Println(passBuffer.String())
}
