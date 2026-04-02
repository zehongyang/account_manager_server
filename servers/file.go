package servers

import (
	"github.com/zehongyang/bee"
	"github.com/zehongyang/bee/config"
	"github.com/zehongyang/bee/logger"
	"github.com/zehongyang/bee/utils"
	"io"
	"mime/multipart"
	"os"
	"path"
)

type FileUploadConfig struct {
	Upload struct {
		Dir string
	}
}
type FileServer struct {
	cfg FileUploadConfig
}

var GetFileServer = utils.Single(func() *FileServer {
	var cfg FileUploadConfig
	err := config.Load(&cfg)
	if err != nil || len(cfg.Upload.Dir) == 0 {
		logger.Fatal().Err(err).Any("cfg", cfg).Msg("GetFileServer")
	}
	return &FileServer{cfg: cfg}
})

func (s *FileServer) Upload(fh *multipart.FileHeader, fid string) (string, error) {
	if fh == nil || len(fid) < 1 {
		return "", bee.ErrFileHeader
	}
	_, err := os.Stat(s.cfg.Upload.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(s.cfg.Upload.Dir, os.ModePerm)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	fn := path.Join(s.cfg.Upload.Dir, fid)
	_, err = os.Stat(fn)
	if err != nil {
		if os.IsNotExist(err) {
			tfn, err := os.Create(fn)
			if err != nil {
				return "", err
			}
			defer tfn.Close()
			src, err := fh.Open()
			if err != nil {
				return "", err
			}
			defer src.Close()
			_, err = io.Copy(tfn, src)
			if err != nil {
				return "", err
			}
		}
	}
	return fn, nil
}
