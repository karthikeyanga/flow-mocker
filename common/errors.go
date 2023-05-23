package common

import (
	"fmt"
)

//TRANSACTION ERROR : START
type InvalidTxnStatusError struct {
	Msg        string
	DisplayMsg string
}

func (e *InvalidTxnStatusError) Error() string {
	return "Ongoing transaction status" + e.Msg
}

func (e *InvalidTxnStatusError) DisplayError() string {
	if e.DisplayMsg == "" {
		return "Invalid transaction status"
	}
	return e.DisplayMsg
}

type DuplicateTransactionRequest struct {
	Msg        string
	DisplayMsg string
}

func (e *DuplicateTransactionRequest) Error() string {
	return "Duplicate Transaction Request" + e.Msg
}

func (e *DuplicateTransactionRequest) DisplayError() string {
	if e.DisplayMsg == "" {
		return "Duplicate transaction request"
	}
	return e.DisplayMsg
}

type InvalidRefundRequest struct {
	Msg        string
	DisplayMsg string
}

func (irr *InvalidRefundRequest) Error() string {
	return "Invalid Refund Request" + irr.Msg
}

func (irr *InvalidRefundRequest) DisplayError() string {
	if irr.DisplayMsg == "" {
		return "Invalid refund request"
	}
	return irr.DisplayMsg
}

type InEligibleTransactionAttempt struct {
	Msg        string
	DisplayMsg string
}

func (iet *InEligibleTransactionAttempt) Error() string {
	return "In eligible Transaction Attempt " + iet.Msg
}

func (iet *InEligibleTransactionAttempt) DisplayError() string {
	return iet.DisplayMsg
}

type InvalidTransactionRequest struct {
	Msg        string
	DisplayMsg string
}

func (itr *InvalidTransactionRequest) Error() string {
	return "Invalid Transaction Request" + itr.Msg
}

func (itr *InvalidTransactionRequest) DisplayError() string {
	return itr.DisplayMsg
}

// TRANSACTION ERROR : END
type ErrInvalidParam struct {
	Msg string
}

func (e *ErrInvalidParam) Error() string {
	return "invalid param" + e.Msg
}

type ErrMultiErrors struct {
	Msg    string
	Errors map[string]error
}

func NewMultiError(msg string) ErrMultiErrors {
	return ErrMultiErrors{
		Msg:    msg,
		Errors: map[string]error{},
	}
}
func (e *ErrMultiErrors) Error() string {
	return fmt.Sprintf(e.Msg)
}

func (e ErrMultiErrors) AddError(key string, err error) {
	if e.Errors == nil {
		e.Errors = map[string]error{}
	}
	e.Errors[key] = err
}
func (e *ErrMultiErrors) Err() error {
	if len(e.Errors) == 0 {
		return nil
	}
	return e
}

type ErrRedirectToUrl struct {
	Url string
}

func (e *ErrRedirectToUrl) Error() string {
	return fmt.Sprintf("Redirect to %s", e.Url)
}

type ErrBadRequest struct {
	Msg string
}

func (e *ErrBadRequest) Error() string {
	return fmt.Sprintf("Bad Request - %s", e.Msg)
}

type ErrUnAuthorised struct {
	Msg string
}

func (e *ErrUnAuthorised) Error() string {
	return fmt.Sprintf("UnAuthorised - %s", e.Msg)
}
