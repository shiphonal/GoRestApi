package save

import (
	res "GoServise/internal/lib/api"
	"GoServise/internal/lib/random"
	"GoServise/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type UrlSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

type Response struct {
	*res.Ans
	Alias string `json:"alias,omitempty"`
}

type Request struct {
	Url   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

func New(log *slog.Logger, save UrlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("fail with decoding json", slog.Any("err", err))
			render.JSON(w, r, res.Error("failed decode request"))
			return
		}
		log.Info("successfully decode request", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {
			var validErr validator.ValidationErrors
			errors.As(err, &validErr)
			log.Error("fail with validation error", slog.Any("err", err))
			render.JSON(w, r, res.Validation(validErr))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(res.LengthAlias)
		}

		id, err := save.SaveURL(req.Url, alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Error("url is already exists", slog.String("url", req.Url), slog.Any("err", err))
				render.JSON(w, r, res.Error("url is already exists"))
			} else {
				log.Error("fail with saving url", slog.String("url", req.Url), slog.Any("err", err))
				render.JSON(w, r, res.Error("failed to save url"))
			}
			return
		}

		log.Info("successfully saving url", slog.Int64("id", id))
		render.JSON(w, r, Response{
			Ans:   res.Ok(),
			Alias: alias,
		})
	}
}
