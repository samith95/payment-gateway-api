package error_constant

var (
	InvalidAuthIdField           = "authorisation id field is not valid"
	InvalidAmount                = "amount cannot be negative"
	InvalidCardExpiryDate        = "expiry date is not valid"
	InvalidCardNumber            = "card number is not valid"
	InvalidCvv                   = "cvv number is not valid"
	InvalidCurrencyCode          = "currency code is invalid"
	AuthorisationFailure         = "authorisation failure"
	CaptureFailure               = "capture failure"
	RefundFailure                = "refund failure"
	CancelledTransaction         = "transaction has been cancelled"
	TransactionRetrievalFailure  = "unable to retrieve authorisation transaction"
	RejectRetrievalFailure       = "unable to retrieve rejects"
	UpdateAvailableAmountFailure = "unable to update available amount"
	TransactionNotFound          = "authorisation transaction not found"
	ExpiredCard                  = "card is expired"
	RequestedAmountNotValid      = "the requested amount cannot be processed"
	TransactionStateInvalid      = "transaction is not in a state that allows this operation"
	OperationNameInvalid         = "passed operation name is invalid"
	UnableToCheckForInvalidState = "unable to check for invalid state"
	UnableToVoidTransaction      = "unable to void transaction"
)
