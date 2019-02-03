package main

import (
	"gopkg.in/doug-martin/goqu.v5"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	_ "gopkg.in/doug-martin/goqu.v5/adapters/postgres"
	_ "gopkg.in/doug-martin/goqu.v5/adapters/mysql"
	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/xo/dburl"
	"os"
	"encoding/json"
	"strconv"
)

type Service struct {
	config           AppConfig
	db               *sql.DB
	database         *goqu.Database
	logger           *log.Logger
	moduleRepository *ModuleRepository
}

func (s *Service) SetLogger(logger *log.Logger) {
	s.logger = logger
}

func (s *Service) createDefaultLogger() *log.Logger {
	return log.New(os.Stdout, "module-server", log.Ldate)
}

func (s *Service) GetLogger() *log.Logger {
	if s.logger == nil {
		newLogger := s.createDefaultLogger()
		s.SetLogger(newLogger)
		return newLogger
	} else {
		return s.logger
	}
}

func (s *Service) init() {
	var err error

	u, err := dburl.Parse(s.config.DatabaseUrl)
	if err != nil {
		panic(err)
	}

	s.db, err = sql.Open(u.Driver, u.DSN)
	if err != nil {
		panic(err.Error())
	}
	s.database = goqu.New(u.Driver, s.db)

	s.moduleRepository = NewModuleRepository(s.database)
}

func NewService(config AppConfig) *Service {
	srv := &Service{config: config}
	srv.init()
	return srv
}

func (s *Service) createRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", s.HealthHandler)

	r.HandleFunc("/modules", s.ModuleListHandler).Methods("GET")
	r.HandleFunc("/modules/{id:[0-9]+}", s.ModuleViewHandler).Methods("GET")
	r.HandleFunc("/modules/history/{pkg}", s.ModuleHistoryHandler).Methods("GET")
	r.HandleFunc("/modules/check", s.ModuleCheckHandler).Methods("POST")

	return r
}

func (s *Service) Listen() {
	s.GetLogger().Fatal(http.ListenAndServe(":8000", s.createRouter()))
}

func (s *Service) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Сделать механизм для обновления и доставки новых модулей. На вход от клиента ожидаем id модуля.
// В ответ выдаём модуль. (Придумать структуру ответа).
func (s *Service) ModuleViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{err})
		return
	}

	module, err := s.moduleRepository.Find(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{err})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.createResponseFromModule(module))
}

func (s *Service) createResponseFromModule(m Module) ModuleRequest {
	return ModuleRequest{
		Data: m.Data,
		Meta: Meta{
			ID:      m.ID,
			Version: m.Version,
			Package: m.Package,
			Extras:  m.Extras,
		},
	}
}

func (s *Service) ModuleHistoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	modules, err := s.moduleRepository.FindHistory(vars["pkg"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{err})
		return
	}

	var deviceResponse []ModuleRequest
	for _, m := range modules {
		deviceResponse = append(deviceResponse, s.createResponseFromModule(m))
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deviceResponse)
}

func (s *Service) ModuleListHandler(w http.ResponseWriter, r *http.Request) {
	modules, err := s.moduleRepository.FindAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{err})
		return
	}

	var deviceResponse []ModuleRequest
	for _, m := range modules {
		deviceResponse = append(deviceResponse, s.createResponseFromModule(m))
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deviceResponse)
}

// Android приложение раз в N минут обращается к серверу передавая следующие данные:
//{
//	"deviceId": "uniqueDeviceID",
//	"installedModules": [
//		{
//			"id": "ads",
//			"version": 8
//		},
//		{
//			"id": "proxy",
//			"version": 3
//		}
//	]
//}
//
//Сравнить клиентские версии модулей с рабочими, если они отличается – отправить клиенту в
// ответ список с idами модулей, которые нужно обновить.
func (s *Service) ModuleCheckHandler(w http.ResponseWriter, r *http.Request) {
	var data DeviceRequest
	json.NewDecoder(r.Body).Decode(&data)

	deviceModules, err := s.moduleRepository.FindAllByDevice(data.DeviceID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{err})
		return
	}

	for _, module := range deviceModules {
		for _, installed := range data.InstalledModules {
			if module.ID != installed.ID {
				continue
			}

			if module.Version == installed.Version {
				continue
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(deviceModules)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
