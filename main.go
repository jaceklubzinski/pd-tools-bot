package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/client"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/extensions"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/incident"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/maintenance"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/oncall"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/schedule"
	"github.com/jaceklubzinski/pd-tools-bot/pkg/services"
	team "github.com/jaceklubzinski/pd-tools-bot/pkg/teams"
	"github.com/kelseyhightower/envconfig"
	"github.com/shomali11/slacker"
)

type envConfig struct {
	pagerdutyAuthToken string `required:"true" split_words:"true"`
	slackAuthToken     string `required:"true" split_words:"true"`
	dupa               string `required:"true" split_words:"true"`
}

func main() {

	var env envConfig
	if err := envconfig.Process("pdbot", &env); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("slack: ", env.slackAuthToken)
	pdclient := pagerduty.NewClient(env.pagerdutyAuthToken)
	conn := client.NewApiClient(pdclient)

	bot := slacker.NewClient(env.dupa)

	bot.Init(func() {
		log.Println("Connected!")
	})

	bot.Err(func(err string) {
		log.Println(err)
	})

	bot.DefaultCommand(func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		response.Reply("Say what? try `help` command")
	})

	bot.DefaultEvent(func(event interface{}) {
		fmt.Println(event)
	})
	oncallMonth := &slacker.CommandDefinition{
		Description: "PagerDuty oncall current month summary",
		Example:     "oncall month PCKO8FO",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			onCall := &oncall.AdminOnDutyList{OnCall: conn}
			pdschedule := request.StringParam("pdschedule", "PCKO8FO")
			today := time.Now()
			start := extensions.BeginningOfMonth(today)
			end := extensions.EndOfMonth(today)
			onCall.UsersOnCallOptions(start.String(), end.String(), pdschedule)
			if err := onCall.UsersOnCall(start, end); err != nil {
				response.ReportError(errors.New("Oops!"))
			}
			oncall := onCall.PrintDutySummary(true)
			response.Typing()
			time.Sleep(time.Second)
			response.Reply(oncall)
		},
	}

	incidentListTeam := &slacker.CommandDefinition{
		Description: "PagerDuty list of triggered and acknowledged incident for specific team",
		Example:     "incident list PHJN9RO",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			incident := incident.Incidents{Incident: conn}
			pdTeam := request.StringParam("pdteam", "PU7IVK3")
			pdTeamList := []string{}
			pdTeamList = append(pdTeamList, pdTeam)
			incidentOutls := incident.GetTeam(pdTeamList)
			response.Typing()
			time.Sleep(time.Second)
			response.Reply(incidentOutls)
		},
	}

	serviceListTeam := &slacker.CommandDefinition{
		Description: "PagerDuty list of services assigned to specific team",
		Example:     "service list PHJN9RO",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			service := &services.Services{Service: conn}
			pdTeam := request.StringParam("pdteam", "PU7IVK3")
			pdTeamList := []string{}
			pdTeamList = append(pdTeamList, pdTeam)
			serviceOutls := service.GetTeam(pdTeamList)
			response.Typing()
			time.Sleep(time.Second)
			response.Reply(serviceOutls)
		},
	}

	maintenanceListTeam := &slacker.CommandDefinition{
		Description: "PagerDuty list of service maintenance to specific team",
		Example:     "maintenace list PHJN9RO",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			maintenance := &maintenance.Maintenances{Maintenance: conn}
			pdTeam := request.StringParam("pdteam", "PU7IVK3")
			pdTeamList := []string{}
			pdTeamList = append(pdTeamList, pdTeam)
			maintenanceOutls := maintenance.GetMaintenance(pdTeamList)
			response.Typing()
			time.Sleep(time.Second)
			response.Reply(maintenanceOutls)
		},
	}

	maintenanceCreateTeam := &slacker.CommandDefinition{
		Description: "PagerDuty create maintenance for specific service",
		Example:     "maintenace create PHJN9RO 4h",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			maintenance := &maintenance.Maintenances{Maintenance: conn}
			pdServiceID := request.StringParam("pdservice", "PR1XCPX")
			toHour := request.StringParam("pdhour", "4h")
			maintenanceCreateOutls := maintenance.CreateMaintenance(pdServiceID, toHour)
			response.Reply(maintenanceCreateOutls)

		},
	}

	scheduleList := &slacker.CommandDefinition{
		Description: "PagerDuty schedule list",
		Example:     "schedule list",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			schedule := &schedule.NewSchedule{Schedule: conn}
			scheduleOutls := schedule.PrintSchedules()
			response.Typing()
			time.Sleep(time.Second)
			response.Reply(scheduleOutls)
		},
	}

	teamList := &slacker.CommandDefinition{
		Description: "PagerDuty team list",
		Example:     "team list",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			team := &team.NewTeam{Team: conn}
			teamOutls := team.PrintTeams()
			response.Typing()
			time.Sleep(time.Second)
			response.Reply(teamOutls)
		},
	}

	bot.Command("oncall month <pdschedule>", oncallMonth)
	bot.Command("incident list <pdteam>", incidentListTeam)
	bot.Command("schedule list", scheduleList)
	bot.Command("team list", teamList)
	bot.Command("service list <pdteam>", serviceListTeam)
	bot.Command("maintenance list <pdteam>", maintenanceListTeam)
	bot.Command("maintenance create <pdservice> <pdhour>", maintenanceCreateTeam)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
