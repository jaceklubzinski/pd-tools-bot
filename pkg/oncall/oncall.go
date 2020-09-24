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

//monthEdges calculate proper month profit for weekend between months
func (u *AdminOnDutyList) monthEdges(id int, start, end time.Time) {
	pdMonthStart := u.AdminsOnDuty[id].start.Month()
	pdMonthEnd := u.AdminsOnDuty[id].end.Month()
	monthStart := start.Month()
	monthEnd := end.Month()
	durationMonthStart := int(pdMonthStart) - int(monthStart)
	durationMonthEnd := int(pdMonthEnd) - int(monthEnd)
	if durationMonthStart < 0 {
		u.AdminsOnDuty[id].start = start
	}
	if durationMonthEnd > 0 {
		u.AdminsOnDuty[id].end = end.AddDate(0, 0, 1)
	}
}

//typeOfDay check type of day
func (u *AdminOnDutyList) typeOfDay(id int) {
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

	dutyDurationDays := u.durationDays(id)
	dayCounter := u.AdminsOnDuty[id].start
	for i := 0; i < int(dutyDurationDays); i++ {
		if hol, _, _ := c.IsHoliday(dayCounter); hol {
			u.AdminsOnDuty[id].holiday++
		} else if c.IsWorkday(dayCounter) {
			u.AdminsOnDuty[id].workday++
		} else if cal.IsWeekend(dayCounter) {
			u.AdminsOnDuty[id].weekend++
		}
		dayCounter = dayCounter.AddDate(0, 0, 1)
	}
}

//durationDays duration in hours for specific user duty
func (u *AdminOnDutyList) durationDays(id int) float64 {
	dutyStartDate := u.AdminsOnDuty[id].start
	dutyEndDate := u.AdminsOnDuty[id].end
	duration := dutyEndDate.Sub(dutyStartDate).Hours() / 24
	roundDuration := math.Round(duration*10) / 10
	return roundDuration
}

//getPdUserID user name do PagerDuty ID
func (u *AdminOnDutyList) getPdUserID(name string) int {
	for id, user := range u.AdminsOnDuty {
		if user.name == name {
			return id
		}
	}
	return -1
}

//defineUserParameter create user with PagerDuty data
func (u *AdminOnDutyList) defineUserParameter(user pagerduty.OnCall, userID *int) {
	var pduser adminOnDuty
	loc, _ := time.LoadLocation("Europe/Warsaw")
	if *userID == -1 {
		pduser.name = user.User.Summary
		pduser.start, _ = time.ParseInLocation(time.RFC3339, user.Start, loc)
		pduser.end, _ = time.ParseInLocation(time.RFC3339, user.End, loc)
		u.AdminsOnDuty = append(u.AdminsOnDuty, pduser)
		*userID = len(u.AdminsOnDuty) - 1
	}
	u.AdminsOnDuty[*userID].start, _ = time.ParseInLocation(time.RFC3339, user.Start, loc)
	u.AdminsOnDuty[*userID].end, _ = time.ParseInLocation(time.RFC3339, user.End, loc)
}

//UsersOnCall operation on users on call
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

//PrintDutySummary duty summary for current month
func (u *AdminOnDutyList) PrintDutySummary(profit bool) (strs string) {
	var total int
	u.DutyUsersProfits()
	for _, user := range u.AdminsOnDuty {
		total = total + user.workday + user.holiday + user.weekend
		if profit {
			strstmp := fmt.Sprintf("Name: %s workdays: %d holidays: %d weekend days: %d profit: %d \n", user.name, user.workday, user.holiday, user.weekend, user.profit)
			strs = strs + strstmp
		} else {
			fmt.Println("Name: ", user.name, " workdays: ", user.workday, " holidays: ", user.holiday, " weekend days: ", user.weekend)
		}
	}
	//fmt.Println("Total days: ", total)
	return strs
}

//PrintTodayDuty today person on duty
func (u *AdminOnDutyList) PrintTodayDuty(schedule string) (strs string) {
	for _, user := range u.AdminsOnDuty {
		strs = fmt.Sprintf("Schedule: `%s` On duty: `%s`\n", schedule, user.name)
	}
	return strs
}

//DutyUsersProfits additional profits info
func (u *AdminOnDutyList) DutyUsersProfits() {
	p := dutyPay{workday: 100, weekend: 180, holiday: 270}
	for userID, user := range u.AdminsOnDuty {
		u.AdminsOnDuty[userID].profit = user.holiday*p.holiday + user.workday*p.workday + user.weekend*p.weekend
	}
}
