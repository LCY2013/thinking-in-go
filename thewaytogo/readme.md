在格式化字符串里，%d 用于格式化整数（%x 和 %X 用于格式化 16 进制表示的数字），%g 用于格式化浮点型（%f 输出浮点数，%e 输出科学计数表示法），%0nd 用于规定输出长度为n的整数，其中开头的数字 0 是必须的。

%n.mg 用于表示数字 n 并精确到小数点后 m 位，除了使用 g 之外，还可以使用 e 或者 f，例如：使用格式化字符串 %5.2e 来输出 3.4 的结果为 3.40e+00。

```
优先级     运算符
 7         ^ !
 6         * / % << >> & &^
 5         + - | ^
 4         == != < <= >= >
 3         <-
 2         &&
 1         ||
```

```
\n：换行符
\r：回车符
\t：tab 键
\u 或 \U：Unicode 字符
\\：反斜杠自身
```
