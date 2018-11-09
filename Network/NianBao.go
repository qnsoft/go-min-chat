package Network

import (
	"io"
	"encoding/binary"
	"fmt"
	"errors"
	"github.com/beego/bee/logger/colors"
	"strconv"
)

const HEAD_SIZE = 2

type Buffer struct {
	reader io.Reader
	buf    []byte
	start  int
	end    int
}

func NewBuffer(reader io.Reader, len int) *Buffer {
	buf := make([]byte, len)
	return &Buffer{reader, buf, 0, 0}
}

// grow 将有用的字节前移
func (b *Buffer) grow() {
	if b.start == 0 {
		return
	}
	copy(b.buf, b.buf[b.start:b.end])
	b.end -= b.start
	b.start = 0
}

func (b *Buffer) len() int {
	return b.end - b.start
}

//返回n个字节，而不产生移位
func (b *Buffer) seek(n int) ([]byte, error) {
	if b.end-b.start >= n {
		buf := b.buf[b.start : b.start+n]
		return buf, nil
	}
	return nil, errors.New("not enough")
}

//舍弃offset个字段，读取n个字段
func (b *Buffer) read(offset, n int) ([]byte) {
	b.start += offset
	buf := b.buf[b.start : b.start+n]
	b.start += n
	return buf
}

//从reader里面读取数据，如果reader阻塞，会发生阻塞
func (b *Buffer) readFromReader() (int, error) {
	n, err := b.reader.Read(b.buf[b.end:])
	a := "读到了n: " + strconv.Itoa(n) + ", b.start: " + strconv.Itoa(b.start) + ", " + "b.end: " + strconv.Itoa(b.end)
	fmt.Println(colors.Green(a))
	if (err != nil) {
		fmt.Println(err.Error())
		panic(err.Error())
		return n, err
	}
	b.end += n
	return n, nil
}

// 读一个处理一个， todo 异步读
func (buffer *Buffer) ThomasRead(msg chan string) (error) {
	for {
		buffer.grow()                     // 移动数据
		_, err := buffer.readFromReader() // 读数据拼接到定额缓存后面
		fmt.Println("b.start:", buffer.start, ", ", "b.end", buffer.end)
		if err != nil {
			fmt.Println(err)
			return err
		}
		// 检查定额缓存里面的数据有几个消息(可能不到1个，可能连一个消息头都不够，可能有几个完整消息+一个消息的部分)
		isBreak := buffer.checkMsg(msg)
		if (isBreak) {
			return nil
		}
	}
}

func (buffer *Buffer) checkMsg(msg chan string) (bool) {
	var isBreak bool
	fmt.Println("111")
	headBuf, err1 := buffer.seek(HEAD_SIZE)
	if err1 != nil { // 一个消息头都不够， 跳出去继续读吧
		return false
	}
	contentSize := int(binary.BigEndian.Uint16(headBuf))
	fmt.Println("contentSize:", contentSize);
	if (buffer.len() >= contentSize-HEAD_SIZE) { // 一个消息体也是够的
		fmt.Println("一次读够了返回去")
		contentBuf := buffer.read(HEAD_SIZE, contentSize) // 把消息读出来，把start往后移
		msg <- string(contentBuf)
		fmt.Println("len(msg): ", len(msg))
		fmt.Println("b.start:", buffer.start, ", ", "b.end", buffer.end)
		// 那么继续，看剩下的还够一个消息不
		isBreak = true
		buffer.checkMsg(msg)
	} else { // 一个消息体不够的， 跳出去继续读吧
		isBreak = false
	}
	return isBreak
}
