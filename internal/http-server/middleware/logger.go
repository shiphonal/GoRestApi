package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

// создает middleware для логирования HTTP-запросов
// Она принимает логгер, создает новый обработчик, который добавляет логирование перед и после вызова следующего обработчика в цепочке.
// Это позволяет легко добавлять и удалять middleware без изменения основного кода.
func NewLogger(log *slog.Logger) func(next http.Handler) http.Handler {

	// --------------------------------------------------------------------------------------- функция логирования обработчика HTTP-запросов ---------------------------------------------------------------------------------------
	return func(next http.Handler) http.Handler {
		log := log.With(slog.String("module", "middleware"))
		log.Info("logger middleware enabled")

		// w http.ResponseWriter - объект для записи ответов от сервера
		// r *http.Request - вопрос к серверу
		// логгирование обработчика и обработка следующего запроса (наш собственный интерфейс обработчика)
		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
			// Обёртка для объекта записи ответов, которая предоставляет новые поля для использования
			// r.ProtoMajor - протокол http (1, 2 и тд), берётся из информации о запросе
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			// после завершения обработки запроса выполнится логгирование
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					// Время, затраченное на обработку запроса
					slog.String("duration", time.Since(t1).String()),
				)
			}()

			// обрабатывает следующий запрос в цепочке, для его дальнейшего логгирования
			next.ServeHTTP(ww, r)
		}

		// возвращает объект обработчика, интерфейс которого мы написали
		// (он был наш собственный, а тут мы его обернули в интерфейс обработчика)
		return http.HandlerFunc(fn)
	}
}
