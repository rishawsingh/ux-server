package crons

//type ScoreCron struct {
//	logger         *internal.Logger
//	productService services.Product
//	inMemScore     *utils.InMemScore
//	job            *gocron.Job
//}
//
//const runEveryThreeHour = 3
//
//func NewScoreCron(logger *internal.Logger, productService services.Product, inMemStore *utils.InMemScore) *ScoreCron {
//	return &ScoreCron{
//		logger:         logger,
//		productService: productService,
//		inMemScore:     inMemStore,
//	}
//}
//
//func (s *ScoreCron) Name() string {
//	return "scoring-job"
//}
//
//func (s *ScoreCron) RunAt() null.Time {
//	return null.Time{}
//}
//
//func (s *ScoreCron) IntervalHour() null.Int {
//	return null.IntFrom(runEveryThreeHour)
//}
//
//func (s *ScoreCron) SetJob(job *gocron.Job) {
//	s.job = job
//	s.job.SingletonMode()
//}
//
//func (s *ScoreCron) GetJob() *gocron.Job {
//	return s.job
//}
//
//func (s *ScoreCron) Run() error {
//	s.logger.Infof("Starting %s", s.Name())
//	start := time.Now()
//	defer func(start time.Time, s *ScoreCron) {
//		s.logger.Infof("%s took %d millisecond to complete", s.Name(), time.Since(start).Milliseconds())
//	}(start, s)
//	productWithAttributes, scoreErr := s.productService.GetAllProductWithAttributes(context.TODO(), nil)
//	if scoreErr != nil {
//		s.logger.Errorf("%s failed with error %+v", s.Name(), scoreErr)
//		return scoreErr
//	}
//	s.inMemScore.Set(productWithAttributes)
//	return nil
//}
