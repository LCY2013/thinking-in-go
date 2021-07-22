### 优雅地重启或停止
你想优雅地重启或停止 web 服务器吗？有一些方法可以做到这一点。

我们可以使用 [fvbock/endless](https://github.com/fvbock/endless) 来替换默认的 ListenAndServe。更多详细信息，请参阅 issue #296。

```text
router := gin.Default()
router.GET("/", handler)
// [...]
endless.ListenAndServe(":4242", router)
```

替代方案:
- [manners](https://github.com/braintree/manners) ：可以优雅关机的 Go Http 服务器。

- [graceful](https://github.com/tylerb/graceful) ：Graceful 是一个 Go 扩展包，可以优雅地关闭 http.Handler 服务器。

- [grace](https://github.com/braintree/manners) ：Go 服务器平滑重启和零停机时间部署。









