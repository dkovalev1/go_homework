package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func drainChannel(in In, stageNo int) {
	/* For the first stage (0) we read from the input which is done in goroutine
	** Close shall be done synchronously to avoid race condition
	** For other stages we can read it asynchronously to allow writer stages to complete
	** with 100ms duration  and to fit into 50ms time limit.
	** to
	 */
	if stageNo == 0 {
		for range in { //nolint
		}
	} else {
		go func() {
			for range in { //nolint
			}
		}()
	}
}

func checkDone(done In, in In, stageNo int) bool {
	select {
	case <-done:
		drainChannel(in, stageNo)
		return true
	default:
	}
	return false
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	inCur := in
	var outCur Out

	for i, stage := range stages {
		myOut := make(Bi)
		go func(in In, done In) {
			defer close(myOut)

			for {
				select {
				case <-done:
					drainChannel(in, i)
					return
				case val, ok := <-in:
					if ok {
						/* read ok, write it to stage and check for while writing */
						select {
						case myOut <- val:
						case <-done:
							drainChannel(in, i)
							return
						}
					} else {
						// channel in closed, we are done
						return
					}
				}
				shouldReturn := checkDone(done, in, i)
				if shouldReturn {
					drainChannel(in, i)
					return
				}
			}
		}(inCur, done)

		outCur = stage(myOut)
		inCur = outCur
	}

	return outCur
}
