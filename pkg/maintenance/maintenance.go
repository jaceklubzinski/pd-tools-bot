package maintenance

import (
	"fmt"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/client"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/extensions"
)

type Maintenances struct {
	Maintenance client.MaintenanceClient
}

func (s *Maintenances) GetMaintenance(teamID []string) (strs string) {
	opts := pagerduty.ListMaintenanceWindowsOptions{
		Filter:  "ongoing",
		TeamIDs: teamID,
	}
	getMaintenance := s.Maintenance.ListMaintenance(opts)
	for _, p := range getMaintenance.MaintenanceWindows {
		strstmp := fmt.Sprintf("ID: %s Services: %s Start Time: %s End Time: %s Description: %s\n", p.APIObject.Summary, p.StartTime, p.EndTime, p.Description)
		strs = strs + strstmp
	}
	return strs
}

func (s *Maintenances) CreateMaintenance(serviceID string, addHour string) (strs string) {
	startToday := time.Now().Format("2006-01-02 15:04:05")
	endHours := extensions.AddDurationToDate(startToday, addHour)
	opts := pagerduty.MaintenanceWindow{
		StartTime: startToday,
		EndTime:   endHours,
		Services: []pagerduty.APIObject{
			pagerduty.APIObject{
				ID:   serviceID,
				Type: "service_reference",
			},
		},
	}
	m := s.Maintenance.CreateMaintenance(opts)
	strs = fmt.Sprintf("%s - %s - %s", m.Services, m.StartTime, m.EndTime)
	return strs
}
