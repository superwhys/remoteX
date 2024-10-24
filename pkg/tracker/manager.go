package tracker

import (
	"time"

	"go.uber.org/atomic"
)

type Snapshot struct {
	UploadBlip   float64 `json:"upload_blip"`
	DownloadBlip float64 `json:"download_blip"`
}

type trafficInfo struct {
	uploadTemp   *atomic.Int64
	downloadTemp *atomic.Int64
	uploadBlip   *atomic.Int64
	downloadBlip *atomic.Int64
}

func initTraffic() *trafficInfo {
	return &trafficInfo{
		uploadTemp:   atomic.NewInt64(0),
		downloadTemp: atomic.NewInt64(0),
		uploadBlip:   atomic.NewInt64(0),
		downloadBlip: atomic.NewInt64(0),
	}
}

type Manager struct {
	*trafficInfo
}

func CreateTrackerManager() *Manager {
	m := &Manager{
		trafficInfo: initTraffic(),
	}
	go m.handleTrafficCalc()

	return m
}

func (m *Manager) PushUploaded(size int64) {
	m.uploadTemp.Add(size)
}

func (m *Manager) PushDownloaded(size int64) {
	m.downloadTemp.Add(size)
}

func (m *Manager) Now() (up, down int64) {
	return m.uploadBlip.Load(), m.downloadBlip.Load()
}

func (m *Manager) Snapshot() *Snapshot {
	return &Snapshot{
		UploadBlip:   float64(m.uploadBlip.Load()) / 1024,
		DownloadBlip: float64(m.downloadBlip.Load()) / 1024,
	}
}

func (m *Manager) ResetStatistic() {
	m.uploadTemp.Store(0)
	m.uploadBlip.Store(0)
	m.downloadTemp.Store(0)
	m.downloadBlip.Store(0)
}

func (m *Manager) handleTrafficCalc() {
	ticker := time.NewTicker(time.Second)

	for range ticker.C {
		m.uploadBlip.Store(m.uploadTemp.Load())
		m.uploadTemp.Store(0)
		m.downloadBlip.Store(m.downloadTemp.Load())
		m.downloadTemp.Store(0)
	}
}
