package nodenormal

import (
	"time"

	"github.com/fananchong/go-xserver/common"
	"github.com/fananchong/go-xserver/common/utils"
	"github.com/fananchong/go-xserver/internal/protocol"
	"github.com/fananchong/go-xserver/internal/utility"
)

// Node : 普通节点
type Node struct {
	*Session
	components []utils.IComponent
}

// NewNode : 普通节点实现类的构造函数
func NewNode(ctx *common.Context, nodeType common.NodeType) *Node {
	node := &Node{
		Session: NewSession(ctx),
	}
	node.Info = &protocol.SERVER_INFO{}
	node.Info.Id = utility.NodeID2ServerID(utility.NewNID())
	node.Info.Type = uint32(nodeType)
	node.Info.Addrs = []string{utils.GetIPInner(ctx), utils.GetIPOuter(ctx)}
	node.Info.Ports = ctx.Config.Network.Port
	// TODO: 后续支持
	// node.Info.Overload
	// node.Info.Version
	ctx.Log.Infoln("NODE ID:", utility.NodeID2UUID(node.GetID()).String())
	return node
}

// Init : 初始化节点
func (node *Node) Init() bool {
	// ping ticker
	pingTicker := utils.NewTickerHelper(node.Ctx, 5*time.Second, node.ping)

	// bind components
	node.components = []utils.IComponent{
		node.Session,
		pingTicker,
	}
	return true
}

// Start : 节点开始工作
func (node *Node) Start() bool {
	for _, v := range node.components {
		if v != nil && v.Start() == false {
			panic("")
		}
	}
	return true
}

// Close : 关闭节点
func (node *Node) Close() {
	for _, v := range node.components {
		v.Close()
	}
}

func (node *Node) ping() {
	msg := &protocol.MSG_MGR_PING{}
	node.Session.SendMsg(uint64(protocol.CMD_MGR_PING), msg)
}
