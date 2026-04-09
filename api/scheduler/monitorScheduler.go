package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	monitorDao "dodevops-api/api/monitor/dao"
	monitorService "dodevops-api/api/monitor/service"
	"dodevops-api/api/monitor/model"

	"github.com/robfig/cron/v3"
)

type MonitorScheduler struct {
	cron         *cron.Cron
	mutex        sync.Mutex
	running      bool
	domainJobs   map[uint]cron.EntryID
	domainSpecs  map[uint]string
}

func NewMonitorScheduler() *MonitorScheduler {
	return &MonitorScheduler{
		cron:        cron.New(),
		domainJobs:  make(map[uint]cron.EntryID),
		domainSpecs: make(map[uint]string),
	}
}

func (m *MonitorScheduler) Start() error {
	m.mutex.Lock()
	if m.running {
		m.mutex.Unlock()
		return nil
	}

	if _, err := m.cron.AddFunc("@every 2m", func() {
		if err := monitorService.NewMonitorAutomationService().RunHostAlertScan(context.Background()); err != nil {
			log.Printf("host alert scan failed: %v", err)
		}
	}); err != nil {
		m.mutex.Unlock()
		return err
	}

	if _, err := m.cron.AddFunc("@every 3m", func() {
		if err := monitorService.NewMonitorAutomationService().RunDBAlertScan(context.Background()); err != nil {
			log.Printf("database alert scan failed: %v", err)
		}
	}); err != nil {
		m.mutex.Unlock()
		return err
	}

	if _, err := m.cron.AddFunc("@every 1m", func() {
		m.reloadDomainSchedules()
	}); err != nil {
		m.mutex.Unlock()
		return err
	}

	m.cron.Start()
	m.running = true
	m.mutex.Unlock()
	m.reloadDomainSchedules()
	log.Println("监控深化调度器启动成功")
	return nil
}

func (m *MonitorScheduler) Stop() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if !m.running {
		return
	}
	m.cron.Stop()
	m.cron = cron.New()
	m.domainJobs = make(map[uint]cron.EntryID)
	m.domainSpecs = make(map[uint]string)
	m.running = false
}

func (m *MonitorScheduler) reloadDomainSchedules() {
	schedules, err := monitorDao.GetEnabledMonitorDomainSchedules()
	if err != nil {
		log.Printf("reload domain schedules failed: %v", err)
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	active := make(map[uint]model.MonitorDomainScheduleEntity, len(schedules))
	for _, item := range schedules {
		if !item.Enabled || item.CronExpr == "" {
			continue
		}
		active[item.ID] = item
	}

	for id, entryID := range m.domainJobs {
		item, ok := active[id]
		if !ok || m.domainSpecs[id] != item.CronExpr {
			m.cron.Remove(entryID)
			delete(m.domainJobs, id)
			delete(m.domainSpecs, id)
		}
	}

	for _, item := range schedules {
		if !item.Enabled || item.CronExpr == "" {
			continue
		}
		if spec, ok := m.domainSpecs[item.ID]; ok && spec == item.CronExpr {
			continue
		}
		scheduleID := item.ID
		entryID, addErr := m.cron.AddFunc(item.CronExpr, func() {
			if err := monitorService.NewMonitorAutomationService().RunSSLDomainScan(context.Background()); err != nil {
				log.Printf("ssl domain scan failed(schedule=%d): %v", scheduleID, err)
			}
		})
		if addErr != nil {
			log.Printf("add ssl domain schedule failed(id=%d): %v", item.ID, addErr)
			continue
		}
		m.domainJobs[item.ID] = entryID
		m.domainSpecs[item.ID] = item.CronExpr
		if next := m.cron.Entry(entryID).Next; !next.IsZero() {
			nextRun := next
			lastRun := time.Now()
			_ = monitorDao.UpdateMonitorDomainScheduleRuntime(item.ID, "running", &lastRun, &nextRun)
		}
	}
}
