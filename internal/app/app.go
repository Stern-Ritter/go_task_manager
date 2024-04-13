package app

import (
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Stern-Ritter/go_task_manager/internal/config"
	"github.com/Stern-Ritter/go_task_manager/internal/service"
	"github.com/Stern-Ritter/go_task_manager/internal/storage"
	"github.com/go-chi/chi"
	"go.uber.org/zap"

	_ "modernc.org/sqlite"
)

func Run(config *config.ServerConfig, logger *zap.Logger) error {
	appPath, err := os.Getwd()
	if err != nil {
		logger.Fatal(err.Error(), zap.String("event", "get absolute path for current process"))
		return err
	}

	db, err := sql.Open(config.DatabaseDriverName, config.DatabaseFile)
	if err != nil {
		logger.Fatal(err.Error(), zap.String("event", "open database connection"))
		return err
	}
	defer db.Close()

	needInitDataBase := !isDatabaseExists(appPath, config.DatabaseFile)
	if needInitDataBase {
		err = initDatabase(db, filepath.Join(appPath, "/resources/database/init.sql"))
		if err != nil {
			logger.Fatal(err.Error(), zap.String("event", "init database schema"))
			return err
		}
	}

	authService := service.NewAuthService(config.RootPassword, logger)
	taskStore := storage.NewTaskStore(db)
	taskService := service.NewTaskService(taskStore, logger)
	server := service.NewServer(authService, taskService, config, logger)

	url := strings.Join([]string{"", strconv.Itoa(config.Port)}, ":")
	r := addRoutes(server, appPath)
	err = http.ListenAndServe(url, r)
	if err != nil {
		server.Logger.Fatal(err.Error(), zap.String("event", "start server"))
	}
	return err
}

func addRoutes(s *service.Server, appPath string) *chi.Mux {
	r := chi.NewRouter()
	filesDir := http.Dir(filepath.Join(appPath, "web"))
	fileServer(r, "/", filesDir)

	r.Route("/api", func(r chi.Router) {
		r.Post("/signin", s.SignInHandler)
		r.Get("/nextdate", s.GetNextDateHandler)

		r.Route("/tasks", func(r chi.Router) {
			r.Use(s.AuthMiddleware)
			r.Get("/", s.GetTasksHandler)
		})

		r.Route("/task", func(r chi.Router) {
			r.Use(s.AuthMiddleware)
			r.Get("/", s.GetTaskHandler)
			r.Post("/", s.AddTaskHandler)
			r.Put("/", s.UpdateTaskHandler)
			r.Delete("/", s.DeleteTaskHandler)
			r.Post("/done", s.CompleteTaskHandler)
		})
	})
	return r
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func isDatabaseExists(path string, dataBaseName string) bool {
	dbFile := filepath.Join(path, dataBaseName)
	_, err := os.Stat(dbFile)
	return err == nil
}

func initDatabase(db *sql.DB, initScriptPath string) error {
	data, err := os.ReadFile(initScriptPath)
	if err != nil {
		return err
	}

	initScript := string(data)
	_, err = db.Exec(initScript)
	return err
}
