### Set Key Value
```text
redis> SET mykey Hello 
"OK"

等价于下面

*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nHello\r\n

*3
$3
SET 
$5
mykey 
$5
Hello

```
