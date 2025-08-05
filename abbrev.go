// abbrev.go - generate abbreviations from a wordlist
//
// (c) 2014 Sudhi Herle <sw-at-herle.net>
//
// Placed in the Public Domain
// This software does not come with any express or implied
// warranty; it is provided "as is". No claim  is made to its
// suitability for any purpose.

package pflag

// Given a wordlist in 'words', generate unique abbreviations of it
// and return as a map[abbrev]word.
// e.g.,
//
//	given a wordlist ["hello", "help", "sync"],
//	Abbrev() returns:
//	  {
//	     "hello": "hello",
//	     "hell":  "hell"
//	     "help":  "help",
//	     "sync":  "sync",
//	     "syn":   "sync",
//	     "sy":    "sync",
//	     "s":     "sync"
//	  }
func abbrev(words []string) map[string]string {
	seen := make(map[string]int)
	table := make(map[string]string)

	for _, w := range words {
		for n := len(w) - 1; n > 0; n -= 1 {
			ab := w[:n]
			seen[ab] += 1

			switch seen[ab] {
			case 1:
				table[ab] = w
			case 2:
				delete(table, ab)
			default:
				goto next
			}
		}
	next:
	}

	// non abbreviations always get entered
	// This has to be done _after_ the loop above; because
	// if there are words that are prefixes of other words in
	// the argument list, we need to ensure we capture them
	// intact.
	for _, w := range words {
		table[w] = w
	}
	return table
}
