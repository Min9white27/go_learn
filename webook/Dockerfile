# 基础镜像
FROM ubuntu:20.04
# 把编译后的打包进来这个镜像，放到工作目录 /app。可以根据个人或公司需求，随便放
COPY webook /app/webook
WORKDIR /app
# CMD 是执行命令，是在启动webook的时候执行
# 最佳实践
# ENTRYPOINT 是在打开这个文件夹的时候就执行一遍
ENTRYPOINT ["/app/webook"]


