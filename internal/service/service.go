package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Haba1234/sysmon/internal/service/cpu"
	"github.com/Haba1234/sysmon/internal/service/loadaverage"
	"github.com/Haba1234/sysmon/internal/sysmon"
)

type collector struct {
	logg       sysmon.Logger
	collection Collection
	la         *sysmon.DataStore
	cpu        *sysmon.DataStore
}

// Collection структура конфигурации.
type Collection struct {
	LoadAverageEnabled bool // LoadAverageEnabled - включение метрики 'load average'.
	CPUEnabled         bool // CPUEnabled - включение метрики 'cpu average'.
	BufSize            int  // BufSize - глубина истории собираемых метрик.
}

// NewCollector конструктор.
func NewCollector(logg sysmon.Logger, collection Collection, f sysmon.Collectors) sysmon.Collector {
	col := &collector{
		logg:       logg,
		collection: collection,
		la:         loadaverage.NewLoadAverage(collection.BufSize, f.LoadAvg),
		cpu:        cpu.NewCPU(collection.BufSize, f.CPU),
	}
	return col
}

// Start запуск сервисов.
func (col *collector) Start(ctx context.Context) error {
	col.logg.Info("service 'collector' running...")
	if col.collection.LoadAverageEnabled {
		col.TaskGetDataAndSave(ctx, "load average", col.logg, col.la)
	}
	if col.collection.CPUEnabled {
		col.TaskGetDataAndSave(ctx, "CPU average", col.logg, col.cpu)
	}
	return nil
}

// GetStatusServices подготовка данных о состоянии сервисов.
func (col *collector) GetStatusServices() sysmon.StatusServices {
	return sysmon.StatusServices{
		La:  getStatus(col.la),
		CPU: getStatus(col.cpu),
	}
}

// GetStatus возврат текущего статуса сервиса.
func getStatus(ds *sysmon.DataStore) sysmon.Status {
	ds.Mu.Lock()
	defer ds.Mu.Unlock()
	return sysmon.Status{
		Counter:    ds.Counter,
		StatusCode: ds.StatusCode,
	}
}

// GetStats подготовка данных по сервисам для отправки клиенту.
func (col *collector) GetStats(m int) sysmon.StatisticsData {
	return sysmon.StatisticsData{
		La:  getValue(m, col.la),
		CPU: getValue(m, col.cpu),
	}
}

// getValue подсчет среднего значения за последние m секунд.
func getValue(m int, ds *sysmon.DataStore) sysmon.OutputData {
	str := make([]string, ds.Total)
	ds.Mu.Lock()
	defer ds.Mu.Unlock()

	if ds.Counter < m { // Данные еще не накопились.
		return sysmon.OutputData{}
	}

	avg := []float64{0.0, 0.0, 0.0}
	for y := 0; y < ds.Total; y++ {
		for i := 0; i < m; i++ {
			avg[y] += readValue(y, i, ds)
		}
		avg[y] /= float64(m)
		str[y] = strconv.FormatFloat(avg[y], 'f', 2, 64)
	}
	return sysmon.OutputData{
		Data:       str,
		Counter:    ds.Counter,
		StatusCode: ds.StatusCode,
	}
}

// readValue читает значение из кольцевого буфера.
func readValue(indexBuf, indexPos int, ds *sysmon.DataStore) float64 {
	var value float64
	laIndex := ds.Index - 1 //nolint:ifshort
	if laIndex < 0 {
		laIndex = ds.BufSize - 1
	}
	index := (indexPos + laIndex) % ds.BufSize
	value = ds.Stats[indexBuf][index]
	return value
}

// shiftIndex вычисляет новое значение индекса в кольцевом буфере.
func shiftIndex(index, counter, bufSize int) (int, int) {
	index++
	if index >= bufSize {
		index = 0
	}
	if counter < bufSize {
		counter++
	}
	return index, counter
}

func addNewValue(ds *sysmon.DataStore) error {
	result, err := ds.FuncDataReq()
	if err != nil {
		return err
	}
	ds.Mu.Lock()
	for i, v := range result {
		ds.Stats[i][ds.Index] = v
	}
	ds.Index, ds.Counter = shiftIndex(ds.Index, ds.Counter, ds.BufSize)
	ds.Mu.Unlock()
	return nil
}

func setStatusCode(ds *sysmon.DataStore, val int) {
	ds.Mu.Lock()
	ds.StatusCode = val
	ds.Mu.Unlock()
}

// TaskGetDataAndSave запускает сборку данных раз в секунду по таймеру.
func (col *collector) TaskGetDataAndSave(ctx context.Context, nameService string, logg sysmon.Logger, ds *sysmon.DataStore) {
	go func() {
		ds := ds
		col.logg.Info(fmt.Sprintf("service '%s' started", nameService))
		ticker := time.NewTicker(time.Second) // Сбор данных раз в секунду.
		setStatusCode(ds, sysmon.ServiceRun)
		for {
			select {
			case <-ctx.Done():
				logg.Info(fmt.Sprintf("service '%s' stopped", nameService))
				setStatusCode(ds, sysmon.ServiceStop)
				ticker.Stop()
				return
			case <-ticker.C:
				if err := addNewValue(ds); err != nil {
					logg.Error(fmt.Sprintf("stoped service '%s'. Data generation error: %v", nameService, err))
					setStatusCode(ds, sysmon.ServiceError)
					ticker.Stop()
					return
				}
			}
		}
	}()
}
