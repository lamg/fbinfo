package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	var file string
	flag.StringVar(&file, "f", "", "Facebook profile web page")
	flag.Parse()

	names, e := parse(file)
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

// file cannot contain spaces after "shortProfiles:"
func parse(file string) (ns []string, e error) {
	var fc []byte
	fc, e = ioutil.ReadFile(file)
	var ui []map[string]string
	if e == nil {
		fs := string(fc)
		word := "shortProfiles:"
		n := strings.LastIndex(fs, word)
		fs = fs[n+len(word):]
		b, i := true, 0
		for cb := 0; b && i != len(fs); i++ {
			if fs[i] == '{' {
				cb = cb + 1
			} else if fs[i] == '}' {
				cb = cb - 1
			}
			b = cb != 0
		}
		if !b {
			// found a matching '}' at i with the '{' at fs[0]
			fs = fs[:i]
			ui, e = parseMap(fs)
		} else {
			e = fmt.Errorf("Not found maching curly bracket")
		}
	}
	if e == nil {
		ns = make([]string, len(ui))
		i := 0
		for _, v := range ui {
			ns[i], i = v["name"], i+1
		}
	}
	return
}

func parseMap(s string) (m []map[string]string, e error) {
	i := 1
	m = make([]map[string]string, 0)
	for e == nil && i < len(s)-1 {
		_, i, e = getID(s, i)
		var ui map[string]string
		if e == nil {
			ui, i, e = getUsrInf(s, i)
		}
		if e == nil {
			m = append(m, ui)
		}
	}

	return
}

func getID(s string, i int) (id string, n int, e error) {
	b := s[i] == '"'
	n, id = i+1, ""
	for b && n != len(s) {
		b = s[n] != '"'
		if b {
			id, n = id+string(s[n]), n+1
		}
	}
	if !b && n+2 < len(s) && s[n+1] == ':' && s[n+2] == '{' {
		n = n + 2
	} else {
		e = fmt.Errorf("Error parsing ID at %d", n)
	}
	// s[n] = '{'
	return
}

func getUsrInf(s string, i int) (ui map[string]string, n int, e error) {
	// s[i] = '{'
	n = i + 1
	keys := []string{
		"id:",
		"name:",
		"firstName:",
		"vanity:",
		"thumbSrc:",
		"uri:",
		"gender:",
		"i18nGender:",
		"additionalName:",
		"type:",
		"is_friend:",
		"is_birthday:",
		"mThumbSrcSmall:",
		"mThumbSrcLarge:",
		"dir:",
		"searchTokens:",
		"alternateName:",
		"is_nonfriend_messenger_contact:",
	}
	ui = make(map[string]string)
	for j := 0; e == nil && j != len(keys); j++ {
		// fmt.Printf("len(keys): %d, j: %d\n", len(keys), j)

		if s[n:n+len(keys[j])] == keys[j] && n < len(s) {
			n = n + len(keys[j])
			// fmt.Printf("len(s): %d, n: %d\n", len(s), n)
			val := ""
			if keys[j] != "searchTokens:" {
				for n != len(s) && s[n] != ',' && s[n] != '}' {
					val, n = val+string(s[n]), n+1
				}
			} else {
				for n != len(s) && s[n] != ']' {
					val, n = val+string(s[n]), n+1
				}
				val, n = val+"]", n+1
			}
			if n == len(s) {
				e = fmt.Errorf("Error reading %s at %d", keys[j], i)
			} else {
				ui[keys[j][:len(keys[j])-1]] = val
				n = n + 1
			}
		}
	}
	if e == nil {
		if n < len(s) && (s[n] == ',' || s[n] == '}') {
			n = n + 1
		} else {
			println(s[n-10 : n+10])
			e = fmt.Errorf("Not found next user info at %d", n)
		}
	}
	return
}

type usrInf struct {
	ID                          string   `json:"id"`
	Name                        string   `json:"name"`
	FirstName                   string   `json:"firstName"`
	Vanity                      string   `json:"vanity"`
	ThumbSrc                    string   `json:"thumbSrc"`
	URI                         string   `json:"uri"`
	Gender                      int      `json:"gender"`
	I18NGender                  int      `json:"i18nGender"`
	Type                        string   `json:"type"`
	IsFriend                    bool     `json:"is_friend"`
	MThumbSrcSmall              []string `json:"mThumbSrcSmall"`
	MThumbSrcLarge              []string `json:"mThumbSrcLarge"`
	Dir                         string   `json:"dir"`
	SearchTokens                []string `json:"searchTokens"`
	AlternateName               string   `json:"alternateName"`
	IsNonfriendMessengerContact bool     `json:"is_nonfriend_messenger_contact"`
}
