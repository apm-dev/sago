package sago

const (
	ZB_SAGO_KEY    = "sago"
	ZB_SAGA_ID     = "saga_id"
	ZB_SAGA_TYPE   = "saga_type"
	ZB_TASK_RESULT = "step_result"
)

func buildSagoKey(sagaType, sagaId string) string {
	return sagaType + ":" + sagaId
}
