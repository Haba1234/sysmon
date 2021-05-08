package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Haba1234/sysmon/internal/logger"
	"github.com/Haba1234/sysmon/internal/service/cpu"
	"github.com/Haba1234/sysmon/internal/service/loadaverage"
)

// StatisticsData подготовленные данные для отправки клиенту.
type StatisticsData struct {
	La      []string // La - load average.
	CPU     []string // CPU - cpu average.
	Counter int      // Counter - кол-во собранных данных.
}

type Collector struct {
	logg       *logger.Logger
	collection Collection
	la         *loadaverage.LoadAverage
	cpu        *cpu.CPU
	StatsData  StatisticsData
}

// Collection структура конфигурации.
type Collection struct {
	LoadAverageEnabled bool // LoadAverageEnabled - включение метрики 'load average'.
	CPUEnabled         bool // CPUEnabled - включение метрики 'cpu average'.
	BufSize            int  // BufSize - глубина истории собираемых метрик.
}

type writeStatistics interface {
	DataRequest() ([]float64, error)
	ShiftIndex()
	WriteValue(out []float64)
}

func NewCollector(logg *logger.Logger, collection Collection) *Collector {
	la := loadaverage.NewLoadAverage(collection.BufSize)
	cp := cpu.NewCPU(collection.BufSize)

	col := &Collector{
		logg:       logg,
		collection: collection,
		la:         la,
		cpu:        cp,
	}
	return col
}

// Start запуск сервиса.
func (col *Collector) Start(ctx context.Context) error {
	col.logg.Info("service 'collector' starting...")
	ticker := time.NewTicker(time.Second) // Сбор данных раз в секунду.
	col.StatsData.Counter = 0

	go func() {
		col.getDataStatistics()
		for {
			select {
			case <-ctx.Done():
				col.logg.Info("service 'collector' stopped")
				ticker.Stop()
				return
			case <-ticker.C:
				col.getDataStatistics()
				if col.StatsData.Counter < col.collection.BufSize {
					col.StatsData.Counter++
				}
			}
		}
	}()
	return nil
}

func (col *Collector) getDataStatistics() {
	if col.collection.LoadAverageEnabled {
		if err := col.addNewValue(col.la); err != nil {
			col.logg.Error(fmt.Sprintf("Service 'load average' - data generation error: %v", err))
		}
	}
	if col.collection.CPUEnabled {
		if err := col.addNewValue(col.cpu); err != nil {
			col.logg.Error(fmt.Sprintf("Service 'CPU average' - data generation error: %v", err))
		}
	}
}

func (col *Collector) addNewValue(i writeStatistics) error {
	out, err := i.DataRequest()
	if err != nil {
		return err
	}
	i.WriteValue(out)
	i.ShiftIndex()
	return nil
}

func (col *Collector) ReadStats(m int) {
	col.StatsData.La = col.la.GetValue(m)
	col.StatsData.CPU = col.cpu.GetValue(m)
}
