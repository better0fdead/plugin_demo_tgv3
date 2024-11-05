package plugin

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/better0fdead/plugin_demo_tgv3/context"
	"github.com/better0fdead/plugin_demo_tgv3/parser"
	"github.com/better0fdead/plugin_demo_tgv3/render"
)

type SendCtx struct {
	Pr    []byte            `json:"Pr,omitempty"`
	Flags map[string]string `json:"flags,omitempty"`
}

type Description struct {
	Desc    string `json:"Desc,omitempty"`
	Version string `json:"Version,omitempty"`
}

func DesirializeData(jsonData []byte) (context.PluginCtx, parser.PackageInfo, error) {
	var req SendCtx
	log := logrus.WithTime(time.Now())
	err := json.Unmarshal(jsonData, &req)
	if err != nil {
		log.Infof("error deserializing sent Data: %s", err.Error())
		return context.PluginCtx{}, parser.PackageInfo{}, err
	}

	pluginCtx := context.SetPluginCtx(req.Flags)

	var parsedPackage parser.PackageInfo
	if len(req.Pr) > 0 {
		err = json.Unmarshal(req.Pr, &parsedPackage)
		if err != nil {
			log.Infof("error deserializing parser data: %s", err.Error())
			return context.PluginCtx{}, parser.PackageInfo{}, err
		}
		for j, service := range parsedPackage.Services {
			for i, method := range service.Methods {
				var parametrs []parser.FieldPkgInfo
				parametrs = append(parametrs, parser.FieldPkgInfo{Name: "ctx", Kind: "context.Context"})
				parsedPackage.Types["context.Context"] = parser.TypeInfo{Name: "context.Context", IsScalar: true, Pkg: "context"}
				parametrs = append(parametrs, method.Parameters...)
				returns := append(method.Returns, parser.FieldPkgInfo{Name: "err", Kind: "error", IsScalar: true})
				parsedPackage.Services[j].Methods[i].Parameters = parametrs
				parsedPackage.Services[j].Methods[i].Returns = returns
			}
		}
	}

	return pluginCtx, parsedPackage, err
}

func Start(about []byte, source, description, version, help string) error {
	// Create a Unix domain socket and listen for incoming connections.
	socket, err := net.Listen("unix", "./plugin.sock")
	if err != nil {
		log.Fatal(err)
	}

	// Cleanup the sockfile.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove("./plugin.sock")
		os.Exit(1)
	}()

	for {
		// Accept an incoming connection.
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Handle the connection in a separate goroutine.
		go func(conn net.Conn) {
			defer conn.Close()

			buf := make([]byte, 500)

			n, err := conn.Read(buf)
			if err != nil {
				log.Fatal(err)
			}

			msgLen, _ := strconv.Atoi(string(buf[:n]))
			msgBuf := make([]byte, msgLen)
			msgLenStr := strconv.Itoa(msgLen)
			conn.Write([]byte(msgLenStr))

			n, err = conn.Read(msgBuf)
			if err != nil {
				log.Fatal(err)
			}

			pluginCtx, parsedPackage, err := DesirializeData(msgBuf[:n])
			if err != nil {
				conn.Write([]byte(err.Error()))
				return
			}
			if pluginCtx.Help {
				conn.Write([]byte(help))
				return
			}
			if pluginCtx.Doc {
				conn.Write([]byte(about))
				return
			}
			if pluginCtx.Source {
				conn.Write([]byte(source))
				return
			}

			if pluginCtx.Flags {
				flags := context.GetPluginFlags()
				conn.Write([]byte(strings.Join(flags, " ")))
			}

			if pluginCtx.Description {
				description := Description{Desc: description, Version: version}
				mDescription, err := json.Marshal(description)
				if err != nil {
					conn.Write([]byte(err.Error()))
					return
				}
				conn.Write(mDescription)
				return
			}

			err = render.Render(pluginCtx, parsedPackage)
			if err != nil {
				conn.Write([]byte(err.Error()))
			} else {
				conn.Write([]byte("done"))
			}

		}(conn)
	}
	return err
}
