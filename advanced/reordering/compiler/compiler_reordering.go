package main

/*
==== 编译器重排 ====
snippet 1
X = 0
for i in range(100):

	X = 1
	print X

snippet 2
X = 1
for i in range(100):

	print X

snippet 1 和 snippet 2 是等价的。

如果这时候，假设有 Processor 2 同时在执行一条指令：

X = 0
P2 中的指令和 snippet 1 交错执行时，可能产生的结果是：111101111..

P2 中的指令和 snippet 2 交错执行时，可能产生的结果是：11100000…

有人说这个例子不够有说服力，我们看看参考资料中的另一个例子:

int a, b;
int foo()

	{
	    a = b + 1;
	    b = 0;
	    return 1;
	}

输出汇编:

mov eax, DWORD PTR b[rip]
add eax, 1
mov DWORD PTR a[rip], eax    // --> store to a
mov DWORD PTR b[rip], 0      // --> store to b
开启 O2 优化后，输出汇编:

mov eax, DWORD PTR b[rip]
mov DWORD PTR b[rip], 0      // --> store to b
add eax, 1
mov DWORD PTR a[rip], eax    // --> store to a
给 a 和 b 的赋值顺序被修改了，可见 compiler 也是可能会修改赋值的顺序的。

在多核心场景下,没有办法轻易地判断两段程序是“等价”的。
*/
func main() {

}
