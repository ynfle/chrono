package chrono

import "time"

type TriggerContext interface {
	LastCompletionTime() time.Time
	LastExecutionTime() time.Time
	LastTriggeredExecutionTime() time.Time
}

type SimpleTriggerContext struct {
	lastCompletionTime         time.Time
	lastExecutionTime          time.Time
	lastTriggeredExecutionTime time.Time
}

func NewSimpleTriggerContext() *SimpleTriggerContext {
	return &SimpleTriggerContext{}
}

func (ctx *SimpleTriggerContext) Update(lastCompletionTime time.Time, lastExecutionTime time.Time, lastTriggeredExecutionTime time.Time) {
	ctx.lastCompletionTime = lastCompletionTime
	ctx.lastExecutionTime = lastExecutionTime
	ctx.lastTriggeredExecutionTime = lastTriggeredExecutionTime
}

func (ctx *SimpleTriggerContext) LastCompletionTime() time.Time {
	return ctx.lastCompletionTime
}

func (ctx *SimpleTriggerContext) LastExecutionTime() time.Time {
	return ctx.lastExecutionTime
}

func (ctx *SimpleTriggerContext) LastTriggeredExecutionTime() time.Time {
	return ctx.lastTriggeredExecutionTime
}

type Trigger interface {
	NextExecutionTime(ctx TriggerContext) time.Time
}

type CronTrigger struct {
	cronExpression *CronExpression
	location       *time.Location
}

func CreateCronTrigger(expression string, location *time.Location) (*CronTrigger, error) {
	cron, err := ParseCronExpression(expression)

	if err != nil {
		return nil, err
	}

	trigger := &CronTrigger{
		cron,
		time.Local,
	}

	if location != nil {
		trigger.location = location
	}

	return trigger, nil
}

func (trigger *CronTrigger) NextExecutionTime(ctx TriggerContext) time.Time {
	now := time.Now()
	lastCompletion := ctx.LastCompletionTime()

	if !lastCompletion.IsZero() {

		lastExecution := ctx.LastTriggeredExecutionTime()

		if !lastExecution.IsZero() && now.Before(lastExecution) {
			now = lastExecution
		}

	}

	originalLocation := now.Location()

	convertedTime := now.In(trigger.location)
	newTime := time.Date(convertedTime.Year(),
		convertedTime.Month(),
		convertedTime.Day(),
		convertedTime.Hour(),
		convertedTime.Minute(),
		convertedTime.Second(),
		0,
		trigger.location)

	next := trigger.cronExpression.NextTime(newTime)
	return next.In(originalLocation)
}
