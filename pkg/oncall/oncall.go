package oncall

import (
	"fmt"
	"math"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/client"
	"github.com/rickar/cal/v2"
	"github.com/rickar/cal/v2/pl"
)

type adminOnDuty struct {
	workday, weekend, holiday int
	profit                    int
	name                      string
	start                     time.Time
	end                       time.Time
}

// AdminOnDutyList Slice of admin on duty
type AdminOnDutyList struct {
	AdminsOnDuty []adminOnDuty
	OnCall       client.OnCallClient
	Options      pagerduty.ListOnCallOptions
}

type dutyPay struct {
	workday, weekend, holiday int
}

// UsersOnCallOptions Set options for PagderDuty
func (u *AdminOnDutyList) UsersOnCallOptions(dutyStartDate, dutyEndDate string, scheduleID string) {
	u.Options.Since = dutyStartDate
	u.Options.Until = dutyEndDate
	u.Options.More = true
	u.Options.Limit = 100
	u.Options.ScheduleIDs = append(u.Options.ScheduleIDs, scheduleID)
}

func (users *AdminOnDutyList) monthEdges(id int, start, end time.Time) {
	pdMonthStart := users.AdminsOnDuty[id].start.Month()
	pdMonthEnd := users.AdminsOnDuty[id].end.Month()
	monthStart := start.Month()
	monthEnd := end.Month()
	durationMonthStart := int(pdMonthStart) - int(monthStart)
	durationMonthEnd := int(pdMonthEnd) - int(monthEnd)
	if durationMonthStart < 0 {
		users.AdminsOnDuty[id].start = start
	}
	if durationMonthEnd > 0 {
		users.AdminsOnDuty[id].end = end.AddDate(0, 0, 1)
	}
}

func (users *AdminOnDutyList) typeOfDay(id int) {
	c := cal.NewBusinessCalendar()
	c.AddHoliday(
		pl.NewYear,
		pl.ThreeKings,
		pl.EasterMonday,
		pl.LabourDay,
		pl.ConstitutionDay,
		pl.CorpusChristi,
		pl.AssumptionBlessedVirginMary,
		pl.AllSaints,
		pl.NationalIndependenceDay,
		pl.ChristmasDayOne,
		pl.ChristmasDayTwo,
	)

	dutyDurationDays := users.durationDays(id)
	dayCounter := users.AdminsOnDuty[id].start
	for i := 0; i < int(dutyDurationDays); i++ {
		if hol, _, _ := c.IsHoliday(dayCounter); hol {
			users.AdminsOnDuty[id].holiday++
		} else if c.IsWorkday(dayCounter) {
			users.AdminsOnDuty[id].workday++
		} else if cal.IsWeekend(dayCounter) {
			users.AdminsOnDuty[id].weekend++
		}
		dayCounter = dayCounter.AddDate(0, 0, 1)
	}
}

func (users *AdminOnDutyList) durationDays(id int) float64 {
	dutyStartDate := users.AdminsOnDuty[id].start
	dutyEndDate := users.AdminsOnDuty[id].end
	duration := dutyEndDate.Sub(dutyStartDate).Hours() / 24
	roundDuration := math.Round(duration*10) / 10
	return roundDuration
}

func (users *AdminOnDutyList) getPdUserID(name string) int {
	for id, user := range users.AdminsOnDuty {
		if user.name == name {
			return id
		}
	}
	return -1
}

func (users *AdminOnDutyList) defineUserParameter(user pagerduty.OnCall, userID *int) {
	var pduser adminOnDuty
	loc, _ := time.LoadLocation("Europe/Warsaw")
	if *userID == -1 {
		pduser.name = user.User.Summary
		pduser.start, _ = time.ParseInLocation(time.RFC3339, user.Start, loc)
		pduser.end, _ = time.ParseInLocation(time.RFC3339, user.End, loc)
		users.AdminsOnDuty = append(users.AdminsOnDuty, pduser)
		*userID = len(users.AdminsOnDuty) - 1
	}
	users.AdminsOnDuty[*userID].start, _ = time.ParseInLocation(time.RFC3339, user.Start, loc)
	users.AdminsOnDuty[*userID].end, _ = time.ParseInLocation(time.RFC3339, user.End, loc)
}

// UsersOnCall operation on users on call
func (u *AdminOnDutyList) UsersOnCall(start, end time.Time) error {
	users, err := u.OnCall.ListOnCalls(u.Options)
	if err != nil {
		return err
	}
	for _, user := range users.OnCalls {
		// PagerDuty has multiple level of escalation for single schedule
		// With all accessible levels PagerDuty return users multiple times
		if user.EscalationLevel == 1 {
			userID := u.getPdUserID(user.User.Summary)
			u.defineUserParameter(user, &userID)
			u.monthEdges(userID, start, end)
			u.typeOfDay(userID)
		}
	}
	return nil
}

// PrintDutySummary duty summary
func (users *AdminOnDutyList) PrintDutySummary(profit bool) (strs string) {
	var total int
	users.DutyUsersProfits()
	for _, user := range users.AdminsOnDuty {
		total = total + user.workday + user.holiday + user.weekend
		if profit {
			fmt.Println("Name: ", user.name, " workdays: ", user.workday, " holidays: ", user.holiday, " weekend days: ", user.weekend, " profit: ", user.profit)

			strstmp := fmt.Sprintf("Name: %s workdays: %d holidays: %d weekend days: %d profit: %d \n", user.name, user.workday, user.holiday, user.weekend, user.profit)
			strs = strs + strstmp
		} else {
			fmt.Println("Name: ", user.name, " workdays: ", user.workday, " holidays: ", user.holiday, " weekend days: ", user.weekend)
		}
	}
	//fmt.Println("Total days: ", total)
	return strs
}

// DutyUsersProfits additional profits info
func (users *AdminOnDutyList) DutyUsersProfits() {
	p := dutyPay{workday: 100, weekend: 180, holiday: 270}
	for userID, user := range users.AdminsOnDuty {
		users.AdminsOnDuty[userID].profit = user.holiday*p.holiday + user.workday*p.workday + user.weekend*p.weekend
	}
}
