package main

import (
	"strings"
	"testing"
	"text/scanner"
)

var src = `shortProfiles:{"100014259181657":{id:"100014259181657",name:"Cocuya González",firstName:"Cocuya",vanity:"Cocuya.González",thumbSrc:"https://scontent-mia3-1.xx.fbcdn.net/v/t1.0-1/p32x32/44837884_499207077231261_9161212237900152832_n.jpg?_nc_cat=100&_nc_ht=scontent-mia3-1.xx&oh=712b9f8fc8d80979a34c6516e386a3d4&oe=5C7F3BC2",uri:"https://www.facebook.com/Cocuya.González",gender:1,i18nGender:2,type:"friend",is_friend:true,mThumbSrcSmall:null,mThumbSrcLarge:null,dir:null,searchTokens:["González","Cocuya"],alternateName:"",is_nonfriend_messenger_contact:false}`

func TestScan(t *testing.T) {
	rd := strings.NewReader(src)
	s := new(scanner.Scanner)
	s.Init(rd)
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		t.Log(s.TokenText())
	}
}

func TestScanData(t *testing.T) {
	rd := strings.NewReader(src)
	s := new(scanner.Scanner)
	s.Init(rd)
	ns, e := scanData(s)
	if e == nil {
		for i, j := range ns {
			t.Logf("%d - %s", i, j)
		}
	} else {
		t.Log(e.Error())
	}
}
