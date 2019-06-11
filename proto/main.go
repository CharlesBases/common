package main

import (
	"bytes"
	"flag"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"

	log "github.com/cihub/seelog"
	"golang.org/x/tools/imports"

	"github.com/CharlesBases/common/proto/parse"
)

var (
	goFile  = flag.String("file", "/Users/sun/go/SourceCode/src/github.com/CharlesBases/common/proto/bll.go", "full path of the file")
	genPath = flag.String("path", "./proto/pb/", "full path of the generate folder")
	pkgName = flag.String("package", "pb", "package name in .proto file")
)

var (
	serFile = "pb.server.go"
	cliFile = "pb.client.go"
	proFile = "pb.proto"
)

func init() {
	logger, _ := log.LoggerFromConfigAsString(`
			<?xml version="1.0" encoding="utf-8" ?>
			<seelog levels="info,error">
				<outputs formatid="main">
					<filter levels="warn,info">
						<console formatid="main"/>
					</filter>
					<filter levels="error,critical">
						<console formatid="error"/>
					</filter>
				</outputs>
				<formats>
					<format id="main" format="[%Date(2006-01-02 15:04:05.000)][%LEV] ==&gt; %Msg%n"/>
					<format id="error" format="%EscM(31)[%Date(2006-01-02 15:04:05.000)][%LEV] ==&gt; %Msg%n%EscM(0)"/>
				</formats>
			</seelog>`)
	log.UseLogger(logger)
}

func main() {
	defer log.Flush()

	flag.Parse()

	serFile = path.Join(*genPath, serFile)
	cliFile = path.Join(*genPath, cliFile)
	proFile = path.Join(*genPath, proFile)

	os.MkdirAll(*genPath, 0666)

	if *goFile == "" {
		_, file, _, _ := runtime.Caller(0)
		goFile = &file
	}

	log.Info("parsing files for go: ", *goFile)

	var infor *parse.File

	fileSet := token.NewFileSet()
	astFile, err := parser.ParseFile(fileSet, *goFile, nil, 0)
	if err != nil {
		log.Error(err)
		return
	}
	gofile := parse.NewFile(*pkgName, "github.com/CharlesBases/common/proto/")
	gofile.ParseFile(astFile)
	if len(gofile.Interfaces) == 0 {
		return
	}
	gofile.ParsePkgStruct(&parse.Package{PkgPath: func() string {
		path, _ := os.Getwd()
		return path
	}()})
	infor = &gofile

	// generate proto file
	profile, err := createFile(proFile)
	if err != nil {
		log.Error(err)
		return
	}
	defer profile.Close()
	infor.GenProtoFile(profile)

	log.Info("run the protoc command ...")
	dir := filepath.Dir(proFile)
	out, err := exec.Command("protoc", "--proto_path="+dir+"/", "--gogo_out=plugins=grpc:"+dir+"/", proFile).CombinedOutput()
	if err != nil {
		log.Error("protoc error: ", string(out))
		return
	}
	log.Info("protoc complete !")

	infor.GoTypeConfig()

	// generate server file
	serfile, err := createFile(serFile)
	if err != nil {
		log.Error(err)
		return
	}
	defer serfile.Close()
	bufferSer := bytes.NewBuffer([]byte{})
	infor.GenServer(bufferSer)
	serfile.Write(bufferSer.Bytes())
	byteSlice, e1 := imports.Process("", bufferSer.Bytes(), nil)
	if e1 != nil {
		log.Error(e1)
		return
	}
	serfile.Truncate(0)
	serfile.Seek(0, 0)
	serfile.Write(byteSlice)

	// generate client file
	clifile, err := createFile(cliFile)
	if err != nil {
		log.Error(err)
		return
	}
	defer clifile.Close()
	bufferCli := bytes.NewBuffer([]byte{})
	infor.GenClient(bufferCli)
	clifile.Write(bufferCli.Bytes())
	byteSlice, e := imports.Process("", bufferCli.Bytes(), nil)
	if e != nil {
		log.Error(e)
		return
	}
	clifile.Truncate(0)
	clifile.Seek(0, 0)
	clifile.Write(byteSlice)

	log.Info("complete!")
}

func createFile(fileName string) (*os.File, error) {
	os.RemoveAll(fileName)
	log.Info("create file: " + fileName)
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Error(err)
		return file, err
	}
	return file, nil
}
