#!/bin/bash

# 获取源码路径
SOURCE_PATH=$1

# 获取可选的输出路径
OUTPUT_PATH=${2:-$(pwd)}

echo "OUTPUT_PATH: " $OUTPUT_PATH
# 检查源码路径是否提供
if [ -z "$SOURCE_PATH" ]; then
    echo "Usage: $0 <source_path> [output_path]"
    exit 1
fi

# 检查源码路径是否存在
if [ ! -d "$SOURCE_PATH" ]; then
    echo "Error: Source path '$SOURCE_PATH' does not exist."
    exit 1
fi

# 检查输出路径是否存在，如果不存在，则创建它
if [ ! -d "$OUTPUT_PATH" ]; then
    mkdir -p "$OUTPUT_PATH"
    if [ $? -ne 0 ]; then
        echo "Failed to create output directory '$OUTPUT_PATH'."
        exit 1
    fi
fi

# 构建输出文件名，考虑到操作系统
OUTPUT_FILE="wiki-user"
if [ "$(uname)" == "Darwin" ] || [ "$(uname)" == "Linux" ]; then
    OUTPUT_FILE="$OUTPUT_PATH/wiki-user"
# shellcheck disable=SC2046
elif [ "$(uname)" == "CYGWIN" ] || [ "$(uname)" == "MINGW" ] || [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
    OUTPUT_FILE="$OUTPUT_PATH/wiki-user.exe"
fi

# 执行 Go 编译
# shellcheck disable=SC2164
cd "$SOURCE_PATH"

# 拉取项目依赖
echo "Fetching dependencies..."
go mod tidy

# 如果可执行文件已存在，则删除
EXECUTABLE=$OUTPUT_PATH/wiki-user
if [ -f "$EXECUTABLE" ]; then
    echo "Removing existing executable $EXECUTABLE"
    rm $EXECUTABLE
fi

# 编译项目
echo "Building the project..."
go build -o "$OUTPUT_FILE"

# 拷贝.yaml文件到输出目录，如果文件不存在则拷贝，如果存在则跳过
echo "Copying YAML configuration files..."
find $SOURCE_PATH -name "*.yaml" | while read yaml_file; do
    base_name=$(basename "$yaml_file")
    if [ ! -f "$OUTPUT_PATH/$base_name" ]; then
         echo "cp   $yaml_file   $OUTPUT_PATH"
         cp "$yaml_file" "$OUTPUT_PATH/" || echo "Failed to copy $yaml_file to $OUTPUT_PATH. Check if the filesystem is write-protected or use sudo."
    else
        echo "Skipping existing file $base_name at destination."
    fi
done

# 检查编译是否成功
if [ $? -eq 0 ]; then
    echo "Build successful. Output at: $OUTPUT_FILE"
else
    echo "Build failed."
    exit 1
fi
echo "Listing all relevant files in the output directory:"


ls -l $OUTPUT_PATH | grep -E '(\.yaml$)'