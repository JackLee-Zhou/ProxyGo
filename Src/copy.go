package Src

import (
	"ProxyGo/ProxyLog"
	"io"
)

// Copy 用于远端和本地端的数据同步
func Copy(src io.ReadWriteCloser, dst io.ReadWriteCloser) (num int, err error) {
	temp := []byte{}
	for {
		// 从远端读
		sn, err := src.Read(temp)
		if err != nil {
			ProxyLog.Err.Println("SrcReadErr ", err)
			return 0, err
		}
		// 表明有数据
		if sn > 0 {
			// 向远端写数据
			wn, err := dst.Write(temp[:sn])
			if err != nil {
				ProxyLog.Err.Println("DstWriteErr ", err)
				return num, err
			}
			// 发送成功才累加
			if wn > 0 {
				// 累加发送的数据个数
				num += wn
			}
			// 发送的个数不同
			if wn != sn {
				ProxyLog.Err.Println("DstSendNotComplete ", io.ErrShortWrite)
				break
			}
		} else {
			ProxyLog.Info.Println("SrcReadEmpty")
			break
		}
	}
	return num, err
}
