package echonetlite

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/exp/slog"
)

var (
	ErrAlreadyStarted  = errors.New("a controller is already started")
	ErrNotQueryMessage = errors.New("not a query message")
)

const (
	broadcastAddress = "224.0.23.0:3610"
)

type Controller struct {
	connectionCancel func()
	receivers        receiverCollection
	currentTid       uint32
	Logger           *slog.Logger
}

func NewController() Controller {
	return Controller{
		receivers: newReceiverCollection(),
		Logger:    slog.Default(),
	}
}

func (c *Controller) CreateFrame() Frame {
	nextTid := atomic.AddUint32(&c.currentTid, 1)
	return Frame{
		Ehd1: 0x10,
		Ehd2: 0x81,
		Tid:  uint16(nextTid & 0xffff),
	}
}

func (c *Controller) Close() error {
	return c.Stop()
}

func (c *Controller) Start() error {
	if c.connectionCancel != nil {
		return ErrAlreadyStarted
	}

	udpAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:3610")
	if err != nil {
		slog.Debug("%+v", err)
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		slog.Debug("%+v", err)
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.connectionCancel = cancel

	c.receivers = newReceiverCollection()

	go c.udpListener(ctx, conn)
	return nil
}

func (c *Controller) udpListener(ctx context.Context, conn *net.UDPConn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn.SetDeadline(time.Now().Add(1 * time.Second))
			n, addr, err := conn.ReadFromUDP(buf)
			if err, ok := err.(net.Error); ok && err.Timeout() {
				continue
			} else if err != nil {
				c.Logger.Debug("[udpListener] error", "err", err)
				continue
			} else if n >= 1024 {
				c.Logger.Debug("[udpListener] large data is arrived, ignore it")
				continue
			}

			data := buf[:n]
			frame, err := DeserializeFrame(data)
			if err != nil {
				c.Logger.Debug("[udpListener] error", "err", err)
				continue
			}
			c.receivers.AcceptAll(*addr, frame)
		}
	}
}

func (c *Controller) Stop() error {
	if c.connectionCancel == nil {
		return nil
	}

	c.connectionCancel()
	c.connectionCancel = nil

	return nil
}

func (c *Controller) QueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		controller: c,
	}
}

func (*Controller) Execute(f Frame) error {
	panic("not implemented")
}

type QueryBuilder struct {
	controller *Controller
	Timeout    time.Duration
}

func (q *QueryBuilder) SetTimeout(duration time.Duration) *QueryBuilder {
	q.Timeout = duration
	return q
}

func (q QueryBuilder) Query(f Frame) ([]QueryResponse, error) {
	if f.Ehd1 != 0x10 || f.Ehd2 != 0x81 || f.Edata.Esv != 0x62 {
		return nil, ErrNotQueryMessage
	}

	udpAddr, err := net.ResolveUDPAddr("udp", broadcastAddress)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	receiver := newReceiver(f.Tid)
	q.controller.receivers.Add(receiver)
	defer q.controller.receivers.Remove(receiver)

	frameBytes, err := f.Serialize()
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(frameBytes)
	if err != nil {
		return nil, err
	}

	time.Sleep(q.Timeout)

	responses := []QueryResponse{}
	for _, a := range receiver.data {
		responses = append(responses, QueryResponse{
			Addr:  a.addr,
			Frame: a.data,
		})
	}

	return responses, nil
}

type QueryResponse struct {
	Addr  net.UDPAddr
	Frame Frame
}

type responseReceiver interface {
	Accept(addr net.UDPAddr, data Frame) bool
}

type receiverData struct {
	addr net.UDPAddr
	data Frame
}

type receiver struct {
	tid  uint16
	data []receiverData
}

func newReceiver(tid uint16) *receiver {
	return &receiver{tid: tid, data: []receiverData{}}
}

func (r *receiver) Accept(addr net.UDPAddr, frame Frame) bool {
	if frame.Tid != r.tid {
		return false
	}
	r.data = append(r.data, receiverData{addr, frame})
	return true
}

type receiverCollection struct {
	m         sync.Mutex
	receivers []responseReceiver
}

func newReceiverCollection() receiverCollection {
	return receiverCollection{
		receivers: []responseReceiver{},
	}
}

func (c *receiverCollection) Add(r responseReceiver) {
	c.m.Lock()
	defer c.m.Unlock()

	c.receivers = append(c.receivers, r)
}

func (c *receiverCollection) Remove(target responseReceiver) {
	c.m.Lock()
	defer c.m.Unlock()

	newSlice := []responseReceiver{}
	for _, r := range c.receivers {
		if r != target {
			newSlice = append(newSlice, r)
		}
	}
	c.receivers = newSlice
}

func (c *receiverCollection) AcceptAll(addr net.UDPAddr, f Frame) {
	c.m.Lock()
	defer c.m.Unlock()

	for _, r := range c.receivers {
		r.Accept(addr, f)
	}
}
