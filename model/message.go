package model

type Message struct {
	Id       int64  `form:"id" json:"id,omitempty"`             //消息ID
	UserId   int64  `form:"userid" json:"userid,omitempty"`     //谁发的
	ObjectId int64  `form:"objectid" json:"objectid,omitempty"` //对端用户ID/群ID
	Cmd      int    `form:"cmd" json:"cmd,omitempty"`           //私聊/群聊心跳
	Type     int    `form:"type" json:"type,omitempty"`         //消息类型
	Content  string `form:"content" json:"content,omitempty"`   //消息的内容
	Image    string `form:"image" json:"image,omitempty"`       //预览图片
	Url      string `form:"url" json:"url,omitempty"`           //服务的URL
	Amount   int    `form:"amount" json:"amount,omitempty"`     //其他和数字相关的
}

/**
消息发送结构体
1、MEDIA_TYPE_TEXT
{id:1,userid:2,objectid:3,cmd:10,type:1,content:"hello"}
2、MEDIA_TYPE_News
{id:1,userid:2,objectid:3,cmd:10,type:2,content:"标题",image:"http://www.baidu.com/a/log,jpg",url:"http://www.a,com/dsturl","memo":"这是描述"}
3、MEDIA_TYPE_VOICE，amount单位秒
{id:1,userid:2,objectid:3,cmd:10,type:3,url:"http://www.a,com/dsturl.mp3",anount:40}
4、MEDIA_TYPE_IMG
{id:1,userid:2,objectid:3,cmd:10,type:4,url:"http://www.baidu.com/a/log,jpg"}
5、MEDIA_TYPE_REDPACKAGR //红包amount 单位分
{id:1,userid:2,objectid:3,cmd:10,type:5,url:"http://www.baidu.com/a/b/c/redpackageaddress?id=100000","amount":300,"memo":"恭喜发财"}
6、MEDIA_TYPE_EMOJ 6
{id:1,userid:2,objectid:3,cmd:10,type:6,"content":"cry"}
7、MEDIA_TYPE_Link 6
{id:1,userid:2,objectid:3,cmd:10,type:7,"url":"http://www.a,com/dsturl.html"}
8、MEDIA_TYPE_VIDEO 8
{id:1,userid:2,objectid:3,cmd:10,type:8,image:"http://www.baidu.com/a/log,jpg",url:"http://www.a,com/a.mp4"}
9、MEDIA_TYPE_CONTACT 9
{id:1,userid:2,objectid:3,cmd:10,type:9,"content":"10086","image":"http://www.baidu.com/a/avatar,jpg","memo":"胡大力"}
*/
