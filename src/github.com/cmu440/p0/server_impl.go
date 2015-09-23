// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0

import (
	"net"
	"strconv"
	"bufio"
)

type multiEchoServer struct {
	workers []workerServer
	count   int
	close   chan bool
	ctrl chan bool
}

type workerServer struct {
	master *multiEchoServer
	conn   net.Conn
	inMsg  chan string
	outMsg chan string
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	rst := &multiEchoServer {
		workers: make([]workerServer, 0),
		count: 0,
		close: make(chan bool),
		ctrl: make(chan bool, 1),
	}

	rst.ctrl <- true
	return rst
}

func (mes *multiEchoServer) Start(port int) error {
	portStr := strconv.Itoa(port)

	ln, err := net.Listen("tcp", ":" + portStr)
	if err != nil {
		return err
	}

	go mes.handleAccept(ln)
	return nil
}

func (mes *multiEchoServer) Close() {
	close(mes.close)
	for _, s := range mes.workers {
		s.conn.Close()
	}
}

func (mes *multiEchoServer) Count() int {
	<-mes.ctrl
	count := mes.count
	mes.ctrl<- true
	return count
}

// TODO: add additional methods/functions below!
func (mes *multiEchoServer) handleAccept(ln net.Listener) {
	for {
		select {
		case <-mes.close:
			ln.Close()
			return
		default:
		}
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			worker := workerServer{
				master: mes,
				conn: conn,
				inMsg: make(chan string, 75),
				outMsg: make(chan string),
			}

			go mes.addWorker(&worker)
	}
}


func (mes *multiEchoServer)addWorker(worker *workerServer) {
	<-mes.ctrl
	mes.workers = append(mes.workers, *worker)
	mes.count++
	mes.ctrl<- true

	go worker.handleConn()
}

func (worker *workerServer)handleConn() {
	go worker.handleReq()
	go worker.listenMsg()
	go worker.handleRes()
}

func (worker *workerServer) handleRes() {
	for inMsg := range worker.inMsg{
		worker.conn.Write([]byte(inMsg))
	}
}

func (worker *workerServer) handleReq() {
	b := bufio.NewReader(worker.conn)
	for {
		line, err := b.ReadBytes('\n')
		if err != nil {
			worker.conn.Close()

			<-worker.master.ctrl
			worker.master.count--
			worker.master.ctrl<- true
			break
		}
		outMsg := append(line[:len(line) - 1], line...) // Trim the first "\n"
		worker.outMsg <- string(outMsg)
	}
}


func (worker *workerServer) listenMsg() {
	for outMsg := range worker.outMsg {
		<-worker.master.ctrl
		for _, w := range worker.master.workers {
			if len(w.inMsg) < 75 {
				w.inMsg <- outMsg
			}
		}
		worker.master.ctrl <- true
	}
}