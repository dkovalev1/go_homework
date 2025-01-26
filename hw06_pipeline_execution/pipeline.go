package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func drainChannel(in In) {
	for range in { //nolint
	}
}

func checkDone(done In, in In) bool {
	select {
	case <-done:
		drainChannel(in)
		return true
	default:
	}
	return false
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Place your code here.

	inCur := in
	var outCur Out

	for _, stage := range stages {
		myOut := make(Bi)
		go func(in In, done In) {
			defer close(myOut)

			for {
				select {
				case <-done:
					drainChannel(in)
					return
				case val, ok := <-in:
					if ok {
						select {
						case myOut <- val:
						case <-done:
							return
						}
					} else {
						// in closed, we done here
						return
					}
				}
				shouldReturn := checkDone(done, in)
				if shouldReturn {
					return
				}
			}
		}(inCur, done)

		outCur = stage(myOut)
		inCur = outCur
	}

	return outCur
}
