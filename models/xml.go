package models

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	//OLD ... данные не изменены
	OLD uint8 = 0
	//NEW ... полностью новые даные
	NEW uint8 = 1
	//CHANGEDVALUE ... данные изменены
	CHANGEDVALUE uint8 = 2
	//CHANGEDORDER ... изменен порядок
	CHANGEDORDER uint8 = 3
	//CHANGEDVALUEORDER ... изменен порядок
	CHANGEDVALUEORDER uint8 = 4
	//DELETED ... изменен порядок
	DELETED uint8 = 5
)

//XMLTree ...
type XMLTree struct {
	Name        string
	Value       string
	ContentHash string
	Status      uint8
	root        *XMLTree
	Branches    []*XMLTree
}

// ReadXMLFromFile ...
func ReadXMLFromFile(path string) *XMLTree {
	bytes, err := ioutil.ReadFile(path)
	CheckError(err)
	return ReadXML(bytes)
}

//ReadXML ...
func ReadXML(xmlBytes []byte) *XMLTree {
	r := strings.NewReader(string(xmlBytes))
	var depth int = 0
	xmlTree := new(XMLTree)
	nextXMLTree := xmlTree
	parser := xml.NewDecoder(r)
	for {
		token, err := parser.Token()
		if err != nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			elmt := xml.StartElement(t)
			name := elmt.Name.Local
			for _, attr := range elmt.Attr {
				name += fmt.Sprintf(" %s=\"%s\"", attr.Name.Local, attr.Value)
			}
			if depth == 0 {
				nextXMLTree.Name = name
				depth = depth + 1
				break
			}
			newBranch := new(XMLTree)
			newBranch.Name = name
			newBranch.root = nextXMLTree
			nextXMLTree.Branches = append(nextXMLTree.Branches, newBranch)
			nextXMLTree = newBranch
			depth = depth + 1

		case xml.EndElement:
			Value := nextXMLTree.Value
			for _, branch := range nextXMLTree.Branches {
				Value += branch.ContentHash
			}
			contentBytes := []byte(Value)
			contentHash := fmt.Sprintf("%x", md5.Sum(contentBytes)) //convert bytes array to HEXstring
			nextXMLTree.ContentHash = contentHash
			depth = depth - 1
			nextXMLTree = nextXMLTree.root
		case xml.CharData:
			bytes := xml.CharData(t)
			chardata := string(bytes)
			if chardata != "\n" {
				nextXMLTree.Value += chardata
			}
			/*
				case xml.Comment:
					fmt.Println("comment")
				case xml.ProcInst:
					fmt.Println("ProcInst")
				case xml.Directive:
					fmt.Println("Directive")
				default:
					fmt.Println("Unknown")
			*/
		}
	}
	return xmlTree
}

//PrintXML ...
func PrintXML(xml XMLTree, depth int) {
	PrintElmt(xml.Name, depth, xml.Status)
	PrintElmt("{", depth, 0)
	if xml.Value != "\n" && xml.Value != "" {
		PrintElmt("'"+xml.Value+"'", depth, 0)
	}
	if len(xml.Branches) != 0 {
		depth++
		for _, nextElem := range xml.Branches {
			PrintXML(*nextElem, depth)
		}
		depth--
	}
	PrintElmt("}", depth, 0)
	depth--
}

//PrintElmt ...
func PrintElmt(s string, depth int, status uint8) {
	var add string
	switch status {
	case OLD:
		add = ""
	case NEW:
		add = "NEW| "
	case CHANGEDVALUE:
		add = "CHANGEDVALUE| "
	case CHANGEDORDER:
		add = "CHANGEDORDER| "
	case CHANGEDVALUEORDER:
		add = "CHANGEDVALUEORDER| "
	case DELETED:
		add = "DELETED| "
	}
	for n := 0; n < depth; n++ {
		fmt.Print("|  ")
	}
	fmt.Println(add + s)
}

//CheckError ...
func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
