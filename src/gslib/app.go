package gslib

import (
	"encoding/binary"
	"fmt"
	"gslib/gen_server"
	"gslib/store"
	. "gslib/utils"
	"io"
	"log"
	"net"
	"runtime"
	"time"
)

func Run() {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("caught panic in main()", x)
		}
	}()

	// runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.GOMAXPROCS(1)

	// go SysRoutine()

	// time.Sleep(1 * time.Second)

	fmt.Println("Server Started!")

	// Init DB Connections
	store.InitDB()

	// player := new(Player)
	// gen_server.Start("test", player, "test")
	// gen_server.Call("test", "HandleWrap", func() interface{} {
	// 	fmt.Println("player: ", player.PlayerId)
	// 	player.Store.LoadData("equips", "54BC69792B897814D763403B")
	// 	// player.Store.Get([]string{"equips"})
	// 	fmt.Println("equips count: ", player.Store.Count([]string{"equips"}))
	// 	equipModel := player.Store.Get([]string{"equips"}, "54BC69792B897814D7634040")
	// 	fmt.Println("equipModel: ", equipModel)
	// 	return "ok"
	// })

	start_tcp_server()
}

func start_tcp_server() {
	l, err := net.Listen("tcp", ":4100")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer func() {
		if x := recover(); x != nil {
			ERR("caught panic in handleClient", x)
		}
	}()

	server_name := conn.RemoteAddr().String()
	gen_server.Start(server_name, new(Player), server_name)

	header := make([]byte, 2)

	for {
		// header
		conn.SetReadDeadline(time.Now().Add(TCP_TIMEOUT * time.Second))
		n, err := io.ReadFull(conn, header)
		if err != nil {
			NOTICE("Connection Closed: ", err)
			break
		}

		// data
		size := binary.BigEndian.Uint16(header)
		data := make([]byte, size)
		n, err = io.ReadFull(conn, data)
		if err != nil {
			WARN("error receiving msg, bytes:", n, "reason:", err)
			break
		}

		gen_server.Cast(server_name, "HandleRequest", data, conn)
	}

	gen_server.Cast(server_name, "removeConn")

}