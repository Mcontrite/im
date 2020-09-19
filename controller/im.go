package controller

import (
	"encoding/json"
	"im/model"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
)

const (
	CMD_SINGLE_MSG = 10
	CMD_ROOM_MSG   = 11
	CMD_HEART      = 0
)

//形成userid和Node的映射关系
type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte //并行转串行,
	GroupSets set.Interface
}

//映射关系表
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

//读写锁
var rwlocker sync.RWMutex

func init() {
	go udpsendproc()
	go udprecvproc()
}

//用来存放发送的要广播的数据
var udpsendchan chan []byte = make(chan []byte, 1024)

// 通过 udp 将消息广播到局域网
func broadMsg(data []byte) {
	udpsendchan <- data
}

// 完成udp数据的发送协程
func udpsendproc() {
	log.Println("start udpsendproc")
	// 使用udp协议拨号
	con, err := net.DialUDP("udp", nil,
		&net.UDPAddr{
			IP:   net.IPv4(192, 168, 0, 255),
			Port: 3000,
		})
	defer con.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	// 通过的到的con发送消息
	// con.Write()
	for {
		select {
		case data := <-udpsendchan:
			_, err = con.Write(data)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

// 完成upd接收并处理功能
func udprecvproc() {
	log.Println("start udprecvproc")
	// 监听udp广播端口
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		log.Println(err.Error())
	}
	// 处理端口发过来的数据
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			log.Println(err.Error())
			return
		}
		//直接数据处理
		dispatch(buf[0:n])
	}
	log.Println("stop updrecvproc")
}

//添加新的群ID到用户的groupset中
func AddGroupId(userId, gid int64) {
	//取得node
	rwlocker.Lock()
	node, ok := clientMap[userId]
	if ok {
		node.GroupSets.Add(gid)
	}
	//clientMap[userId] = node
	rwlocker.Unlock()
	//添加gid到set
}

//后端调度逻辑处理
func dispatch(data []byte) {
	// 解析data为message
	msg := model.Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// 根据cmd对逻辑进行处理
	switch msg.Cmd {
	case CMD_SINGLE_MSG:
		sendMsg(msg.Dstid, data)
	case CMD_ROOM_MSG:
		// 群聊转发逻辑
		for userId, v := range clientMap {
			if v.GroupSets.Has(msg.Dstid) {
				//自己排除,不发送
				if msg.Userid != userId {
					v.DataQueue <- data
				}
			}
		}
	case CMD_HEART:
		// 一般啥都不做
	}
}

// 发送消息
func sendMsg(userId int64, msg []byte) {
	rwlocker.RLock()
	node, ok := clientMap[userId]
	rwlocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

//ws发送协程
func sendproc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

//ws接收协程
func recvproc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			return
		}
		dispatch(data)
		//把消息广播到局域网
		//broadMsg(data)
		log.Printf("[ws]<=%s\n", data)
	}
}

//检测是否有效
func checkToken(userId int64, token string) bool {
	//从数据库里面查询并比对
	user := userService.Find(userId)
	return user.Token == token
}

// ws://127.0.0.1/chat?id=1&token=xxxx
func Chat(writer http.ResponseWriter, request *http.Request) {
	// 检验接入是否合法
	//checkToken(userId int64,token string)
	query := request.URL.Query()
	id := query.Get("id")
	token := query.Get("token")
	userId, _ := strconv.ParseInt(id, 10, 64)
	isvalida := checkToken(userId, token)
	log.Println(id, token, isvalida)
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// 获得conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	// 获取用户全部群Id
	comIds := contactService.SearchComunityIds(userId)
	for _, v := range comIds {
		node.GroupSets.Add(v)
	}
	// userid和node形成绑定关系
	rwlocker.Lock()
	clientMap[userId] = node
	rwlocker.Unlock()
	go sendproc(node) // 完成发送逻辑
	go recvproc(node) // 完成接收逻辑
	log.Printf("<-%d\n", userId)
	sendMsg(userId, []byte("hello,world!"))
}
