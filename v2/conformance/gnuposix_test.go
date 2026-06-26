//go:build conformance

package conformance_test

// POSIX/GNU command-line argument syntax conformance.
//
// Besides being a drop-in for the stdlib flag package, pflag's README promises
// compatibility with "the GNU extensions to the POSIX recommendations for
// command-line options". The two authoritative sources are:
//
//   - POSIX "Utility Syntax Guidelines" 1-14 (the base rules):
//     https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap12.html
//   - GNU "Argument Syntax" (the long-option extension layer):
//     https://www.gnu.org/software/libc/manual/html_node/Argument-Syntax.html
//
// POSIX is authoritative: where this spec conflicts with the stdlib-flag drop-in
// spec (flag_go*_test.go), POSIX wins and the conflict is recorded explicitly in
// divergences.json (enforced by the oracle in hack/oracle). There is no upstream
// test to vendor (these rules are
// prose), so each case below is hand-written and cites the Guideline (Gn) or GNU
// rule it checks. Like the rest of the package it builds only under the
// "conformance" tag and is expected to fail until v2 implements the native flag
// API.

import (
	"io"
	"slices"
	"testing"

	. "github.com/spf13/pflag/v2"
)

func newGNUFlagSet() *FlagSet {
	fs := NewFlagSet("gnuposix", ContinueOnError)
	fs.SetOutput(io.Discard)
	return fs
}

func gnuParse(t *testing.T, fs *FlagSet, args ...string) {
	t.Helper()
	if err := fs.Parse(args); err != nil {
		t.Fatalf("Parse(%q) returned unexpected error: %v", args, err)
	}
}

func wantRemainingArgs(t *testing.T, fs *FlagSet, want ...string) {
	t.Helper()
	if got := fs.Args(); !slices.Equal(got, want) {
		t.Errorf("non-option args = %q, want %q", got, want)
	}
}

// TestShortOptions covers the POSIX short-option guidelines.
func TestShortOptions(t *testing.T) {
	// G4: "All options should be preceded by the '-' delimiter character."
	// G3: "Each option name should be a single alphanumeric character."
	t.Run("single short option", func(t *testing.T) {
		fs := newGNUFlagSet()
		v := fs.BoolP("verbose", "v", false, "")
		gnuParse(t, fs, "-v")
		if !*v {
			t.Errorf("-v: verbose = false, want true")
		}
		wantRemainingArgs(t, fs)
	})

	// G5: "One or more options without option-arguments ... should be accepted
	// when grouped behind one '-' delimiter."
	t.Run("grouped short options", func(t *testing.T) {
		fs := newGNUFlagSet()
		a := fs.BoolP("all", "a", false, "")
		b := fs.BoolP("brief", "b", false, "")
		c := fs.BoolP("count", "c", false, "")
		gnuParse(t, fs, "-abc")
		if !*a || !*b || !*c {
			t.Errorf("-abc: got a=%v b=%v c=%v, want all true", *a, *b, *c)
		}
	})

	// G6: "Each option and option-argument should be a separate argument",
	// except (Utility Argument Syntax item 2) an option-argument may be in the
	// same token as its option.
	t.Run("short option argument as separate token", func(t *testing.T) {
		fs := newGNUFlagSet()
		o := fs.StringP("output", "o", "", "")
		gnuParse(t, fs, "-o", "file.txt")
		if *o != "file.txt" {
			t.Errorf("-o file.txt: output = %q, want %q", *o, "file.txt")
		}
	})
	t.Run("short option argument in same token", func(t *testing.T) {
		fs := newGNUFlagSet()
		o := fs.StringP("output", "o", "", "")
		gnuParse(t, fs, "-ofile.txt")
		if *o != "file.txt" {
			t.Errorf("-ofile.txt: output = %q, want %q", *o, "file.txt")
		}
	})

	// G5: "... followed by at most one option that takes an option-argument,
	// should be accepted when grouped behind one '-' delimiter."
	t.Run("cluster ending in an option that takes an argument", func(t *testing.T) {
		fs := newGNUFlagSet()
		v := fs.BoolP("verbose", "v", false, "")
		o := fs.StringP("output", "o", "", "")
		gnuParse(t, fs, "-vofile.txt")
		if !*v || *o != "file.txt" {
			t.Errorf("-vofile.txt: got verbose=%v output=%q, want true and %q", *v, *o, "file.txt")
		}
	})
}

