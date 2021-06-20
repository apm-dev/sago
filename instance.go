package sago

type SagaInstance struct {
	id                       string
	sagaType                 string
	lastRequestID            string
	serializedSagaData       []byte
	stateName                string
	destinationsAndResources map[string]string
	endState                 bool
	compensating             bool
}

func NewSagaInstance(sagaID, sagaType, stateName, lastReqID string, serializedData []byte, destAndRes map[string]string) *SagaInstance {
	return &SagaInstance{
		id:                       sagaID,
		sagaType:                 sagaType,
		lastRequestID:            lastReqID,
		serializedSagaData:       serializedData,
		stateName:                stateName,
		destinationsAndResources: destAndRes,
	}
}

func (si *SagaInstance) ID() string {
	return si.id
}
func (si *SagaInstance) SetID(id string) {
	si.id = id
}

func (si *SagaInstance) SagaType() string {
	return si.sagaType
}
func (si *SagaInstance) SetSagaType(st string) {
	si.sagaType = st
}

func (si *SagaInstance) LastRequestID() string {
	return si.lastRequestID
}
func (si *SagaInstance) SetLastRequestID(id string) {
	si.lastRequestID = id
}

func (si *SagaInstance) SerializedSagaData() []byte {
	return si.serializedSagaData
}
func (si *SagaInstance) SetSerializedSagaData(data []byte) {
	si.serializedSagaData = data
}

func (si *SagaInstance) StateName() string {
	return si.stateName
}
func (si *SagaInstance) SetStateName(sn string) {
	si.stateName = sn
}

func (si *SagaInstance) DestinationsAndResources() map[string]string {
	return si.destinationsAndResources
}
func (si *SagaInstance) SetDestinationsAndResources(dr map[string]string) {
	si.destinationsAndResources = dr
}

func (si *SagaInstance) IsEndState() bool {
	return si.endState
}
func (si *SagaInstance) SetEndState(e bool) {
	si.endState = e
}

func (si *SagaInstance) IsCompensating() bool {
	return si.compensating
}
func (si *SagaInstance) SetCompensating(c bool) {
	si.compensating = c
}
