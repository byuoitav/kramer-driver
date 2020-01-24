package via

import (
	"encoding/xml"
	"fmt"
)

type via struct {
	Address  string
	Username string
	Password string
}

type Command struct {
	XMLName  xml.Name `xml:"P"`
	Username string   `xml:"UN"`
	Password string   `xml:"Pwd"`
	Command  string   `xml:"Cmd"`
	Param1   string   `xml:"P1,omitempty"`
	Param2   string   `xml:"P2,omitempty"`
	Param3   string   `xml:"P3,omitempty"`
	Param4   string   `xml:"P4,omitempty"`
	Param5   string   `xml:"P5,omitempty"`
	Param6   string   `xml:"P6,omitempty"`
	Param7   string   `xml:"P7,omitempty"`
	Param8   string   `xml:"P8,omitempty"`
	Param9   string   `xml:"P9,omitempty"`
	Param10  string   `xml:"P10,omitempty"`
}

func main() {
	fmt.Println("vim-go")
}
