package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/scanner"
)

func main() {
	var file string
	flag.StringVar(&file, "f", "", "Facebook profile web page")
	flag.Parse()

	names, e := scan(file)
	n := 0
	if e == nil {
		for i, j := range names {
			fmt.Printf("%d - %s\n", i, j)
		}
	} else {
		fmt.Fprintln(os.Stderr, e.Error())
		n = 1
	}
	os.Exit(n)
}

/*
data = `"shortProfiles"` ":" dict.
dict = "{" {id ":" profile} "}".
profile = "{" key ":" val { "," key ":" val } "}".
var = number | string | string_list.
*/

func scan(file string) (ns []string, e error) {
	var rd io.Reader
	rd, e = os.Open(file)
	s := new(scanner.Scanner)
	if e == nil {
		s.Init(rd)
		s.Filename = file
		s.Error = func(n *scanner.Scanner, msg string) {
			if !strings.HasSuffix(msg, "illegal char escape") &&
				!strings.HasSuffix(msg, "illegal octal number") {
				fmt.Fprintf(os.Stderr, "%s %s\n", n.Position.String(), msg)
			}
		}
		ns, e = scanData(s)
	}
	return
}

func scanData(s *scanner.Scanner) (ns []string, e error) {
	b := true
	for tok := s.Scan(); b && tok != scanner.EOF; tok = s.Scan() {
		b = s.TokenText() != "shortProfiles"
	}
	s.Scan()
	b = s.TokenText() == "{"
	var m []map[string]string
	if b {
		m, e = scanUsrInfo(s)
	} else if e == nil {
		e = fmt.Errorf("No '{' found")
	}

	if e == nil {
		ns = make([]string, len(m))
		i := 0
		for _, v := range m {
			ns[i], i = v["name"], i+1
		}
	}
	return
}

func scanUsrInfo(s *scanner.Scanner) (m []map[string]string, e error) {
	m = make([]map[string]string, 0)
	b := true
	for b {
		s.Scan() // id scanned
		s.Scan() // : scanned
		b = s.TokenText() == ":"
		usr := make(map[string]string)
		if b {
			s.Scan()
			b = s.TokenText() == "{"
		} else {
			e = fmt.Errorf("No ':' found")
		}
		for b {
			s.Scan()
			field := s.TokenText()
			s.Scan()
			b = s.TokenText() == ":"
			var val string
			if b {
				s.Scan()
				if s.TokenText() == "[" {
					// a list
					val, e = scanList(s)
				} else {
					val = s.TokenText()
				}
				if e == nil {
					usr[field] = val
					s.Scan()
					b = s.TokenText() == ","
				}
			}
		}
		b = s.TokenText() == "}"
		if b {
			m = append(m, usr)
			s.Scan()
			b = s.TokenText() == ","
		} else {
			e = fmt.Errorf("No '}' found at %s", s.Position.String())
		}
	}
	b = s.TokenText() == "}"
	return
}

func scanList(s *scanner.Scanner) (v string, e error) {
	v = ""
	b, eof := true, false
	for b && !eof {
		v = v + s.TokenText()
		tok := s.Scan()
		eof = tok == scanner.EOF
		b = s.TokenText() != "]"
	}
	if eof {
		e = fmt.Errorf("EOF")
	}
	return
}
