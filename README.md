## go实现的网络框架

### 将消息处理拆解成一下几个模块，均可自定义替换任意子模块
1. 封包/解包，默认采用```长度``` + ```数据流```
2. 消息编码/解码：、
   - 编码：将要发送的消息，解析为字节流
   - 解码：将收到的字节流，解析为所需的消息
   - 目前支持基本的字符串消息、protobuf消息
3. 收消息接口：解析出的消息会传递给这个接口
4. 事件处理接口：建立连接、关闭连接等等事件将会传递到该接口，还可以自定义事件
### 优化记录
###### 1. SendMessage时，一定是调用方将message转换为[]byte数据序列，避免在其他线程转换时引发并发问题
###### 2. **socket写数据优化**。socket写数据核心流程：message序列化->写入socket，有些socket库，存在重复拷贝的问题，基于核心流程，减少不必要的拷贝和内存分配，进行如下优化：
1. **每个连接定义一个write_buffer，用来减少tcp报文数量**：假如同时发送多个消    息时，不需要每个消息都要立即写入socket，只需每个消息调用session的WriteMessage接口暂时写入到write_buffer，最后调用Flush时统一写入到socket
2. **减少message序列化时内存分配**，session中定义`cacheBufferList chan []byte`缓存数据列表，每次序列化时从其中拿一个buffer，写入到socket后返还回去，如果取不到再使用`make([]byte, xxx_size)直接分配buffer。估测90%以上的消息都能通过该流程简化（比直接从内存中make快，并且减少内存碎片）
3. `**减少内存重复拷贝**：如果是直接Send数据，那么无需拷贝到writebuffer，可直接写入到socket

###### 4. **socket读数据优化**。socket读数据核心流程：socket读数据到buffer->根据buffer序列化message
**每个连接定义一个read_buffer和read_decode_buffer**：socket读入到read_buffer，然后直接根据read_buffer中的数据解析出message，循环复用read_buffer，并且减少拷贝
###### 5. 分包策略
目前大部分解包/封包都会写入消息的长度，一般都是用int（4个字节）表示，采用[SOCKET封包和解包](https://blog.csdn.net/bestans/article/details/103188695)动态策略可以给绝大部分消息节省2/3个字节
###### 6. 减少read和write协程中select的chan数量 
由于select操作每增加一个chan的查询，会带来一定得查询消耗，为了提高read和write效率，因此建议这俩协程select的chan数量严格控制（看到有些socket库在select加上quitChan，可以去除，提高性能）
###### 7. EncodeMessage时使用sync.Pool分配writeBuffer，减少GC压力
