package biz

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/cloudwebrtc/nats-discovery/pkg/discovery"
	nrpc "github.com/cloudwebrtc/nats-grpc/pkg/rpc"
	"github.com/nats-io/nats.go"
	log "github.com/pion/ion-log"
	biz "github.com/pion/ion/pkg/grpc/biz"
	islb "github.com/pion/ion/pkg/grpc/islb"
	sfu "github.com/pion/ion/pkg/grpc/sfu"
	"github.com/pion/ion/pkg/proto"
	"github.com/pion/ion/pkg/util"
)

// BizServer represents an BizServer instance
type BizServer struct {
	biz.UnimplementedBizServer
	sfu.UnimplementedSFUServer
	nc       *nats.Conn
	elements []string
	roomLock sync.RWMutex
	rooms    map[string]*Room
	closed   chan bool
	islbcli  islb.ISLBClient
	nodeLock sync.RWMutex
	nodes    map[string]*discovery.Node
}

// newBizServer creates a new avp server instance
func newBizServer(dc string, nid string, elements []string, nc *nats.Conn) *BizServer {
	return &BizServer{
		nc:       nc,
		elements: elements,
		rooms:    make(map[string]*Room),
		closed:   make(chan bool),
		nodes:    make(map[string]*discovery.Node),
	}
}

func (s *BizServer) start() error {
	go s.stat()
	return nil
}

func (s *BizServer) close() {
	close(s.closed)
}

// getRoom get a room by id
func (s *BizServer) getRoom(id string) *Room {
	s.roomLock.Lock()
	defer s.roomLock.Unlock()
	r := s.rooms[id]
	return r
}

/*func (s *BizServer) BridgeTest(sid, uid string) {

	if s.islbcli == nil {
		for nid, _ := range s.nodes {
			ncli := nrpc.NewClient(s.nc, nid)
			s.islbcli = islb.NewISLBClient(ncli)
		}
		resp, err := s.islbcli.FindNode(context.TODO(), &islb.FindNodeRequest{
			Service: proto.ServiceSFU,
			Sid:     sid,
		})
		if err == nil && len(resp.Nodes) > 0 {

			peer := &Peer{
				uid:  uid,
				sid:  sid,
				info: []byte(""),
			}

			r := s.getRoom(sid)
			if r == nil {
				r = newRoom(sid, resp.GetNodes()[0].Nid)
				s.roomLock.RLock()
				s.rooms[sid] = r
				s.roomLock.RUnlock()
			}
			r.addPeer(peer)
		}
	}
}*/

//Signal bridge SFU signaling between client and sfu node.
func (s *BizServer) Signal(sstream sfu.SFU_SignalServer) error {
	var peer *Peer = nil
	var cstream sfu.SFU_SignalClient = nil
	reqCh := make(chan *sfu.SignalRequest)
	replyCh := make(chan *sfu.SignalReply)
	errCh := make(chan error)

	defer func() {
		if cstream != nil {
			cstream.CloseSend()
		}
		close(reqCh)
		close(replyCh)
		close(errCh)
	}()

	go func() {
		for {
			req, err := sstream.Recv()
			if err != nil {
				log.Errorf("Singal server stream.Recv() err: %v", err)
				errCh <- err
				return
			}
			reqCh <- req
		}
	}()

	for {
		select {
		case err, _ := <-errCh:
			return err
		case req, ok := <-reqCh:

			if !ok {
				return io.EOF
			}

			if cstream != nil {
				cstream.Send(req)
				break
			}

			switch payload := req.Payload.(type) {
			case *sfu.SignalRequest_Join:
				//s.BridgeTest(payload.Join.Sid, payload.Join.Uid)
				r := s.getRoom(payload.Join.Sid)
				if r != nil {
					peer = r.getPeer(payload.Join.Uid)
					if peer != nil {
						//Use nats-grpc or grpc
						cli := sfu.NewSFUClient(nrpc.NewClient(s.nc, r.NID()))
						var err error
						cstream, err = cli.Signal(context.Background())
						if err != nil {
							log.Errorf("Singal cli.Signal() err: %v", err)
							return err
						}

						go func() {
							for {
								reply, err := cstream.Recv()
								if err != nil {
									log.Errorf("Singal client stream.Recv() err: %v", err)
									errCh <- err
									return
								}
								replyCh <- reply
							}
						}()

						cstream.Send(req)
						break
					} else {
						return fmt.Errorf("peer [%v] not found", payload.Join.Uid)
					}
				} else {
					return fmt.Errorf("session [%v] not found", payload.Join.Sid)
				}
			}
		case reply, ok := <-replyCh:
			if ok {
				sstream.Send(reply)
				break
			}
			return io.EOF
		}
	}
}

