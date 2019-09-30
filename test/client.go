package test

import (
	"bytes"
	"fmt"
	"net"
	"reflect"
	"sync"

	"github.com/gogo/protobuf/proto"
)

const (
	constBufferLimit = 8096
)

//NewClientTest 创建测试客户端
func NewClientTest() *Client {
	r := &Client{b: bytes.NewBuffer([]byte{}), out: make(chan []byte, 32), stoped: make(chan bool, 1),
		e: make(map[reflect.Type]NetEventMethod)}
	r.b.Grow(constBufferLimit)
	return r
}

//AnalysisMethod 数据包解析方法
type AnalysisMethod func(b *bytes.Buffer) (string, uint64, []byte, error)

//NetEventMethod 网络事件方法
type NetEventMethod func(message interface{})

//NetEventHandle 网络事件
type NetEventHandle struct {
	Handle uint64
	Name   string
	Data   []byte
}

//Client xxx
type Client struct {
	Analysis AnalysisMethod
	s        net.Conn
	b        *bytes.Buffer
	out      chan []byte
	stoped   chan bool
	w        sync.WaitGroup
	e        map[reflect.Type]NetEventMethod
}

//Connect 连接目标
func (c *Client) Connect(host string) error {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}

	c.s = conn
	c.w.Add(2)
	go c.recv()
	go c.send()
	return nil
}

func (c *Client) Write(d []byte) {
	select {
	case stop := <-c.stoped:
		if stop {
			return
		}
	case c.out <- d:
	}
}

//RegisterMethod 注册事件方法
func (c *Client) RegisterMethod(k interface{}, f NetEventMethod) {
	c.e[reflect.TypeOf(k)] = f
}

//Wait 挂起
func (c *Client) Wait() {
	c.w.Wait()
	close(c.stoped)
	close(c.out)
}

//Shutdown 关闭客户端
func (c *Client) Shutdown() {
	if c.s != nil {
		c.s.Close()
	}
}

func (c *Client) recv() {
	defer c.w.Done()
	tmpbuf := make([]byte, 4096)
	for {
		n, err := c.s.Read(tmpbuf)
		if err != nil {
			fmt.Printf("连接错误 Read:bytes-%d error-%+v\n", n, err)
			goto end
		}
		c.write(tmpbuf[:n])
	}
end:
	c.stoped <- true
	c.s.Close()
}

func (c *Client) send() {
	defer c.w.Done()
	for {
		select {
		case stop := <-c.stoped:
			if stop {
				goto end
			}

		case o := <-c.out:
			n, err := c.s.Write(o)
			if err != nil {
				goto end
			}

			if n < len(o) {
				fmt.Printf("写入数据未完成:%d,%d\n", len(o), n)
			}
		}
	}
end:
}

func (c *Client) write(d []byte) error {
	if d == nil || len(d) == 0 {
		return nil
	}

	var (
		space  int
		writed int
		wby    int
		pos    int

		h    uint64
		v    []byte
		name string
		err  error
	)

	for {
		space = c.b.Cap() - c.b.Len()
		wby = len(d) - writed
		if space > 0 && wby > 0 {
			if space > wby {
				space = wby
			}

			_, err = c.b.Write(d[pos : pos+space])
			if err != nil {
				return fmt.Errorf("error close client recv: write buffer error %+v socket %+v", err, c.s)
			}

			pos += space
			writed += space
		}

		for {
			// Decomposition of Packets
			name, h, v, err = c.Analysis(c.b)
			if err != nil {
				return fmt.Errorf("error close client %+v", err)
			}

			if v != nil {
				msgType := proto.MessageType(name)
				fmt.Println(msgType)
				if msgType != nil {
					f, success := c.e[msgType]
					if success {
						f(&NetEventHandle{Handle: h, Name: name, Data: v})
					}
				}
				continue
			}

			if writed >= len(d) {
				return nil
			}
			break
		}
	}
}
