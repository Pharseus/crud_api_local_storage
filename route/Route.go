package route

import (
	"crud-api-local-storage/controllers"
	_ "crud-api-local-storage/docs" // Import docs yang akan di-generate
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func StartRoute() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API BERHASIL - Swagger: http://localhost:8080/swagger/index.html"))
	})

	r.Route("/files", func(r chi.Router) {
		r.Post("/upload", controllers.UploadFile)
		r.Get("/view/{path}", controllers.DownloadFile)
		r.Get("/", controllers.ListFiles)
		r.Delete("/{path}", controllers.DeleteFiles)
	})

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Swagger UI: http://localhost:8080/swagger/index.html")
	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		panic(err)
	}
}
