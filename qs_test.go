package qs

import (
	"fmt"
	"testing"
)

func Test_parse(t *testing.T) {

	qs := "user.name=Alice&user[name]=Jack&user.address.city=Seoul&user.address.zip=07300&tags=admin,editor&string=foo&string=bar"
	root, err := Parse(qs)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", root)

	user, _ := root.GetChild("user")
	name, _ := user.GetChild("name")
	fmt.Printf("%#v\n", name.Values())

	address, _ := user.GetChild("address")
	city, _ := address.GetChild("city")
	fmt.Printf("%#v\n", city.Values())

	zip, _ := address.GetChild("zip")
	fmt.Printf("%#v\n", zip.Values())

	tags, _ := root.GetChild("tags")
	fmt.Printf("%#v\n", tags.Values())

	string_, _ := root.GetChild("string")
	fmt.Printf("%#v\n", string_.Values())

}
