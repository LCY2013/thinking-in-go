package main

/*
对于真正意外的情况，那些表示不可恢复的程序错误，例如索引越界、不可恢复的环境问题、栈溢出，我们才使用 panic。对于其他的错误情况，我们应该是期望使用 error 来进行判定。
you only need to check the error value if you care about the result.  -- Dave
This blog post from Microsoft’s engineering blog in 2005 still holds true today, namely:
My point isn’t that exceptions are bad. My point is that exceptions are too hard and I’m not smart enough to handle them.

简单。
考虑失败，而不是成功(Plan for failure, not success)。
没有隐藏的控制流。
完全交给你来控制 error。
Error are values。
*/
