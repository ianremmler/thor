package thor

import (
	"io"
	"strconv"

	"github.com/skelterjohn/gopp"
)

const grammar = `
ignore: /^#.*\n/
ignore: /^(?:[ \t\n])+/

Doc => {type=Node} {field=Key} {/} <Table>
Table => {field=Type} {table} {field=Kids} <<Node>>*
Node => {type=Node} {field=Key} <id> ':' <Value>
Node => {type=Node} <Value>
Value => {type=Node} '{' <Table> '}'
Value => {type=Node} {field=Type} {num} {field=Val} <num>
Value => {type=Node} {field=Type} {bool} {field=Val} <bool>
Value => {type=Node} {field=Type} {str} {field=Val} <str>

num = /([-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?)/
bool = /(true|false)/
str = /"((?:[^"\\]|\\.)*)"/
id = /([\pL][\pL\pN\-_]*)/
`

var (
	decFact *gopp.DecoderFactory
)

func init() {
	var err error
	decFact, err = gopp.NewDecoderFactory(grammar, "Doc")
	if err != nil {
		panic(err)
	}
	decFact.RegisterType(Node{})
}

type Node struct {
	Key  string
	Val  string
	Type string
	Kids []Node
}

func (n Node) String() string {
	s := ""
	switch n.Type {
	case "table":
		for i := range n.Kids {
			s += n.Kids[i].String()
			if i < len(n.Kids)-1 {
				s += " "
			}
		}
	case "str":
		s = strconv.Quote(n.Val)
	default:
		s = n.Val
	}
	if n.Type == "table" && n.Key != "/" {
		s = "{" + s + "}"
	}
	if n.Key != "" && n.Key != "/" {
		s = n.Key + ":" + s
	}
	return s
}

func Parse(in io.Reader) (*Node, error) {
	dec := decFact.NewDecoder(in)
	tree := &Node{}
	err := dec.Decode(tree)
	return tree, err
}