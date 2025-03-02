package remove

import (
	res "GoServise/internal/lib/api"
	"GoServise/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.remove.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("Missing alias")
			render.JSON(w, r, res.Error("missing alias"))
			return
		}

		err := urlDeleter.DeleteURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Error("URL not found", slog.Any("err", err))
				render.JSON(w, r, res.Error("URL not found"))
			} else {
				log.Error("Failed with delete url", slog.Any("err", err))
				render.JSON(w, r, res.Error("failed with delete url"))
			}
			return
		}

		log.Info("Removed URL", slog.Any("alias", alias))
		render.JSON(w, r, res.Ok())
	}
}
