package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/einride/nobil-go/cmd/attributegen/internal/attributegen"
)

func main() {
	log.SetFlags(0)
	response, err := http.Get("https://nobil.no/admin/attributes.php")
	if err != nil {
		log.Panic(err)
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		_ = response.Body.Close()
		log.Fatal(err)
	}
	_ = response.Body.Close()
	body := string(data)
	var page attributegen.AttributesPage
	if err := page.UnmarshalHTML(body); err != nil {
		log.Fatal(err)
	}
	fmt.Print("package nobil\n")
	fmt.Println()
	fmt.Println("// Code generated by attributegen. DO NOT EDIT.")
	fmt.Println()
	fmt.Println("import (")
	fmt.Println("\t\"fmt\"")
	fmt.Println("\t\"strconv\"")
	fmt.Println(")")
	fmt.Println()
	fmt.Println("const (")
	for _, attribute := range page.Attributes {
		fmt.Print("\tAttributeID_", attribute.GoIdent(), " AttributeID = ", strconv.Quote(strconv.Itoa(attribute.ID)), "\n")
	}
	fmt.Println(")")
	for _, attribute := range page.Attributes {
		if len(attribute.Values) < 2 {
			continue
		}
		fmt.Println()
		fmt.Println("type", attribute.GoIdent(), " int")
		fmt.Println()
		fmt.Println("const (")
		for _, value := range attribute.Values {
			ident := attribute.GoIdent() + "_" + value.GoIdent()
			fmt.Print("\t", ident, " ", attribute.GoIdent(), " = ", value.ID, "\n")
		}
		fmt.Println(")")
		fmt.Println()
		fmt.Print("func (a ", attribute.GoIdent(), ") String() string {\n")
		fmt.Println("\tswitch a {")
		for _, value := range attribute.Values {
			ident := attribute.GoIdent() + "_" + value.GoIdent()
			fmt.Print("\tcase ", ident, ":\n")
			fmt.Print("\t\t", `return "`, value.Name, `"`, "\n")
		}
		fmt.Println("\tcase -1:")
		fmt.Print("\t\t", `return ""`, "\n")
		fmt.Println("\tdefault:")
		fmt.Print("\t\t", `return fmt.Sprintf("`, attribute.GoIdent(), `(%d)", a)`, "\n")
		fmt.Println("\t}")
		fmt.Println("}")
		fmt.Println()
		attributesType := "map[AttributeID]*Attribute"
		fmt.Print("func (a *", attribute.GoIdent(), ") unmarshalAttributes(attrs ", attributesType, ") {\n")
		fmt.Print("\tattr, ok := attrs[", strconv.Quote(strconv.Itoa(attribute.ID)), "]\n")
		fmt.Println("\tif !ok {")
		fmt.Println("\t\t*a = -1")
		fmt.Println("\t\treturn")
		fmt.Println("\t}")
		fmt.Print("\tv, err := strconv.Atoi(string(attr.ValueID))\n")
		fmt.Println("\tif err != nil {")
		fmt.Println("\t\t*a = -1")
		fmt.Println("\t\treturn")
		fmt.Println("\t}")
		fmt.Print("\t*a = ", attribute.GoIdent(), "(v)\n")
		fmt.Println("}")
	}
}