//Join process biz request.
func (s *BizServer) Join(stream biz.Biz_JoinServer) error {
	var r *Room = nil
	var peer *Peer = nil

	defer func() {
		if peer != nil && r != nil {
			peer.Close()
			r.delPeer(peer.UID())
		}
	}()

	for {
		req, err := stream.Recv()
		if err != nil {
			log.Errorf("Biz stream.Recv() err: %v", err)
			return err
		}

		log.Infof("Biz request => %v", req.String())

		switch payload := req.Payload.(type) {
		case *biz.JoinRequest_Join:
			sid := payload.Join.Sid
			uid := payload.Join.Uid
			success := false
			reason := "unkown error."

			if s.islbcli == nil {
				s.nodeLock.Lock()
				defer s.nodeLock.Unlock()
				for _, node := range s.nodes {
					if node.Service == proto.ServiceISLB {
						ncli := nrpc.NewClient(s.nc, node.NID)
						s.islbcli = islb.NewISLBClient(ncli)
						break
					}
				}
			}

			if s.islbcli != nil {
				r = s.getRoom(sid)
				if r == nil {
					resp, err := s.islbcli.FindNode(context.TODO(), &islb.FindNodeRequest{
						Service: proto.ServiceSFU,
						Sid:     sid,
					})
					if err == nil && len(resp.Nodes) > 0 {
						r = newRoom(sid, resp.GetNodes()[0].Nid)
						s.roomLock.RLock()
						s.rooms[sid] = r
						s.roomLock.RUnlock()
					} else {
						reason = fmt.Sprintf("islbcli.FindNode(serivce = sfu, sid = %v) err %v", sid, err)
					}
				}
				if r != nil {
					peer = &Peer{
						uid:  uid,
						sid:  sid,
						info: payload.Join.Info,
					}
					r.addPeer(peer)
					success = true
					reason = "join success."
				} else {
					reason = fmt.Sprintf("room sid = %v not found", sid)
				}
			} else {
				reason = fmt.Sprintf("islb node not found")
			}

			stream.Send(&biz.JoinReply{
				Payload: &biz.JoinReply_Result{
					Result: &biz.JoinResult{
						Success: success,
						Reason:  reason,
					},
				},
			})

			break
		case *biz.JoinRequest_Leave, *biz.JoinRequest_Msg:
			if peer != nil {
				reply, err := peer.handleRequest(req)
				if err != nil {
					log.Errorf("Biz request err %v", err)
					return err
				}
				stream.Send(reply)
			}
		default:
			break
		}
		log.Infof("req => %v", req.String())
	}
}

// watchNodes watch islb nodes up/down
func (s *BizServer) watchNodes(state discovery.NodeState, node *discovery.Node) {
	s.nodeLock.Lock()
	defer s.nodeLock.Unlock()
	id := node.NID
	service := node.Service
	if state == discovery.NodeUp {
		log.Infof("Service up: "+service+" node id => [%v], rpc => %v", id, node.RPC.Protocol)
		if _, found := s.nodes[id]; !found {
			s.nodes[id] = node
		}
	} else if state == discovery.NodeDown {
		log.Infof("Service down: "+service+" node id => [%v]", id)
		delete(s.nodes, id)
	}
}

