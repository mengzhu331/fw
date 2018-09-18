package web

//Conf sgs web server configuration
type Conf struct {
	Port        int
	EPM         EndPointMap
	WSReadBuff  int
	WSWriteBuff int
	CPS         int
}

func (me *Conf) epmv(key int) string {
	if me == nil || me.EPM == nil || me.EPM[key] == "" {
		return demp()[key]
	}
	return me.EPM[key]
}

func (me *Conf) cps() int {
	if me == nil || me.CPS == 0 {
		return _defaultConf.CPS
	}
	return me.CPS
}

func (me *Conf) port() int {
	if me == nil || me.Port == 0 {
		return _defaultConf.Port
	}
	return me.Port
}

func (me *Conf) wsReadBuff() int {
	if me == nil || me.WSReadBuff == 0 {
		return _defaultConf.WSReadBuff
	}
	return me.WSReadBuff
}

func (me *Conf) wsWriteBuff() int {
	if me == nil || me.WSWriteBuff == 0 {
		return _defaultConf.WSWriteBuff
	}
	return me.WSWriteBuff
}

var _defaultConf = Conf{
	Port:        9090,
	EPM:         nil,
	WSReadBuff:  1024,
	WSWriteBuff: 1024,
	CPS:         3,
}
