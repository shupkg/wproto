package cmd

import (
	"fmt"
	"github.com/shupkg/wproto/gen"
	"github.com/spf13/pflag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func New() *ProtoGen {
	return &ProtoGen{
		Out:    ".",
		Protoc: "/usr/local/bin/protoc",
		Mod:    "",
	}
}

func WithFlag(flag *pflag.FlagSet) *ProtoGen {
	g := New()
	flag.SortFlags = false
	flag.StringVarP(&g.Out, "out", "o", g.Out, "输出目录")
	flag.StringVar(&g.Mod, "mod", g.Mod, "根包名")
	flag.BoolVar(&g.Pb, "pb", g.Pb, "生成 golang pb 文件")
	flag.BoolVar(&g.Grpc, "grpc", g.Grpc, "生成 golang grpc 服务文件")
	flag.BoolVar(&g.GrpcNoUnimplementedServers, "grpc.no_un_impl", g.GrpcNoUnimplementedServers, "grpc requireUnimplementedServers=false")
	flag.BoolVar(&g.GosTs, "gos.ts", g.GosTs, "生成 typescript gos 客户端文件")
	flag.BoolVar(&g.GosModel, "gos.model", g.GosModel, "生成 gos model 文件")
	flag.BoolVar(&g.GosRpc, "gos.rpc", g.GosRpc, "生成 gos rpc 文件")
	flag.StringVar(&g.Protoc, "protoc", g.Protoc, "protoc 路径")
	flag.StringSliceVar(&g.Files, "files", g.Files, "proto文件(夹)")
	flag.StringSliceVar(&g.ProtocFlags, "flags", g.ProtocFlags, "附加更多的命名到protoc")
	return g
}

type ProtoGen struct {
	Out                        string   //输出文件夹
	Mod                        string   //跟路径包名，只有包名以他开头的才会编码
	Pb                         bool     //生成pb代码
	Grpc                       bool     //生成grpc代码
	GrpcNoUnimplementedServers bool     //生成 requireUnimplementedServers
	GosModel                   bool     //生成模型
	GosTs                      bool     //生成ts
	GosRpc                     bool     //生成xxx
	Files                      []string //proto 文件
	Protoc                     string   //protoc 路径
	ProtocFlags                []string //附加的Protoc命令
}

func (g *ProtoGen) Walk() []string {
	var root, _ = filepath.Abs(".")
	var files []string
	for _, f := range g.Files {
		f, _ = filepath.Abs(f)
		_ = filepath.Walk(f, func(truePath string, info os.FileInfo, err error) error {
			if !strings.HasPrefix(truePath, root) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if info != nil && !info.IsDir() && filepath.Ext(truePath) == ".proto" {
				path, _ := filepath.Rel(root, truePath)
				files = append(files, path)
			}

			return nil
		})
	}
	return files
}

func (g *ProtoGen) genParams() string {
	param := gen.Param{}
	if g.GosTs {
		param.Set("gos_ts", "1")
	}

	if g.GosRpc {
		param.Set("gos_rpc", "1")
	}

	if g.GosModel {
		param.Set("gos_model", "1")
	}

	if g.Mod == "" {
		param.Set("paths", "source_relative")
	} else {
		param.Set("module", g.Mod)
	}

	return param.Encode()
}

func (g *ProtoGen) Run() error {
	var args []string
	bin, _ := os.Executable()
	args = append(args, "-I.")

	if g.Out == "" {
		g.Out = "."
	}

	if g.Pb {
		args = append(args, fmt.Sprintf("--go_out=%s", g.Out))
		if g.Mod == "" {
			args = append(args, "--go_opt=paths=source_relative")
		} else {
			args = append(args, fmt.Sprintf("--go_opt=module=%s", g.Mod))
		}
	}

	if g.Grpc {
		args = append(args, fmt.Sprintf("--go-grpc_out=%s", g.Out))
		var grpcOpt = "--go-grpc_opt="
		if g.Mod == "" {
			grpcOpt += "paths=source_relative"
		} else {
			grpcOpt += fmt.Sprintf("module=%s", g.Mod)
		}
		if g.GrpcNoUnimplementedServers {
			grpcOpt += ",requireUnimplementedServers=false"
		}
		args = append(args, grpcOpt)
	}

	if g.GosModel || g.GosRpc || g.GosTs {
		args = append(args, "--plugin=protoc-gen-go-gos="+bin)
		args = append(args, fmt.Sprintf("--go-gos_out=%s", g.Out))
		args = append(args, fmt.Sprintf("--go-gos_opt=%s", g.genParams()))
	}

	args = append(args, g.ProtocFlags...)
	args = append(args, g.Walk()...)
	protoCmd := exec.Command(g.Protoc, args...)
	protoCmd.Stderr = os.Stderr
	log.Println(protoCmd.String())
	return protoCmd.Run()
}
