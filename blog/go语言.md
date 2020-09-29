#### verb(动词)

```
%d                  int变量
%x, %o, %b          分别为16进制，8进制，2进制形式的int
%f, %g, %e          浮点数: 3.141593 3.141592653589793 3.141593e+00 
%t                  布尔变量:true 或 false
%c                  rune (Unicode码点)，Go语言里特有的Unicode字符类型
%s                  string
%q                  带双引号的字符串 "abc" 或 带单引号的 rune 'c'
%v                  会将任意变量以易读的形式打印出来
%T                  打印变量的类型
%%                  字符型百分比标志(%符号本身，没有其他操作)
```

#### 关键字,总计25个

```
break
case
chan
const
continue
default 
func
defer         
go
else          
goto
fallthrough   
if
for           
import
interface
map
package
range
return
select
struct
switch
type
var
```

#### 内建的常量、类型和函数

```
内建常量: true false iota nil

内建类型:
int int8 int16 int32 int64
uint uint8 uint16 uint32 uint64 uintptr
float32 float64 complex128 complex64
bool byte rune string error

内建函数:
make len cap new append copy close delete
complex real imag
panic recover
```

