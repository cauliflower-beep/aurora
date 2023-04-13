#!/bin/bash
protoc --go_out=. *.proto

# 加上下面三行 可以保持窗口不闪退，这样就有时间看命令行日志了
echo 按任意键继续
read -n 1
echo 继续运行