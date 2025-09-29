import os
import threading
import socket
import time
from collections import OrderedDict
from concurrent.futures import ThreadPoolExecutor
import queue

class LRUTimedCache:
    def __init__(self, ttl: int = 3600, max_size: int = 1000, cleanup_interval: int = 300):
        self.cache = OrderedDict()
        self.ttl = ttl
        self.max_size = max_size
        self.cleanup_interval = cleanup_interval
        self.last_cleanup = time.time()
        self.lock = threading.Lock()
    
    def should_skip(self, username: str) -> bool:
        with self.lock:
            current_time = time.time()
            self._auto_cleanup(current_time)
            
            if username in self.cache:
                timestamp = self.cache[username]
                if current_time - timestamp <= self.ttl:
                    self.cache.move_to_end(username)
                    return True
                else:
                    del self.cache[username]
            
            self.cache[username] = current_time
            self.cache.move_to_end(username)
            
            if len(self.cache) > self.max_size:
                self.cache.popitem(last=False)
            
            return False
    
    def _auto_cleanup(self, current_time: float):
        if current_time - self.last_cleanup >= self.cleanup_interval:
            self._cleanup_expired(current_time)
            self.last_cleanup = current_time
    
    def _cleanup_expired(self, current_time: float):
        expired = []
        for username, timestamp in self.cache.items():
            if current_time - timestamp > self.ttl:
                expired.append(username)
            else:
                break
        
        for username in expired:
            del self.cache[username]

class CommandExecutor:
    def __init__(self, max_concurrent=10):
        self.max_concurrent = max_concurrent
        self.executor = ThreadPoolExecutor(max_workers=max_concurrent)
        self.pending_tasks = queue.Queue()
        self.completed_tasks = 0
        self.lock = threading.Lock()
        self.running = True
        
        # 启动任务处理线程
        self.process_thread = threading.Thread(target=self._process_queue, daemon=True)
        self.process_thread.start()
    
    def submit_command(self, username):
        """提交命令到执行器"""
        self.pending_tasks.put(username)
    
    def _process_queue(self):
        """处理队列中的任务"""
        while self.running:
            try:
                # 从队列获取用户名（阻塞式，最多等待1秒）
                username = self.pending_tasks.get(timeout=1)
                
                # 提交到线程池执行
                future = self.executor.submit(self._execute_command, username)
                future.add_done_callback(lambda f: self._task_completed(username))
                
            except queue.Empty:
                continue
            except Exception as e:
                print(f"处理任务时发生错误: {e}")
    
    def _execute_command(self, username):
        """执行单个命令"""
        try:
            print(f"[执行中] 开始处理用户: {username}")
            exit_code = os.system(f"py get_meta_data.py {username}")
            if exit_code == 0:
                print(f"[成功] 用户 {username} 处理完成")
            else:
                print(f"[失败] 用户 {username} 处理失败，退出码: {exit_code}")
            return exit_code
        except Exception as e:
            print(f"[异常] 处理用户 {username} 时发生错误: {e}")
            return -1
    
    def _task_completed(self, username):
        """任务完成回调"""
        with self.lock:
            self.completed_tasks += 1
        self.pending_tasks.task_done()
        print(f"[完成] 用户 {username} 任务已完成，总计完成: {self.completed_tasks}")
    
    def wait_completion(self):
        """等待所有任务完成"""
        self.pending_tasks.join()
        print("所有任务执行完成！")
    
    def shutdown(self):
        """关闭执行器"""
        self.running = False
        self.executor.shutdown(wait=True)

# 创建全局实例
lru_cache = LRUTimedCache(ttl=3600, max_size=500)
command_executor = CommandExecutor(max_concurrent=10)

def call_with_lru(username: str):
    """使用LRU缓存和并发控制的调用函数"""
    if lru_cache.should_skip(username):
        print(f"[跳过] 用户 {username} 在一小时内已处理")
        return
    
    print(f"[提交] 用户 {username} 已提交到执行队列")
    # 将任务提交到执行器队列
    command_executor.submit_command(username)

def handle_client(client_socket, client_address):
    """处理单个客户端连接的线程函数"""
    print(f"[连接] 客户端 {client_address} 已连接")
    try:
        # 接收客户端发送的数据
        data = client_socket.recv(1024)
        if not data:
            print(f"[空数据] 客户端 {client_address} 发送了空数据")
            return
            
        received_str = data.decode('utf-8').strip()
        print(f"[接收] 来自 {client_address} 的字符串: {repr(received_str)}")
        
        # 使用接收到的字符串执行调用
        call_with_lru(received_str)
        
    except UnicodeDecodeError:
        print(f"[解码错误] 来自 {client_address} 的数据无法解码为 UTF-8")
    except Exception as e:
        print(f"[错误] 处理客户端 {client_address} 请求时发生错误: {e}")
    finally:
        # 关闭当前客户端连接
        client_socket.close()
        print(f"[关闭] 与客户端 {client_address} 的连接已关闭")

def start_tcp_listener(host='127.25.9.19', port=8080):
    """启动TCP监听服务器"""
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    
    try:
        server_socket.bind((host, port))
        server_socket.listen(5)
        print(f"[启动] TCP 监听器已启动在 {host}:{port}")
        print("[等待] 等待客户端连接...")
        
        while True:
            client_sock, client_addr = server_socket.accept()
            client_thread = threading.Thread(target=handle_client, args=(client_sock, client_addr))
            client_thread.daemon = True
            client_thread.start()
            print(f"[活跃] 当前连接数: {threading.active_count() - 1}")
            
    except KeyboardInterrupt:
        print("\n[中断] 收到中断信号，服务器关闭中...")
    except Exception as e:
        print(f"[错误] 服务器运行出错: {e}")
    finally:
        # 等待所有任务完成
        command_executor.wait_completion()
        command_executor.shutdown()
        server_socket.close()
        print("[关闭] 服务器已关闭")

if __name__ == '__main__':
    start_tcp_listener()