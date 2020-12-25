package controller

import (
	"encoding/json"
	"im/model"
	"im/service"
	"log"
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
	DataQueue chan []byte //并行数据转串行
	GroupSets set.Interface
}

var rwlocker sync.RWMutex                                  //读写锁
var userNodeMap map[int64]*Node = make(map[int64]*Node, 0) //映射关系表
var udpsendchan chan []byte = make(chan []byte, 1024)      //存放发送的要广播的数据

// 通过 udp 将消息广播到局域网
// func broadMsg(data []byte) {
// 	udpsendchan <- data
// }

func checkToken(userId int64, token string) bool {
	user := service.GetUserByID(userId)
	return user.Token == token
}

// ws://127.0.0.1/im?id=1&token=xxxx
func Im(writer http.ResponseWriter, request *http.Request) {
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
	groupIds := service.SearchGroupsIDs(userId)
	for _, v := range groupIds {
		node.GroupSets.Add(v)
	}
	// userid和node形成绑定关系
	rwlocker.Lock()
	userNodeMap[userId] = node
	rwlocker.Unlock()
	go sendproc(node) // 完成发送逻辑
	go recvproc(node) // 完成接收逻辑
	log.Printf("<-%d\n", userId)
	sendMsg(userId, []byte("hello,world!"))
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
		// 返回3个参数：type(int), buffer([]byte), err
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

//后端调度逻辑处理
func dispatch(data []byte) {
	msg := model.Message{}
	// 解析data为message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// 根据cmd对逻辑进行处理
	switch msg.Cmd {
	case CMD_SINGLE_MSG:
		sendMsg(msg.ObjectId, data)
	case CMD_ROOM_MSG:
		for userId, node := range userNodeMap {
			if node.GroupSets.Has(msg.ObjectId) {
				//自己排除,不发送
				if msg.UserId != userId {
					node.DataQueue <- data
				}
			}
		}
	case CMD_HEART:
		// 一般啥都不做
	}
}

func sendMsg(userId int64, msg []byte) {
	rwlocker.RLock()
	node, ok := userNodeMap[userId]
	rwlocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

// func init() {
// 	go udpsendproc()
// 	go udprecvproc()
// }

// // udp数据的发送协程
// func udpsendproc() {
// 	log.Println("start udpsendproc")
// 	// 使用udp协议拨号
// 	udpcon, err := net.DialUDP("udp", nil, &net.UDPAddr{
// 		IP:   net.IPv4(192, 168, 0, 255),
// 		Port: 3000,
// 	})
// 	defer udpcon.Close()
// 	if err != nil {
// 		log.Println(err.Error())
// 		return
// 	}
// 	// 通过得到的con发送消息
// 	// udpcon.Write()
// 	for {
// 		select {
// 		case data := <-udpsendchan:
// 			_, err = udpcon.Write(data)
// 			if err != nil {
// 				log.Println(err.Error())
// 				return
// 			}
// 		}
// 	}
// }

// // upd数据的接收处理协程
// func udprecvproc() {
// 	log.Println("start udprecvproc")
// 	// 监听udp广播端口
// 	udpcon, err := net.ListenUDP("udp", &net.UDPAddr{
// 		IP:   net.IPv4zero,
// 		Port: 3000,
// 	})
// 	defer udpcon.Close()
// 	if err != nil {
// 		log.Println(err.Error())
// 	}
// 	// 处理端口发过来的数据
// 	for {
// 		var buf [512]byte
// 		n, err := udpcon.Read(buf[0:])
// 		if err != nil {
// 			log.Println(err.Error())
// 			return
// 		}
// 		//直接数据处理
// 		dispatch(buf[0:n])
// 	}
// 	log.Println("stop updrecvproc")
// }
