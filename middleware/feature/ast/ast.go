package ast

import (
	"fmt"
	"go/ast"
	"reflect"
)

type printVisitor struct {
}

func (p *printVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		fmt.Println(node)
		return p
	}
	val := reflect.ValueOf(node)
	typ := reflect.TypeOf(node)
	if typ.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}
	fmt.Printf("val: %+v, type: %s\n", val.Interface(), typ.Name())
	return p
}
