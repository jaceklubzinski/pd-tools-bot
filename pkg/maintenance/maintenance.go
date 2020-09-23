package maintenance

import (
	"fmt"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/client"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/extensions"
)

//Maintenances PagerDuty client
type Maintenances struct {
	Maintenance client.MaintenanceClient
}

//Get List maintenace windows from PagerDuty
func (s *Maintenances) Get(teamID []string) (strs string, err error) {
	opts := pagerduty.ListMaintenanceWindowsOptions{
		Filter:  "ongoing",
		TeamIDs: teamID,
	}
	getMaintenance, err := s.Maintenance.ListMaintenance(opts)
	if err != nil {
		return "", err
	}
	for _, p := range getMaintenance.MaintenanceWindows {
		strstmp := fmt.Sprintf("Services: %s Start Time: %s End Time: %s Description: %s\n", p.APIObject.Summary, p.StartTime, p.EndTime, p.Description)
		strs = strs + strstmp
	}
	return strs, nil
}

//Create maintanence window for PagerDuty service
func (s *Maintenances) Create(serviceID string, addHour string) (strs string, err error) {
	startToday := time.Now().Format("2006-01-02 15:04:05")
	endHours, err := extensions.AddDurationToDate(startToday, addHour)
	if err != nil {
		return "", err
	}
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
	m, err := s.Maintenance.CreateMaintenance(opts)
	if err != nil {
		return "", err
	}
	strs = fmt.Sprintf("%s - %s - %s", m.Services, m.StartTime, m.EndTime)
	return strs, nil
}
