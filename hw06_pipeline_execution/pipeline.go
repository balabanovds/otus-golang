package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

type IPipeline interface {
	Pipe(stage Stage) IPipeline
	Exec() Out
}

type pipeline struct {
	stages     []Stage
	dataStream In
	doneStream In
}

func NewPipeline(inCh, doneCh In, stages ...Stage) IPipeline {
	return &pipeline{
		stages:     stages,
		dataStream: inCh,
		doneStream: doneCh,
	}
}

// Pipe add new stage to pipeline.
// In our example it is not so necessary.
// But added to satisfy IPipeline interface.
func (p *pipeline) Pipe(stage Stage) IPipeline {
	p.stages = append(p.stages, stage)
	return p
}

// Exec run stages one by one.
func (p *pipeline) Exec() Out {
	for _, stage := range p.stages {
		stageStream := make(Bi)

		go func(stageStream Bi, inStream In) {
			defer close(stageStream)

			for {
				select {
				case <-p.doneStream:
					return
				case data, ok := <-inStream:
					if !ok {
						return
					}
					if data != nil {
						select {
						case <-p.doneStream:
							return
						case stageStream <- data:
						}
					}
				}
			}
		}(stageStream, p.dataStream)

		p.dataStream = stage(stageStream)
	}

	return p.dataStream
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	return NewPipeline(in, done, stages...).Exec()
}
