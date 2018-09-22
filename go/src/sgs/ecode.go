package sgs

import "er"

const (
	_E_SGS_RUNNER = 0x1000

	_E_LOAD_CONF_FAIL = _E_SGS_RUNNER | er.IMPT_RECOVERABLE | er.ET_SERVICE | 0x1
)
