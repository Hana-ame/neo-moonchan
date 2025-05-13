import sys
import requests
from io import BytesIO
import os
import json
from PIL import Image, ImageFilter, ImageOps
from PIL import ImageEnhance

input_url = sys.argv[1]
# print(input_url)
def main(image):
    image = Image.open("input.jpg").convert("L")  # 转为灰度图
    filtered_image = image.filter(ImageFilter.MedianFilter(size=3))  # 窗口大小为3x3
    edge_image = filtered_image.filter(ImageFilter.FIND_EDGES)
    enhancer = ImageEnhance.Contrast(edge_image)
    edge_image = enhancer.enhance(2.0)  # 对比度增强2倍
    inverted_image = ImageOps.invert(edge_image)
    # inverted_image.save("output.jpg")
    return inverted_image



def process_and_upload(network_input_url, api_upload_url):
    try:
        # 1. 从网络读取输入图片
        response = requests.get(network_input_url)
        response.raise_for_status()  # 检查HTTP错误[1,5](@ref)
        
        # 2. 图像处理（示例：灰度化）
        with Image.open(BytesIO(response.content)) as img:
            # img_gray = img.convert("L")  # 转为灰度图[1,3](@ref)
            img_gray = main(img)
            
            # 3. 保存到内存缓冲区
            buffer = BytesIO()
            img_gray.save(buffer, format="JPEG")  # 格式根据API要求调整[5](@ref)
            buffer.seek(0)
            
        # 4. 通过PUT上传到API
        headers = {"Content-Type": "image/jpeg"}  # 必须匹配实际格式[7](@ref)
        upload_response = requests.put(
            api_upload_url,
            data=buffer.getvalue(),
            headers=headers
        )
        upload_response.raise_for_status()
        
        # 5. 解析并打印endpoint
        result = upload_response.json()
        # print(result)
        print(f"https://upload.moonchan.xyz/api/{result['id']}/image.jpg")  # 根据API响应字段调整
        
    except requests.exceptions.RequestException as e:
        print(f"网络错误: {e}")
    except Exception as e:
        print(f"处理失败: {e}")

# 示例调用
process_and_upload(
    network_input_url=input_url,
    api_upload_url="https://upload.moonchan.xyz/api/upload"
)
