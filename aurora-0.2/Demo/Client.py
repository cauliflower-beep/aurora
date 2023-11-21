import socket
import time

def main():
    print("client starting...")
    time.sleep(1)

    # 创建一个TCP socket
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # 连接到远程服务器
    try:
        client_socket.connect(("127.0.0.1", 8999))
    except socket.error as e:
        print(f"Client start error: {e}")
        return

    # 数据交互循环
    while True:
        # 发送数据到服务器
        client_socket.sendall(b"hello aurora v0.2! ")

        # 从服务器接收数据
        data = client_socket.recv(512)
        if not data:
            print("No data received. Exiting.")
            break

        print(f"Server callback: {data.decode()}")

        # 模拟处理其他事情，阻塞一秒钟
        time.sleep(1)

    # 关闭连接
    client_socket.close()

if __name__ == "__main__":
    main()
