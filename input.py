import json

input_file = "prompt.json"
output_file ="decoded.json"

try:
    # 读取JSON文件（必须指定utf-8编码）
    with open(input_file, "r", encoding="utf-8") as f:
        decoded_data = json.load(f)  # 自动解码为Python对象（字典/列表）

except json.JSONDecodeError as e:
    print(f"JSON解析错误：{e}")
except FileNotFoundError:
    print(f"文件 {input_file} 不存在")
    
try:
    # 写入格式化后的JSON（缩进美化）
    with open(output_file, "w", encoding="utf-8") as f:
        json.dump(decoded_data, f, ensure_ascii=False, indent=2)
    print(f"已保存至 {output_file}")
except IOError:
    print("文件写入失败")