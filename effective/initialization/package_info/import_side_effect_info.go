/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2021-06-11
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package package_info

import _ "net/http/pprof"

/*
Import for side effect
为副作用而导入

像前例中 fmt 或 io 这种未使用的导入总应在最后被使用或移除: 空白赋值会将代码标识为工 作正在进行中。
但有时导入某个包只是为了其副作用， 而没有任何明确的使用。
例如，在 net/http/pprof 包的 init 函数中记录了 HTTP 处理程序的调试信息。
它有个可导出的 API， 但 大部分客户端只需要该处理程序的记录和通过 Web 叶访问数据。
只为了其副作用来导入该包， 只需将包重命名为空白标识符:

这种导入格式能明确表示该包是为其副作用而导入的，因为没有其它使用该包的可能:
在此 文件中，它没有名字。(若它有名字而我们没有使用，编译器就会拒绝该程序。)
*/
