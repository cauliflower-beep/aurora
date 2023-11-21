const net = require('net');

console.log("client starting...");
setTimeout(() => {
    // 创建一个TCP socket
    const clientSocket = new net.Socket();

    // 连接到远程服务器
    clientSocket.connect(8999, '127.0.0.1', () => {
        console.log('Connected to server');

        // 数据交互循环
        setInterval(() => {
            // 发送数据到服务器
            clientSocket.write('hello aurora v0.2! ');

            // 从服务器接收数据
            clientSocket.on('data', (data) => {
                console.log(`Server callback: ${data}`);
            });

            // 模拟处理其他事情，阻塞一秒钟
            setTimeout(() => {}, 1000);
        }, 1000);
    });

    // 处理连接错误
    clientSocket.on('error', (err) => {
        console.error(`Client start error: ${err.message}`);
        clientSocket.end();
    });

    // 处理连接关闭
    clientSocket.on('close', () => {
        console.log('Connection closed');
    });
}, 1000);
