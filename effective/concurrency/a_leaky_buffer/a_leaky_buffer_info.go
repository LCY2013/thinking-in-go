/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2021-06-17
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package a_leaky_buffer

/*
A leaky buffer  可能泄露的缓冲区

并发编程的工具甚至能很容易地表达非并发的思想。
这里有个提取自 RPC 包的例子。
客户端 Go 程从某些来源，可能是网络中循环接收数据。
为避免分配和释放缓冲区，它保存了一个空闲链表，使用一个带缓冲信道表示。
若信道为空，就会分配新的缓冲区。 一旦消息缓冲区就绪，它将通过 serverChan 被发送到服务器。
serverChan.

var freeList = make(chan *Buffer, 100)
var serverChan = make(chan *Buffer)

func client() {
	for {
		var b *Buffer
		// 若缓冲区可用就用它，不可用就分配个新的。
		select {
		case b = <-freeList:
		// 获取一个，不做别的。
		default:
			// 非空闲，因此分配一个新的。
			b = new(Buffer)
		}
		load(b) // 从网络中读取下一条信息
		serverChan <- b // 发送至服务器。
	}
}

服务器从客户端循环接收每个消息，处理它们，并将缓冲区返回给空闲列表。
func server() {
	for {
		b := <-serverChan // 等待工作。
        process(b)
		// 若缓冲区有空间就重用它。
		select {
		case freeList <- b:
			// 将缓冲区放大空闲列表中，不做别的。
		default:
			// 空闲列表已满，保持就好。
		}
	}
}

客户端试图从 freeList 中获取缓冲区;若没有缓冲区可用，它就将分配一个新的。
服务器将 b 放回空闲列表 freeList 中直到列表已满，此时缓冲区将被丢弃，并被垃圾回收器回收。
(select 语句中的 default 子句在没有条件符合时执行，这也就意味着 selects 永远不会被阻塞。)
依靠带缓冲的信道和垃圾回收器的记录， 我们仅用短短几行代码就构建了一个可能导致缓冲区槽位泄露的空闲列表。
*/
