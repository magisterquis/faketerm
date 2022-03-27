// Package faketerm can be used when either a term.Terminal or just an
// io.Reader/Writer/ReadWriter is needed.
package faketerm

/*
 * faketerm.go
 * Fake terminal which acts like a real one
 * By J. Stuart McMurray
 * Created 20220327
 * Last Modified 20220327
 */

import (
	"bufio"
	"io"
	"strings"
	"sync"
)

// Term is an interface which is satisfied by both term.Terminal and FakeTerm.
// It can be use wherever a term.Terminal or an io.Reade/Writer/ReadWriter may
// be used.  Please see the document for term.Terminal for a description of
// the methods.
type Term interface {
	ReadLine() (line string, err error)
	ReadPassword(prompt string) (line string, err error)
	SetBracketedPasteMode(on bool)
	SetPrompt(prompt string)
	SetSize(width, height int) error
	Write(buf []byte) (n int, err error)
}

// FakeTerm is a Term with an underlying io.ReadWriter.  Its methods are
// analogs of term.Terminals with differences noted.
type FakeTerm struct {
	w  io.Writer
	wL sync.Mutex
	s  *bufio.Scanner
	rL sync.Mutex
}

// New returns a new FakeTerm, ready for use
func New(r io.Reader, w io.Writer) *FakeTerm {
	return &FakeTerm{
		w: w,
		s: bufio.NewScanner(r),
	}
}

func (f *FakeTerm) ReadLine() (line string, err error) {
	f.rL.Lock()
	defer f.rL.Unlock()
	/* Wait for a line to be available. */
	if !f.s.Scan() {
		err := f.s.Err()
		if nil == err {
			return "", io.EOF
		}
		return "", f.s.Err()
	}
	return strings.TrimRight(f.s.Text(), "\r\n"), nil
}

// ReadPassword is a thin wrapper around f.ReadLine.  The prompt is ignored.
func (f *FakeTerm) ReadPassword(prompt string) (line string, err error) {
	return f.ReadLine()
}

// SetBracketedPasteMode is a no-op
func (f *FakeTerm) SetBracketedPasteMode(on bool) {}

// SetPrompt is a no-op.
func (f *FakeTerm) SetPrompt(prompt string) {}

// SetSize is a no-op.
func (f *FakeTerm) SetSize(width, height int) error { return nil }

func (f *FakeTerm) Write(buf []byte) (n int, err error) {
	f.wL.Lock()
	defer f.wL.Unlock()
	return f.w.Write(buf)
}
