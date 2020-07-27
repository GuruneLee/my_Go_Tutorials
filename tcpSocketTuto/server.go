package main

import (
   "io"
   "log"
   "net"
)

func main() {
  l, err := net.Listen("tcp", ":8000")
  if err != nil {
    log.Println(err)
  }
  defer l.Close()

  for {
    conn, err := l.Accept()
    if err != nil {
      log.Println(err)
      continue
    }
    defer conn.Close()
    go ConnHandler(conn)
  }
}

func ConnHandler(conn net.Conn) {
  recvBuf := make([]byte, 4096)
  for {
    n, err := conn.Read(recvBuf)
    if err != nil {
      if err == io.EOF {
        log.Println(err)
        return
      }
      log.Println(err)
      return
    }

    if n>0 {
      data := recvBuf[:n]
      log.Println(string(data))
      _, err = conn.Write(data[:n])
      if err != nil {
        log.Println(err)
        return
      }
    }
  }
}
