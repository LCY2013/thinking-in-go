module demo2

go 1.19

require (
	github.com/lcy2013/workerpool v1.0.0
)

replace (
	github.com/lcy2013/workerpool v1.0.0 => ../workerpool2
)
