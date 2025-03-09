package refactor

import (
	"GoServise/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	OldAlias string `json:"oldAlias"`
	NewAlias string `json:"newAlias"`
}

type URLPatcher interface {
	PatchURL(oldAlias string, newAlias string) (bool, error)
}

func New(log *slog.Logger, urlGetter URLPatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.refactor.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_is", middleware.GetReqID(r.Context())),
		)

		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("fail with decoding json", slog.Any("err", err))
			return
		}
		resURL, err := urlGetter.PatchURL(req.OldAlias, req.NewAlias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Error("URL not founded", slog.Any("err", err))
			}
			log.Error("fail with decoding json", slog.Any("err", err))
			return
		}
		render.JSON(w, r, resURL)
	}
}
