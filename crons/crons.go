package crons

//Module exports cron dependencies
//var Module = fx.Options(
//	fx.Provide(NewScoreCron),
//	fx.Provide(NewCron),
//)
//
//type ICron interface {
//	Run() error
//	Name() string
//	RunAt() null.Time
//	IntervalHour() null.Int
//	SetJob(job *gocron.Job)
//	GetJob() *gocron.Job
//}
//
//type Cron struct {
//	logger    *internal.Logger
//	crons     []ICron
//	scheduler *gocron.Scheduler
//	location  *time.Location
//}
//
//func NewCron(
//	logger *internal.Logger,
//	env *internal.EnvConfig,
//
//	// add your new crons here
//	scoreCron *ScoreCron,
//) (*Cron, error) {
//	location, locErr := time.LoadLocation(env.Timezone)
//	if locErr != nil {
//		return nil, locErr
//	}
//	cron := Cron{
//		logger:    logger,
//		scheduler: gocron.NewScheduler(location),
//		location:  location,
//		crons: []ICron{
//			// add your new crons here
//			scoreCron,
//		},
//	}
//	return &cron, nil
//}
//
//// Setup sets up the cron jobs
//func (c *Cron) Setup() {
//	now := time.Now().In(c.location)
//	for i := range c.crons {
//		if c.crons[i].IntervalHour().Valid && !c.crons[i].RunAt().Valid {
//			// hourly cron runs
//			job, jobErr := c.scheduler.Every(c.crons[i].IntervalHour().Int).Hour().Tag(c.crons[i].Name()).Do(c.crons[i].Run)
//			if jobErr != nil {
//				c.logger.With(jobErr).Errorf("failed to assign cron job: %s", c.crons[i].Name())
//				continue
//			}
//			c.crons[i].SetJob(job)
//		} else if c.crons[i].IntervalHour().Valid && c.crons[i].RunAt().Valid {
//			// scheduled cron runs
//			scheduler := c.scheduler.Every(c.crons[i].IntervalHour().Int).Hour().At(c.crons[i].RunAt().Time.In(c.location))
//			job, jobErr := scheduler.Tag(c.crons[i].Name()).Do(c.crons[i].Run)
//			if jobErr != nil {
//				c.logger.With(jobErr).Errorf("failed to assign cron job: %s", c.crons[i].Name())
//				continue
//			}
//			c.crons[i].SetJob(job)
//		}
//	}
//	c.scheduler.StartAsync()
//	for i := range c.crons {
//		c.logger.Infof("%s job is scheduled to run at: %s, now: %s", c.crons[i].Name(), c.crons[i].GetJob().NextRun(), now)
//	}
//}
