package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
	"github.com/olebedev/when/rules/ru"
)

const (
	taskName  = "TERMINAL_REMINDER"
	taskValue = "1"
)

var (
	args    []string = os.Args
	argPath string   = args[0]
	argTime string
	now     time.Time = time.Now()
)

func printErrorMessage(errMessage string, parameters ...string) {
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	fmt.Printf(errMessage+"\n", parameters)
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
}

func exitIfArgumentsMissing() {
	if len(args) < 3 {
		printErrorMessage("Usage: %s <hh:mm> <text message>", argPath)
		os.Exit(1)
	}
}

func exitIfCannotParse(err error) {
	if err != nil {
		printErrorMessage(err.Error())
		os.Exit(2)
	}
}

func exitIfTimeIsNil(parsedTime *when.Result) {
	if parsedTime == nil {
		printErrorMessage("Cannot parse time from", args...)
		os.Exit(2)
	}
}

func exitIfREminderTimeIsInThePast(parsedTime *when.Result, now time.Time) {
	if now.After(parsedTime.Time) {
		printErrorMessage("Cannot schedule a reminder for a past time [%t], now: [%t]",
			parsedTime.Time.GoString(), now.GoString())
		os.Exit(3)
	}
}

func exitIfCannotSchedule(err error) {
	if err != nil {
		printErrorMessage(err.Error())
		os.Exit(5)
	}
}

func exitIfCommandCannotBeStarted(command *exec.Cmd) {
	if err := command.Start(); err != nil {
		fmt.Println(err)
		os.Exit(4)
	}
}

func main() {
	exitIfArgumentsMissing()
	argTime = args[1]

	whenParser := when.New(nil)
	whenParser.Add(en.All...)
	whenParser.Add(ru.All...)
	whenParser.Add(common.All...)

	parsedTime, err := whenParser.Parse(argTime, now)
	exitIfCannotParse(err)
	exitIfTimeIsNil(parsedTime)
	exitIfREminderTimeIsInThePast(parsedTime, now)

	durationTillReminder := parsedTime.Time.Sub(now)
	if os.Getenv(taskName) == taskValue {
		time.Sleep(durationTillReminder)
		err := beeep.Alert("Reminder", strings.Join(args[2:], " "), "assets/information.png")
		exitIfCannotSchedule(err)
	} else {
		command := exec.Command(argPath, args[1:]...)
		command.Env = append(os.Environ(), fmt.Sprintf("%s=%s", taskName, taskValue))
		exitIfCommandCannotBeStarted(command)
		fmt.Println("Reminder will be displayed after", durationTillReminder.Round(time.Second), "seconds")
		os.Exit(0)
	}
}
