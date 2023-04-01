local argc = tonumber(argc)
local argvLocal = {}

for i = 0, argc - 1 do
    table.insert(argvLocal, i+1, tostring(argv[i]))
end

-- 使用动态参数执行一些操作
-- ...

length = #argvLocal

-- 将结果传递回Go应用程序
--result = "some result: " .."".. length
result = table.concat(argvLocal, ",") ..":".. length