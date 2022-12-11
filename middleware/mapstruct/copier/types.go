package copier

type BasicSrc struct {
	Name    string
	Age     int
	CNumber complex64
}

type BasicDst struct {
	Name    string
	Age     int
	CNumber complex64
}

type SimpleSrc struct {
	Name    string
	Age     *int
	Friends []string
}

type SimpleDst struct {
	Name    string
	Age     *int
	Friends []string
}

type EmbedSrc struct {
	SimpleSrc
	*BasicSrc
}

type EmbedDst struct {
	SimpleSrc
	*BasicSrc
}

type ComplexSrc struct {
	Simple SimpleSrc
	Embed  *EmbedSrc
	BasicSrc
}

type ComplexDst struct {
	Simple SimpleDst
	Embed  *EmbedDst
	BasicSrc
}

type SpecialSrc struct {
	Arr [3]float32
	M   map[string]int
}

type SpecialDst struct {
	Arr [3]float32
	M   map[string]int
}

type InterfaceSrc interface {
}

type InterfaceDst interface {
}

type NotMatchSrc struct {
	Simple SimpleSrc
	Embed  *EmbedSrc
	BasicSrc
	S struct {
		A string
	}
}

type NotMatchDst struct {
	Simple SimpleDst
	Embed  *EmbedDst
	BasicSrc
	S struct {
		A int
	}
}

type MultiPtrSrc struct {
	Name    string
	Age     **int
	Friends []string
}

type MultiPtrDst struct {
	Name    string
	Age     **int
	Friends []string
}

type DiffSrc struct {
	A string
	B int
	c SimpleSrc
	F BasicSrc
}
type DiffDst struct {
	A string
	B int
	d SimpleSrc
	G BasicSrc
}

type SimpleEmbedDst struct {
	SimpleSrc
}

type ArraySrc struct {
	A []SimpleSrc
}

type ArrayDst struct {
	A []SimpleSrc
}

type ArrayDst1 struct {
	A []SimpleDst
}

type MapSrc struct {
	A map[string]SimpleSrc
}

type MapDst struct {
	A map[string]SimpleSrc
}

type MapDst1 struct {
	A map[string]SimpleDst
}

type SpecialSrc1 struct {
	A int
}

type aliasInt int
type SpecialDst1 struct {
	A aliasInt
}

type aliasInt1 = int
type SpecialDst2 struct {
	A aliasInt1
}
