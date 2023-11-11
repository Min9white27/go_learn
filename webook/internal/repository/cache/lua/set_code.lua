-- 你的验证码子在 Redis 上的 key
-- phone_code:login:15200000000
local key = KEYS[1]
-- 验证次数，我们一个验证码，最多重复三次，这个记录还可以验证几次
-- phone_code:login:15200000000:cnt
local cntKey = key..":cnt"
-- 你的验证码 123456
local val = ARGV[1]
-- 过期时间
local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    --    key 存在，但是没有过期时间
    return -2
    --    540 为 9 分钟 600 - 60
elseif ttl == -2 or ttl < 540 then
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    --    完美符合预期
    return 0
else
    --    发送太频繁
    return -1
end