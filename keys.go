package sago

const (
	ZB_SAGO_KEY  = "sago"
	ZB_SAGA_ID   = "saga_id"
	ZB_SAGA_TYPE = "saga_type"
)

func buildSagoKey(sagaType, sagaId string) string {
	return sagaType + ":" + sagaId
}

func removeReservedVariableKeys(vars map[string]interface{}) {
	delete(vars, ZB_SAGO_KEY)
	delete(vars, ZB_SAGA_ID)
	delete(vars, ZB_SAGA_TYPE)
}
