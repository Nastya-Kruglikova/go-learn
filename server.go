package main

import (
    "bufio"
    "log"
    "net"
    "sync"
    "net/textproto"

)

type FuncHandler func(net.Conn)

func (h FuncHandler) Serve(conn net.Conn) {
    h(conn)
}

type Handler interface {
    Serve(conn net.Conn)
}

type Server struct {
    connLock    sync.RWMutex
    users       []*User
    listener    net.Listener
    rooms       []*Room
}

type User struct {
    connection net.Conn
    name []byte
}

type Room struct {
    connLock    sync.RWMutex
    users []*User
    name string
    messages chan string
}

func (room *Room) SendAll(data string) {
    for _, user := range room.users {
        room.connLock.RLock()
        user.connection.Write([]byte(data))
        room.connLock.RUnlock()
    }
}

func (room *Room) WaitForMessages() {
    go func() {
        for {
            if message, more := <- room.messages; more {
                room.SendAll(message)
            } else {
                return
            }
        }
    }()
}

func NewRoom(name string) *Room {
    room := &Room{
        name: name,
        messages: make(chan string),
    }

    room.WaitForMessages()

    return room
}

func (s *Server) echo(conn net.Conn) {
    r := bufio.NewReader(conn)
    tp := textproto.NewReader(r)

    for {
        line, err := tp.ReadLine()
        if err != nil {
            log.Println("Error reading: ", err)
            return
        }
        for _, r := range s.rooms {
            r.messages <- line
        }
    }
}

func (s *Server) ListenAndServe(address string, handler Handler) error {
    listener, err := net.Listen("tcp", address)
    if err != nil {
        return err
    }
    s.listener = listener

    for {
        conn, err := listener.Accept()
        if err != nil {
            return err
        }
        conn.Write([]byte("Enter your name:"))
        r := bufio.NewReader(conn)
        line, err := r.ReadBytes('\n')
        if err != nil {
            log.Fatal(err)
        }
        s.connLock.Lock()
        s.users = append(s.users, &User{conn, line})
        s.connLock.Unlock()

        for _, r := range s.rooms {
            conn.Write(append([]byte("you have joined room "), []byte(r.name)...))
            r.connLock.Lock()
            r.users = append(r.users, &User{conn, line})
            r.connLock.Unlock()
        }

        go func() {
            handler.Serve(conn)
            conn.Close()
        }()
    }
}

func main() {
    server := &Server{}
    room := NewRoom("First Room")
    server.rooms = append(server.rooms, room)
    server.ListenAndServe(":8080", FuncHandler(server.echo))
}