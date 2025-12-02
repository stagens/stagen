package stagen

import (
	"context"
	"errors"
	"io"
	"os"
	"path"

	"github.com/gabriel-vasile/mimetype"
	"github.com/pixality-inc/golang-core/http"
	"github.com/pixality-inc/golang-core/storage"
	"github.com/valyala/fasthttp"
)

func (s *Impl) Web(ctx context.Context) error {
	router := http.NewRouter()

	router.GET("/__admin/{filepath:*}", s.httpAdminHandler)

	router.GET("/{filepath:*}", s.httpHandler)

	middlewares := []http.Middleware{
		s.httpNotFoundMiddleware,
		http.NewCorsMiddleware("*").Handle,
		http.NewRequestMetadataMiddleware().Handle,
		http.RequestLogHandler,
	}

	requestHandler := router.Handle()

	for _, middleware := range middlewares {
		requestHandler = middleware(requestHandler)
	}

	httpServer := http.New(
		"http",
		s.config.Http(),
		requestHandler,
	)

	if err := httpServer.ListenAndServe(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Impl) httpNotFoundMiddleware(originalHandler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		originalHandler(ctx)

		if ctx.Response.StatusCode() == fasthttp.StatusNotFound {
			ctx.Response.ResetBody()

			s.httpNotFound(ctx)
		}
	}
}

func (s *Impl) httpAdminHandler(ctx *fasthttp.RequestCtx) {
	httpPath, ok := ctx.Request.UserValue("filepath").(string)
	if !ok {
		s.httpBadRequest(ctx, "path is not a string", nil)

		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	_, _ = ctx.WriteString("Admin:\n" + httpPath) //nolint:errcheck
}

func (s *Impl) httpHandler(ctx *fasthttp.RequestCtx) {
	httpPath, ok := ctx.Request.UserValue("filepath").(string)
	if !ok {
		s.httpBadRequest(ctx, "path is not a string", nil)

		return
	}

	buildDir := s.buildDir()

	uri := "/" + httpPath

	filesToCheck := []string{
		path.Join(buildDir, uri),
		path.Join(buildDir, uri, "/index.html"),
	}

	if !s.config.Settings().UseUriHtmlFileExtension() {
		filesToCheck = append(filesToCheck, path.Join(buildDir, uri+".html"))
	}

	for index, filename := range filesToCheck {
		filesToCheck[index] = path.Clean(filename)
	}

	// @todo!!!!
	localStorage, ok := s.storage.(storage.LocalStorage)
	if !ok {
		s.httpInternalServerError(ctx, ErrStorageIsNotALocalStorage.Error(), ErrStorageIsNotALocalStorage)

		return
	}

	for _, filename := range filesToCheck {
		// @todo!!!!
		filePath, err := localStorage.LocalPath(ctx, filename)
		if err != nil {
			s.httpInternalServerError(ctx, err.Error(), err)

			return
		}

		file, err := os.Open(filePath)
		if errors.Is(err, os.ErrNotExist) {
			continue
		} else if err != nil {
			s.httpInternalServerError(ctx, "can't open file", err)

			return
		}

		stat, err := file.Stat()
		if err != nil {
			s.httpInternalServerError(ctx, "can't stat file", err)

			return
		}

		if stat.IsDir() {
			continue
		}

		s.httpHandleFile(ctx, filePath, file)

		return
	}

	ctx.SetStatusCode(fasthttp.StatusNotFound)
}

func (s *Impl) httpHandleFile(ctx *fasthttp.RequestCtx, filename string, file io.ReadCloser) {
	log := s.log.GetLogger(ctx)

	defer func() {
		if fErr := file.Close(); fErr != nil {
			log.WithError(fErr).Error("can't close file")
		}
	}()

	buf, err := io.ReadAll(file)
	if err != nil {
		s.httpInternalServerError(ctx, "can't read file", err)

		return
	}

	mimeType, err := mimetype.DetectFile(filename)
	if err != nil {
		s.httpInternalServerError(ctx, "can't detect mime type", err)

		return
	}

	contentType := mimeType.String()

	ext := path.Ext(filename)
	switch ext {
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "text/javascript"
	}

	s.httpOk(ctx, contentType, buf)
}

func (s *Impl) httpInternalServerError(ctx *fasthttp.RequestCtx, message string, err error) {
	// @todo
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)

	_, _ = ctx.WriteString("Internal Server Error:\n" + message + "\n" + err.Error()) //nolint:errcheck
}

func (s *Impl) httpBadRequest(ctx *fasthttp.RequestCtx, message string, err error) {
	// @todo 400
	ctx.SetStatusCode(fasthttp.StatusBadRequest)

	_, _ = ctx.WriteString("Bad Request:\n" + message) //nolint:errcheck
}

func (s *Impl) httpNotFound(ctx *fasthttp.RequestCtx) {
	// @todo 404
	ctx.SetStatusCode(fasthttp.StatusNotFound)

	_, _ = ctx.WriteString("Not Found") //nolint:errcheck
}

func (s *Impl) httpOk(ctx *fasthttp.RequestCtx, contentType string, content []byte) {
	ctx.Response.Header.Set("Content-Type", contentType)

	ctx.SetStatusCode(fasthttp.StatusOK)

	_, _ = ctx.Write(content) //nolint:errcheck
}
