#!/bin/bash

# 设置默认的GIN_MODE，如果未提供参数则使用
GIN_MODE=${1:-"debug"}

# 设置环境变量
export GIN_MODE

# 检查是否有提供配置文件路径参数
if [ -z "$2" ]; then
    # 使用当前路径下的config.yaml文件（如果存在的话）
    CONFIG_PATH="./config.yaml"
    if [ ! -f "$CONFIG_PATH" ]; then
        echo "Configuration file not found in the default path, running without a configuration file..."
        CONFIG_PATH=""
    fi
else
    # 使用提供的路径
    CONFIG_PATH=$2
    # 检查配置文件是否存在
    if [ ! -f "$CONFIG_PATH" ]; then
        echo "Configuration file '$CONFIG_PATH' not found. Exiting..."
        exit 1
    fi
fi

# 根据配置文件路径是否设置来运行应用
if [ -z "$CONFIG_PATH" ]; then
    ./wiki-user
else
    ./wiki-user -c "$CONFIG_PATH"
fi