// TestLongOptions covers the GNU long-option extension (not part of base POSIX).
func TestLongOptions(t *testing.T) {
	// GNU: long options consist of "--" followed by alphanumeric characters and
	// dashes.
	t.Run("long option", func(t *testing.T) {
		fs := newGNUFlagSet()
		v := fs.BoolP("verbose", "v", false, "")
		gnuParse(t, fs, "--verbose")
		if !*v {
			t.Errorf("--verbose: verbose = false, want true")
		}
	})
	t.Run("long option name may contain dashes", func(t *testing.T) {
		fs := newGNUFlagSet()
		n := fs.BoolP("dry-run", "n", false, "")
		gnuParse(t, fs, "--dry-run")
		if !*n {
			t.Errorf("--dry-run: dry-run = false, want true")
		}
	})

	// GNU: "To specify an argument for a long option, write --name=value."
	t.Run("long option argument with =", func(t *testing.T) {
		fs := newGNUFlagSet()
		o := fs.StringP("output", "o", "", "")
		gnuParse(t, fs, "--output=file.txt")
		if *o != "file.txt" {
			t.Errorf("--output=file.txt: output = %q, want %q", *o, "file.txt")
		}
	})
	// GNU getopt_long also accepts the argument as the following token.
	t.Run("long option argument as separate token", func(t *testing.T) {
		fs := newGNUFlagSet()
		o := fs.StringP("output", "o", "", "")
		gnuParse(t, fs, "--output", "file.txt")
		if *o != "file.txt" {
			t.Errorf("--output file.txt: output = %q, want %q", *o, "file.txt")
		}
	})
}

// TestOptionTerminatorAndDash covers the two special tokens.
func TestOptionTerminatorAndDash(t *testing.T) {
	// G10: "The first '--' argument that is not an option-argument should be
	// accepted as a delimiter indicating the end of options." Anything after is
	// an operand, even if it begins with '-'.
	t.Run("double dash terminates option processing", func(t *testing.T) {
		fs := newGNUFlagSet()
		v := fs.BoolP("verbose", "v", false, "")
		gnuParse(t, fs, "-v", "--", "-x", "file")
		if !*v {
			t.Errorf("-v before --: verbose = false, want true")
		}
		wantRemainingArgs(t, fs, "-x", "file")
	})

	// G13: a single '-' operand denotes standard input/output, i.e. it is an
	// ordinary operand, not an option.
	t.Run("single hyphen is an operand", func(t *testing.T) {
		fs := newGNUFlagSet()
		v := fs.BoolP("verbose", "v", false, "")
		gnuParse(t, fs, "-v", "-", "file")
		if !*v {
			t.Errorf("-v: verbose = false, want true")
		}
		wantRemainingArgs(t, fs, "-", "file")
	})
}

// TestOptionOrdering covers operand ordering and option repetition.
func TestOptionOrdering(t *testing.T) {
	// GNU extension to G9 ("All options should precede operands"): GNU
	// implementations reorder argv so options and operands may be interspersed.
	// pflag uses this GNU behavior by default.
	t.Run("options interspersed with operands (GNU default)", func(t *testing.T) {
		fs := newGNUFlagSet()
		v := fs.BoolP("verbose", "v", false, "")
		o := fs.StringP("output", "o", "", "")
		gnuParse(t, fs, "src1", "-v", "src2", "-o", "out", "src3")
		if !*v || *o != "out" {
			t.Errorf("interspersed: got verbose=%v output=%q, want true and %q", *v, *o, "out")
		}
		wantRemainingArgs(t, fs, "src1", "src2", "src3")
	})

	// G9 (strict POSIX / POSIXLY_CORRECT): option processing stops at the first
	// operand. pflag opts into this via SetInterspersed(false).
	t.Run("strict POSIX ordering stops at first operand", func(t *testing.T) {
		fs := newGNUFlagSet()
		fs.SetInterspersed(false)
		v := fs.BoolP("verbose", "v", false, "")
		gnuParse(t, fs, "src1", "-v")
		if *v {
			t.Errorf("interspersing disabled: -v after an operand should not be parsed")
		}
		wantRemainingArgs(t, fs, "src1", "-v")
	})

	// G11: "The order of different options relative to one another should not
	// matter." Options may also appear multiple times.
	t.Run("an option may appear multiple times", func(t *testing.T) {
		fs := newGNUFlagSet()
		n := fs.CountP("verbose", "v", "")
		gnuParse(t, fs, "-v", "-v", "-vv")
		if *n != 4 {
			t.Errorf("repeated -v: count = %d, want 4", *n)
		}
	})
	t.Run("options may be supplied in any order", func(t *testing.T) {
		fs := newGNUFlagSet()
		a := fs.BoolP("all", "a", false, "")
		o := fs.StringP("output", "o", "", "")
		gnuParse(t, fs, "-o", "out", "-a")
		if !*a || *o != "out" {
			t.Errorf("got all=%v output=%q, want true and %q", *a, *o, "out")
		}
	})
}

// TestLongOptionAbbreviation encodes the one GNU rule pflag has historically
// chosen NOT to implement:
//
//	"the user can abbreviate the option name as long as the abbreviation is
//	 unique."
//
// This is a pflag-vs-GNU divergence (not a stdlib conflict), recorded in
// divergences.json. It is kept here as the spec: when v2 reaches this point,
// decide deliberately — implement unambiguous abbreviation, or change this to
// t.Skip documenting the divergence. Do not just delete it.
func TestLongOptionAbbreviation(t *testing.T) {
	fs := newGNUFlagSet()
	verbose := fs.BoolP("verbose", "v", false, "")
	gnuParse(t, fs, "--verb")
	if !*verbose {
		t.Errorf("--verb: should be accepted as an unambiguous abbreviation of --verbose")
	}
}
