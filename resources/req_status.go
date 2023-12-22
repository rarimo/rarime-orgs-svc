package resources

import "encoding/json"

type RequestStatus int32

const (
	RequestStatus_Created   RequestStatus = 0
	RequestStatus_Accepted  RequestStatus = 1
	RequestStatus_Filled    RequestStatus = 2
	RequestStatus_Approved  RequestStatus = 3
	RequestStatus_Rejected  RequestStatus = 4
	RequestStatus_Submitted RequestStatus = 5

	requestStatus_Created_Str   = "created"
	requestStatus_Accepted_Str  = "accepted"
	requestStatus_Filled_Str    = "filled"
	requestStatus_Approved_Str  = "approved"
	requestStatus_Rejected_Str  = "rejected"
	requestStatus_Submitted_Str = "submitted"
)

var requestStatusIntStr = map[RequestStatus]string{
	RequestStatus_Created:   requestStatus_Created_Str,
	RequestStatus_Accepted:  requestStatus_Accepted_Str,
	RequestStatus_Filled:    requestStatus_Filled_Str,
	RequestStatus_Approved:  requestStatus_Approved_Str,
	RequestStatus_Rejected:  requestStatus_Rejected_Str,
	RequestStatus_Submitted: requestStatus_Submitted_Str,
}

var requestStatusStrInt = map[string]RequestStatus{
	requestStatus_Created_Str:   RequestStatus_Created,
	requestStatus_Accepted_Str:  RequestStatus_Accepted,
	requestStatus_Filled_Str:    RequestStatus_Filled,
	requestStatus_Approved_Str:  RequestStatus_Approved,
	requestStatus_Rejected_Str:  RequestStatus_Rejected,
	requestStatus_Submitted_Str: RequestStatus_Submitted,
}

func RequestStatusFromString(s string) RequestStatus {
	return requestStatusStrInt[s]
}

func (s RequestStatus) String() string {
	return requestStatusIntStr[s]
}

func (s RequestStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(Flag{
		Name:  requestStatusIntStr[s],
		Value: int32(s),
	})
}

func (s *RequestStatus) UnmarshalJSON(b []byte) error {
	var res Flag
	err := json.Unmarshal(b, &res)
	if err != nil {
		return err
	}

	*s = RequestStatus(res.Value)
	return nil
}

func (s RequestStatus) Int() int {
	return int(s)
}
func (s RequestStatus) Int16() int16 {
	return int16(s)
}
