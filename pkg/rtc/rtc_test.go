package rtc

import (
	"fmt"
	"testing"
	"time"

	"github.com/pion/ion/pkg/rtc/rtpengine"
	"github.com/pion/ion/pkg/rtc/transport"
	"github.com/pion/rtp"
)

func TestRTPEngineAcceptAndRead(t *testing.T) {
	connCh, err := rtpengine.Serve(6789)
	if err != nil {
		t.Fatal("TestRTPEngineAcceptAndRead ", err)
	}

	go func() {
		for {
			select {
			case rtpTransport := <-connCh:
				fmt.Println("accept new conn from connCh", rtpTransport.RemoteAddr().String())
				go func() {
					for {
						// must read otherwise can't get new conn
						pkt, _ := rtpTransport.ReadRTP()
						fmt.Println("read rtp", pkt)
					}
				}()
			}
		}
	}()

	for i := 0; i < 1; i++ {
		rawPkt := []byte{
			0x90, 0xe0, 0x69, 0x8f, 0xd9, 0xc2, 0x93, 0xda, 0x1c, 0x64,
			0x27, 0x82, 0x00, 0x01, 0x00, 0x01, 0xFF, 0xFF, 0xFF, 0xFF, 0x98, 0x36, 0xbe, 0x88, 0x9e,
		}

		rtp := &rtp.Packet{}
		rtpTransport := transport.NewOutRTPTransport("awsome", "0.0.0.0:6789")
		if err := rtp.Unmarshal(rawPkt); err == nil {
			rtpTransport.WriteRTP(rtp)
		} else {
			fmt.Println("rtpTransport.WriteRTP ", err)
		}
		time.Sleep(time.Second)
	}
}

func TestRTPEngineAcceptKCPAndRead(t *testing.T) {
	connCh, err := rtpengine.ServeWithKCP(1234, "key", "salt")
	if err != nil {
		t.Fatal("TestRTPEngineAcceptKCPAndRead ", err)
	}
	go func() {
		for {
			select {
			case rtpTransport := <-connCh:
				fmt.Println("accept new conn over kcp from connCh", rtpTransport.RemoteAddr().String())
				go func() {
					for {
						// must read otherwise can't get new conn
						pkt, _ := rtpTransport.ReadRTP()
						fmt.Println("read rtp over kcp", pkt)
					}
				}()
			}
		}
	}()

	for i := 0; i < 1; i++ {
		rawPkt := []byte{
			0x90, 0xe0, 0x69, 0x8f, 0xd9, 0xc2, 0x93, 0xda, 0x1c, 0x64,
			0x27, 0x82, 0x00, 0x01, 0x00, 0x01, 0xFF, 0xFF, 0xFF, 0xFF, 0x98, 0x36, 0xbe, 0x88, 0x9e,
		}

		rtp := &rtp.Packet{}
		rtpTransport := transport.NewOutRTPTransportWithKCP("awsome", "0.0.0.0:1234", "key", "salt")
		if err := rtp.Unmarshal(rawPkt); err == nil {
			rtpTransport.WriteRTP(rtp)
		} else {
			fmt.Println("rtpTransport.WriteRTP ", err)
		}
		time.Sleep(time.Second)
	}
}
