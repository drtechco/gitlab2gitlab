package services

import (
	"context"
	"drtech.co/gl2gl/orm"
	"github.com/sirupsen/logrus"
	"sync"
)

type ErrorCodeService struct {
	logger      *logrus.Entry
	i18nService *I18nService
}

var (
	_errLock       sync.Mutex
	_errServiceMap = make(map[string]*ErrorCodeService)
	_errCodeMap    = make(map[int]string)
)

func GetErrorCodeService(lang string) *ErrorCodeService {
	_errLock.Lock()
	defer _errLock.Unlock()
	s, ok := _errServiceMap[lang]
	if ok {
		return s
	} else {
		errService := &ErrorCodeService{}
		errService.logger = logrus.WithField("Name", "ErrorCodeService")
		errService.i18nService = GetI18nService(lang)
		_errServiceMap[lang] = errService
		return errService
	}
}

func (s ErrorCodeService) Msg(errCode int) string {
	_errLock.Lock()
	defer _errLock.Unlock()
	i18nKey, ok := _errCodeMap[errCode]
	if ok {
		return s.i18nService.Get(i18nKey)
	}

	ecM := orm.DbQuery().ErrorCode
	err := ecM.WithContext(context.Background()).
		Where(ecM.Code.Eq(int32(errCode))).Select(ecM.I18nKey).Limit(1).Scan(&i18nKey)

	if err != nil {
		return "errCode find out err"
	}
	_errCodeMap[errCode] = i18nKey
	return s.i18nService.Get(i18nKey)
}
