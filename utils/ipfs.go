package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/ipfs/kubo/core/coreiface/options"
	"io"
	"log"
)

func StoreFileIntoIpfs(context *gin.Context, file io.Reader) (string, error) {
	client, err := rpc.NewLocalApi()
	if err != nil {
		return "", err
	}
	ipfsPath, err := client.Unixfs().Add(context, files.NewReaderFile(file), options.Unixfs.Pin(true))
	if err != nil {
		return "", err
	}

	return ipfsPath.RootCid().String(), nil
}

func FetchFileFromIpfs(context *gin.Context, cidStr string) (files.File, error) {

	client, err := rpc.NewLocalApi()
	if err != nil {
		return nil, err
	}
	cidObj, err := cid.Parse(cidStr)
	if err != nil {
		return nil, err
	}
	ipfsPath, err := path.NewPath(path.FromCid(cidObj).String())
	if err != nil {
		return nil, err
	}

	fileNode, err := client.Unixfs().Get(context, ipfsPath)
	if err != nil {
		return nil, err
	}

	file, ok := fileNode.(files.File)
	if !ok {
		return nil, errors.New("ipfs file failed to get")
	}
	log.Printf("fetch file: %s", cidStr)
	return file, nil
}
