package redirect

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

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_is", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("Missing alias")
			render.JSON(w, r, res.Error("Missing alias"))
			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Error("URL not founded", slog.Any("err", err))
				render.JSON(w, r, res.Error("URL not founded"))
			} else {
				log.Error("Failed to get URL", slog.Any("err", err))
				render.JSON(w, r, res.Error("Failed to get URL"))
			}
			return
		}

		log.Info("Got URL", slog.String("url", resURL))
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
