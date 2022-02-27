package dump

import (
	"fmt"
	"unsafe"
)

type eface struct {
	_type *_type
	data  unsafe.Pointer
}

type tflag uint8
type nameOff int32
type typeOff int32
type _type struct {
	size       uintptr
	ptrdata    uintptr // size of memory prefix holding all pointers
	hash       uint32
	tflag      tflag
	align      uint8
	fieldAlign uint8
	kind       uint8
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal func(unsafe.Pointer, unsafe.Pointer) bool
	// gcdata stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, gcdata is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	gcdata    *byte
	str       nameOff
	ptrToThis typeOff
}
type iface struct {
	tab  *itab
	data unsafe.Pointer
}
type itab struct {
	inter *interfacetype
	_type *_type
	hash  uint32 // copy of _type.hash. Used for type switches.
	_     [4]byte
	fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}

type interfacetype struct {
	typ     _type
	pkgpath name
	mhdr    []imethod
}

type imethod struct {
	name nameOff
	ityp typeOff
}

// name is an encoded type name with optional extra data.
// See reflect/type.go for details.
type name struct {
	bytes *byte
}

const ptrSize = unsafe.Sizeof(uintptr(0))

// DumpEface 用于输出空接口类型变量的内部表示信息；
func DumpEface(i interface{}) {
	ptrToEface := (*eface)(unsafe.Pointer(&i))
	fmt.Printf("eface: %+v\n", *ptrToEface)
	if ptrToEface._type != nil {
		// dump _type info
		fmt.Printf("\t _type: %+v\n", *(ptrToEface._type))
	}
	if ptrToEface.data != nil {
		// dump data
		switch i.(type) {
		case int:
			dumpInt(ptrToEface.data)
		case float64:
			dumpFloat64(ptrToEface.data)
		case T:
			dumpT(ptrToEface.data)
		// other cases ... ...
		default:
			fmt.Printf("\t unsupported data type\n")
		}
	}
	fmt.Printf("\n")
}

// DumpItabOfIface 用于输出非空接口类型变量的 tab 字段信息；
func DumpItabOfIface(ptrToIface unsafe.Pointer) {
	p := (*iface)(ptrToIface)
	fmt.Printf("iface: %+v\n", *p)
	if p.tab != nil {
		// dump itab
		fmt.Printf("\t itab: %+v\n", *(p.tab))
		// dump inter in itab
		fmt.Printf("\t\t inter: %+v\n", *(p.tab.inter))
		// dump _type in itab
		fmt.Printf("\t\t _type: %+v\n", *(p.tab._type))
		// dump fun in tab
		funPtr := unsafe.Pointer(&(p.tab.fun))
		fmt.Printf("\t\t fun: [")
		for i := 0; i < len((*(p.tab.inter)).mhdr); i++ {
			tp := (*uintptr)(unsafe.Pointer(uintptr(funPtr) + uintptr(i)*ptrSize))
			fmt.Printf("0x%x(%d),", *tp, *tp)
		}
		fmt.Printf("]\n")
	}
}

// DumpDataOfIface 用于输出非空接口类型变量的 data 字段信息；
func DumpDataOfIface(i interface{}) {
	// this is a trick as the data part of eface and iface are same
	ptrToEface := (*eface)(unsafe.Pointer(&i))
	if ptrToEface.data != nil {
		// dump data
		switch i.(type) {
		case int:
			dumpInt(ptrToEface.data)
		case float64:
			dumpFloat64(ptrToEface.data)
		case T:
			dumpT(ptrToEface.data)
		// other cases ... ...
		default:
			fmt.Printf("\t unsupported data type\n")
		}
	}
	fmt.Printf("\n")
}

func dumpFloat64(data unsafe.Pointer) {
	var p *float64 = (*float64)(data)
	fmt.Printf("\t data float: %+v\n", *p)
}

func dumpInt(data unsafe.Pointer) {
	var p *int = (*int)(data)
	fmt.Printf("\t data int: %+v\n", *p)
}

type T int

func (t T) Error() string {
	return "bad error"
}

func dumpT(dataOfIface unsafe.Pointer) {
	var p *T = (*T)(dataOfIface)
	fmt.Printf("\t data: %+v\n", *p)
}
