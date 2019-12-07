package biz

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pion/ion/pkg/db"
	"github.com/pion/ion/pkg/log"
	"github.com/pion/ion/pkg/mq"
	"github.com/pion/ion/pkg/proto"
	"github.com/pion/ion/pkg/util"
)

const (
	redisKeyTTL     = 1500 * time.Millisecond
	redisLongKeyTTL = 24 * time.Hour
)

var (
	amqp               *mq.Amqp
	redis              *db.Redis
	streamAddCache     = make(map[string]bool)
	streamAddCacheLock sync.RWMutex
)

// Init func
func Init(mqURL string, config db.Config) {
	amqp = mq.New(proto.IslbID, mqURL)
	redis = db.NewRedis(config)
	handleRPCMsgs()
	handleBroadCastMsgs()
}

func getPIDFromMID(mid string) string {
	return strings.Split(mid, "#")[0]
}

func handleRPCMsgs() {
	rpcMsgs, err := amqp.ConsumeRPC()
	if err != nil {
		log.Errorf(err.Error())
		return
	}

	go func() {
		for m := range rpcMsgs {
			msg := util.Unmarshal(string(m.Body))

			from := m.ReplyTo
			corrID := m.CorrelationId
			if from == proto.IslbID {
				continue
			}
			method := util.Val(msg, "method")
			// log.Infof("rpc from=%v method=%v msg=%v", from, method, msg)
			if method == "" {
				continue
			}
			switch method {
			case proto.IslbOnStreamAdd:
				rid := util.Val(msg, "rid")
				pid := util.Val(msg, "pid")
				mid := util.Val(msg, "mid")
				ssrcPt := util.Unmarshal(util.Val(msg, "mediaInfo"))
				key := rid + "/pub/node/" + pid
				redis.HSetTTL(key, from, "", redisKeyTTL)
				key = rid + "/pub/media/" + mid
				for ssrc, pt := range ssrcPt {
					redis.HSetTTL(key, ssrc, pt, redisKeyTTL)
				}
				//receive more than one streamAdd in 1s, only send once
				if _, ok := streamAddCache[mid]; !ok {
					streamAddCacheLock.Lock()
					streamAddCache[mid] = true
					streamAddCacheLock.Unlock()
					time.AfterFunc(1*time.Second, func() {
						streamAddCacheLock.Lock()
						delete(streamAddCache, mid)
						streamAddCacheLock.Unlock()
					})
					infoMap := redis.HGetAll(rid + "/peer/info/" + pid)
					for info := range infoMap {
						onStreamAdd := util.Map("method", proto.IslbOnStreamAdd, "rid", rid, "pid", pid, "mid", mid, "info", info)
						amqp.BroadCast(onStreamAdd)
					}
				}
			case proto.IslbKeepAlive:
				rid := util.Val(msg, "rid")
				pid := util.Val(msg, "pid")
				mid := util.Val(msg, "mid")
				ssrcPt := util.Unmarshal(util.Val(msg, "mediaInfo"))
				key := rid + "/pub/node/" + pid
				redis.HSetTTL(key, from, "", redisKeyTTL)
				key = rid + "/pub/media/" + mid
				for ssrc, pt := range ssrcPt {
					redis.HSetTTL(key, ssrc, pt, redisKeyTTL)
				}
			case proto.IslbGetPubs:
				rid := util.Val(msg, "rid")
				skipPid := util.Val(msg, "pid")
				key := rid + "/pub/media/*"
				log.Infof("key=%s skipPid=%v", key, skipPid)
				for _, k := range redis.Keys(key) {
					log.Infof("key=%s k=%s skipPid=%v", key, k, skipPid)
					pid := strings.Split(strings.Split(k, "#")[0], "/")[3]
					if !strings.Contains(k, skipPid) {
						mid := strings.Split(k, "/")[3]
						ssrcs := "{"
						for ssrc, pt := range redis.HGetAll(rid + "/pub/media/" + mid) {
							ssrcs += fmt.Sprintf("%s:%s, ", ssrc, pt)
						}
						ssrcs += "}"
						info := redis.HGetAll(rid + "/peer/info/" + pid)
						for i := range info {
							resp := util.Map("response", proto.IslbGetPubs, "rid", rid, "pid", pid, "mid", mid, "mediaInfo", ssrcs, "info", i)
							log.Infof("amqp.RpcCall from=%s resp=%v corrID=%s", from, resp, corrID)
							amqp.RpcCall(from, resp, corrID)
						}
					}
				}
			case proto.IslbOnStreamRemove:
				rid := util.Val(msg, "rid")
				mid := util.Val(msg, "mid")
				pid := getPIDFromMID(mid)
				key := rid + "/pub/media/" + mid
				redis.Del(key)
				onStreamRemove := util.Map("rid", rid, "method", proto.IslbOnStreamRemove, "pid", pid, "mid", mid)
				log.Infof("amqp.BroadCast onStreamRemove=%v", onStreamRemove)
				amqp.BroadCast(onStreamRemove)
			case proto.IslbClientOnJoin:
				rid := util.Val(msg, "rid")
				id := util.Val(msg, "id")
				info := util.Val(msg, "info")
				onJoin := util.Map("method", proto.IslbClientOnJoin, "rid", rid, "id", id, "info", info)
				key := rid + "/peer/info/" + id
				log.Infof("redis.HSetTTL %v %v", key, info)
				redis.HSetTTL(key, info, "", redisLongKeyTTL)
				log.Infof("amqp.BroadCast onJoin=%v", onJoin)
				amqp.BroadCast(onJoin)
			case proto.IslbClientOnLeave:
				rid := util.Val(msg, "rid")
				id := util.Val(msg, "id")
				key := rid + "/pub/node/" + id
				redis.Del(key)
				onLeave := util.Map("rid", rid, "method", proto.IslbClientOnLeave, "id", id)
				log.Infof("amqp.BroadCast onLeave=%v", onLeave)
				amqp.BroadCast(onLeave)
			case proto.IslbGetMediaInfo:
				rid := util.Val(msg, "rid")
				mid := util.Val(msg, "mid")
				info := redis.HGetAll(rid + "/pub/media/" + mid)
				infoStr := util.MarshalStrMap(info)
				resp := util.Map("response", proto.IslbGetMediaInfo, "mid", mid, "info", infoStr)
				log.Infof("amqp.RpcCall from=%s resp=%v corrID=%s", from, resp, corrID)
				amqp.RpcCall(from, resp, corrID)
			case proto.IslbRelay:
				rid := util.Val(msg, "rid")
				mid := util.Val(msg, "mid")
				pid := getPIDFromMID(mid)
				info := redis.HGetAll(rid + "/pub/node/" + pid)
				for ip := range info {
					method := util.Map("method", proto.IslbRelay, "pid", pid, "sid", from, "mid", mid)
					log.Infof("amqp.RpcCall ip=%s, method=%v", ip, method)
					amqp.RpcCall(ip, method, "")
				}
			case proto.IslbUnrelay:
				rid := util.Val(msg, "rid")
				mid := util.Val(msg, "mid")
				info := redis.HGetAll(rid + "/pub/node/" + mid)
				for ip := range info {
					method := util.Map("method", proto.IslbUnrelay, "mid", mid, "sid", from)
					log.Infof("amqp.RpcCall ip=%s, method=%v", ip, method)
					amqp.RpcCall(ip, method, "")
				}
				// time.Sleep(time.Millisecond * 10)
				resp := util.Map("response", proto.IslbUnrelay, "mid", mid, "sid", from)
				log.Infof("amqp.RpcCall from=%s resp=%v corrID=%s", from, resp, corrID)
				amqp.RpcCall(from, resp, corrID)
			}
		}
	}()

}

func handleBroadCastMsgs() {
	broadCastMsgs, err := amqp.ConsumeBroadcast()
	if err != nil {
		log.Errorf(err.Error())
	}

	go func() {
		for m := range broadCastMsgs {
			msg := util.Unmarshal(string(m.Body))
			log.Infof("broadcast msg=%v", msg)
			from := util.Val(msg, "_from")
			id := util.Val(msg, "id")
			method := util.Val(msg, "method")
			if from == proto.IslbID || id == "" || method == "" {
				continue
			}
			switch method {
			}
		}
	}()
}
