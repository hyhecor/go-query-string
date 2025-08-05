package qs

import (
	"iter"
	"net/url"
	"strings"
)

type Node struct {
	parent   *Node
	children []*Node
	name     string
	values   []string
}

func NewNode(parent *Node, name string) *Node {
	return &Node{
		parent: parent,
		name:   name,
	}
}

func (n Node) Parent() *Node {
	return n.parent
}

func (n Node) Name() string {
	return n.name
}

func (n Node) Children() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		for _, cn := range n.children {
			if !yield(cn) {
				return
			}
		}
	}
}

func (n Node) Children2() iter.Seq2[string, *Node] {
	return func(yield func(string, *Node) bool) {
		for _, cn := range n.children {
			if !yield(cn.name, cn) {
				return
			}
		}
	}
}

func (n Node) GetChild(name string) (*Node, bool) {
	for _, c := range n.children {
		if c.name == name {
			return c, true
		}
	}

	return nil, false
}

func (n *Node) SetChild(child *Node) {
	n.children = append(n.children, child)
}

func (n Node) Values() []string {
	return n.values
}

func (n *Node) AppendValues(values ...string) {
	n.values = append(n.values, values...)
}

func Parse(qs string) (*Node, error) {
	// 1. parse into flat map
	flat, err := url.ParseQuery(qs)
	if err != nil {
		return nil, err
	}

	root := NewNode(nil, "")

	for rawKey, vals := range flat {
		for _, value := range vals {

			// 2. split key on [ ], ] or .(dot)
			parts := strings.FieldsFunc(rawKey, func(r rune) bool {
				return r == '[' || r == ']' || r == '.'
			})

			cursor := root

			// 3. keep chaining child nodes
			for _, part := range parts {
				if part == "" {
					continue
				}

				if next, ok := cursor.GetChild(part); ok {
					cursor = next
				} else {
					newNode := NewNode(cursor, part)
					cursor.SetChild(newNode)
					cursor = newNode
				}
			}

			// 4. append values to current node
			cursor.AppendValues(strings.Split(value, ",")...)
		}
	}

	return root, nil
}
