import sys
import json

def process_lines(lines: list[str]):
    actors_set = set()
    
    # 逐行读取标准输入
    for line in lines:
        print(line)
        line = line.strip()
        if not line:
            continue
        
        # 分割时间戳和JSON部分（按第一个空格分割）
        try:
            date_part, time_part, json_part = line.split(' ', 2)
        except ValueError:
            continue  # 跳过格式不正确的行
        
        # 解析JSON
        try:
            data = json.loads(json_part)
            if data.get('type') == 'Block':
                actor = data.get('actor')
                if actor:  # 确保actor字段存在且非空
                    actors_set.add(actor)
        except json.JSONDecodeError:
            pass  # 跳过无效的JSON
    
    # 输出排序后的结果
    for actor in sorted(actors_set):
        print(actor)

if __name__ == "__main__":
    with open('log.txt') as f:
        process_lines(f.readlines())
        