// stat peers
func (s *BizServer) stat() {
	t := time.NewTicker(util.DefaultStatCycle)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			break
		case <-s.closed:
			log.Infof("stop stat")
			return
		}

		var info string
		s.roomLock.RLock()
		for sid, room := range s.rooms {
			info += fmt.Sprintf("room: %s\npeers: %d\n", sid, room.count())
		}
		s.roomLock.RUnlock()
		if len(info) > 0 {
			log.Infof("\n----------------signal-----------------\n" + info)
		}
	}
}

/*
func (s *Server) handle(msg interface{}) (interface{}, error) {
	log.Infof("handle incoming message: %T, %+v", msg, msg)

	switch v := msg.(type) {
	case *proto.SfuOfferMsg:
		s.handleSFUMessage(v.SID, v.UID, msg)
	case *proto.SfuTrickleMsg:
		s.handleSFUMessage(v.SID, v.UID, msg)
	case *proto.SfuICEConnectionStateMsg:
		s.handleSFUMessage(v.SID, v.UID, msg)
	default:
		log.Warnf("unkonw message: %v", msg)
	}

	return nil, nil
}

func (s *Server) handleSFUMessage(sid string, uid proto.UID, msg interface{}) {
	if r := s.getRoom(sid); r != nil {
		if p := r.getPeer(uid); p != nil {
			p.handleSFUMessage(msg)
		} else {
			log.Warnf("peer not exits, sid=%s, uid=%s", sid, uid)
		}
	} else {
		log.Warnf("room not exits, sid=%s, uid=%s", sid, uid)
	}
}

func (s *Server) broadcast(msg interface{}) (interface{}, error) {
	log.Infof("handle islb message: %T, %+v", msg, msg)

	var sid string
	var uid proto.UID

	switch v := msg.(type) {
	case *proto.FromIslbStreamAddMsg:
		sid, uid = v.SID, v.UID
	case *proto.ToClientPeerJoinMsg:
		sid, uid = v.SID, v.UID
	case *proto.IslbPeerLeaveMsg:
		sid, uid = v.SID, v.UID
	case *proto.IslbBroadcastMsg:
		sid, uid = v.SID, v.UID
	default:
		log.Warnf("unkonw message: %v", msg)
	}

	if r := s.getRoom(sid); r != nil {
		r.send(msg, uid)
	} else {
		log.Warnf("room not exits, sid=%s, uid=%s", sid, uid)
	}

	return nil, nil
}



// // getRoomsByPeer a peer in many room
// func (s *server) getRoomsByPeer(uid proto.UID) []*room {
// 	var result []*room
// 	s.roomLock.RLock()
// 	defer s.roomLock.RUnlock()
// 	for _, r := range s.rooms {
// 		if p := r.getPeer(uid); p != nil {
// 			result = append(result, r)
// 		}
// 	}
// 	return result
// }

// delPeer delete a peer in the room
func (s *Server) delPeer(sid string, uid proto.UID) {
	log.Infof("delPeer sid=%s uid=%s", sid, uid)
	room := s.getRoom(sid)
	if room == nil {
		log.Warnf("room not exits, sid=%s, uid=%s", sid, uid)
		return
	}
	if room.delPeer(uid) == 0 {
		s.roomLock.RLock()
		delete(s.rooms, sid)
		s.roomLock.RUnlock()
	}
}

// addPeer add a peer to room
func (s *Server) addPeer(sid string, peer *Peer) {
	log.Infof("addPeer sid=%s uid=%s", sid, peer.uid)
	room := s.getRoom(sid)
	if room == nil {
		room = newRoom(sid)
		s.roomLock.Lock()
		s.rooms[sid] = room
		s.roomLock.Unlock()
	}
	room.addPeer(peer)
}

// // getPeer get a peer in the room
// func (s *server) getPeer(sid string, uid proto.UID) *peer {
// 	log.Infof("getPeer sid=%s uid=%s", sid, uid)
// 	r := s.getRoom(sid)
// 	if r == nil {
// 		log.Infof("room not exits, sid=%s uid=%s", sid, uid)
// 		return nil
// 	}
// 	return r.getPeer(uid)
// }

*/
