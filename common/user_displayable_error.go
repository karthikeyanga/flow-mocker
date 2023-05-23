package common

type ErrUserDisplayable struct {
	displayErrMsg      string
	displaySolutionMsg string
	errMsg             string
	err                error
	errCode            string
	showInPage         bool
}

func (this ErrUserDisplayable) Error() string {
	if this.errMsg != "" {
		return this.errMsg
	}
	if this.err != nil {
		return this.err.Error()
	}
	if this.displayErrMsg != "" {
		return this.displayErrMsg
	}
	return "Some Error Occurred"
}

func NewErrUserDisplayable(displayMsg, errMsg, solution, errCode string) *ErrUserDisplayable {
	return &ErrUserDisplayable{
		displayErrMsg:      displayMsg,
		errMsg:             errMsg,
		displaySolutionMsg: solution,
		errCode:            errCode,
	}
}

func NewErrUserDisplayableFromError(displayMsg string, err error, solution, errCode string) *ErrUserDisplayable {
	return &ErrUserDisplayable{
		displayErrMsg:      displayMsg,
		err:                err,
		displaySolutionMsg: solution,
		errCode:            errCode,
	}
}

func (this ErrUserDisplayable) Problem() string {
	return this.displayErrMsg
}

func (this ErrUserDisplayable) ErrCode() string {
	return this.errCode
}
func (this ErrUserDisplayable) Solution() string {
	return this.displaySolutionMsg
}
func (this *ErrUserDisplayable) ShowInPage() {
	this.showInPage = true
}

func (this ErrUserDisplayable) HasTobeShowedInPage() bool {
	return this.showInPage
}
