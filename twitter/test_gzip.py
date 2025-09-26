import os
import glob
import gzip

# 获取当前目录下所有的 json 文件
json_files = glob.glob("*.json")
print("当前文件夹中的 JSON 文件有：")
for file in json_files:
    print(f"  - {file}")

# 遍历每个 json 文件并进行压缩
for json_file in json_files:
    try:
        gz_file = json_file + ".gz"  # 压缩后的文件名
        with open(json_file, "r", encoding='utf-8') as reader:
            with gzip.open(gz_file, "wt", encoding="utf-8", compresslevel=9) as writer:
                data = reader.read()
                writer.write(data)
        print(f"✅ 成功压缩: {json_file} -> {gz_file}")
    except Exception as e:
        print(f"❌ 压缩文件 {json_file} 时出错: {str(e)}")

print("处理完成！")