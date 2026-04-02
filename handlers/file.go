package handlers

import (
	"account_manager/proto/pb"
	"account_manager/servers"
	"github.com/zehongyang/bee"
	"github.com/zehongyang/bee/logger"
	"net/http"
)

func FileUploadQuery() bee.Handler {
	fs := servers.GetFileServer()
	return func(ctx bee.IContext) {
		var h Header
		err := ctx.BindHeader(&h)
		if err != nil || len(h.FileHash) < 1 {
			logger.Error().Err(err).Any("h", h).Msg("FileUploadQuery")
			ctx.ResponseError(http.StatusBadRequest)
			return
		}
		fh, err := ctx.FormFile("file")
		if err != nil {
			logger.Error().Err(err).Any("h", h).Msg("FileUploadQuery")
			ctx.ResponseError(http.StatusBadRequest)
			return
		}
		fn, err := fs.Upload(fh, h.FileHash)
		if err != nil {
			logger.Error().Err(err).Any("h", h).Msg("FileUploadQuery")
			ctx.ResponseError(http.StatusInternalServerError)
			return
		}
		ctx.ResponseOk(&pb.FileUploadQueryResponse{FileName: fn})
	}
}
