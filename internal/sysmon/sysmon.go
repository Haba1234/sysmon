package sysmon

import (
	"context"
	"sync"
)

const (
	ServiceRun   = 1
	ServiceStop  = 2
	ServiceError = 3
)

// DataStore хранилище статистики системы.
type DataStore struct {
	Mu          *sync.Mutex
	Stats       [][]float64     // Stats - слайс для хранения статистики.
	BufSize     int             // BufSize - максимальная глубина данных.
	Index       int             // Index - текущая позиция для новых значений.
	Counter     int             // Counter - кол-во собранных данных.
	StatusCode  int             // StatusCode - сервис: true (запущен), false (остановлен), 3 (остановлен с ошибкой).
	Total       int             // Total - кол-во считываемых значений метрики.
	FuncDataReq DataRequestFunc // FuncDataReq - функция, читающая определенную метрику системы.
}

// DataRequestFunc чтение данных метрики.
type DataRequestFunc func() ([]float64, error)

// Collectors - набор функций, возвращающих свои метрики. Передается для настройки сборщика.
type Collectors struct {
	LoadAvg DataRequestFunc
	CPU     DataRequestFunc
}

// OutputData структура данных по сервису для клиента.
type OutputData struct {
	Data       []string // Data - усредненные данные за заданный интервал времени.
	Counter    int      // Counter - кол-во собранных данных (максимум = BufSize).
	StatusCode int      // StatusCode - 1 (запущен), 2 (остановлен), 3 (остановлен с ошибкой).
}

// StatisticsData подготовленные данные для отправки клиенту.
type StatisticsData struct {
	La  OutputData // La - load average.
	CPU OutputData // CPU - cpu average.
}

// Status структура состояний сервиса.
type Status struct {
	Counter    int // Counter - кол-во собранных данных (максимум = BufSize).
	StatusCode int // StatusCode - 1 (запущен), 2 (остановлен), 3 (остановлен с ошибкой).
}

// StatusServices текущие состояния сервисов.
type StatusServices struct {
	La  Status // La состояние сервиса load average.
	CPU Status // CPU состояние сервиса CPU average.
}

// Logger интерфейсы.
type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg, pkg string)
}

// Collector интерфейсы.
type Collector interface {
	Start(ctx context.Context) error
	GetStats(m int) StatisticsData
	GetStatusServices() StatusServices
}

// Server gRPC интерфейсы.
type Server interface {
	Start(ctx context.Context, addr string) error
	Stop(ctx context.Context) error
}
