package Network

import (
	"io"
	"encoding/binary"
	"fmt"
	"errors"
	"time"
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
	b.grow() // 有用的字节前移
	n, err := b.reader.Read(b.buf[b.end:])
	a := "读到了n: " + strconv.Itoa(n) + ", b.start: " + strconv.Itoa(b.start) + ", " + "b.end: " + strconv.Itoa(b.end)
	fmt.Println(colors.Green(a))
	if (err != nil) {
		fmt.Println(err.Error())
		return n, err
	}
	b.end += n
	return n, nil
}

// 读一个处理一个， todo 异步读
func (buffer *Buffer) ThomasRead() ([]byte, error) {
	for {
		_, err := buffer.readFromReader()
		fmt.Println("b.start:", buffer.start, ", ", "b.end", buffer.end)
		if err != nil {
			fmt.Println(err)
			return []byte(""), err
		}
		for {
			headBuf, err := buffer.seek(HEAD_SIZE)
			if err != nil {
				break
			}
			contentSize := int(binary.BigEndian.Uint16(headBuf))
			fmt.Println("contentSize:", contentSize);
			if (buffer.len() >= contentSize-HEAD_SIZE) {
				fmt.Println("一次读够了返回去")
				contentBuf := buffer.read(HEAD_SIZE, contentSize)
				fmt.Println("b.start:", buffer.start, ", ", "b.end", buffer.end)
				//buffer.start -= contentSize
				return contentBuf, nil
			}
			fmt.Println("一次读的不够")
			time.Sleep(1 * time.Second)
			break
		}
	}
}
