package services

import (
	"context"
	"drtech.co/gl2gl/orm"
	"github.com/sirupsen/logrus"

	"sync"
)

type I18nService struct {
	logger    *logrus.Entry
	lang      string
	langValue map[string]string
	lock      sync.Mutex
}

var (
	_i18nLock       sync.Mutex
	_i18nServiceMap = make(map[string]*I18nService)
)

func GetI18nService(lang string) *I18nService {
	_i18nLock.Lock()
	defer _i18nLock.Unlock()
	s, ok := _i18nServiceMap[lang]
	if ok {
		return s
	} else {
		_i18nService := &I18nService{}
		_i18nService.logger = logrus.WithField("Name", "I18nService")
		_i18nService.langValue = make(map[string]string)
		_i18nService.lock = sync.Mutex{}
		_i18nService.lang = lang
		_i18nServiceMap[lang] = _i18nService
		return _i18nService
	}

}

func (s *I18nService) GetLangList() []string {
	langM := orm.DbQuery().Lang
	var langKeys []string
	err := langM.WithContext(context.Background()).Select(langM.Key).Scan(&langKeys)
	if err != nil {
		s.logger.Error(err)
	}
	return langKeys
}

func (s *I18nService) Get(key string) string {
	_i18nLock.Lock()
	defer _i18nLock.Unlock()
	i18nValue, ok := s.langValue[key]
	if ok {
		return i18nValue
	}
	var v string
	i18nM := orm.DbQuery().I18n
	err := i18nM.WithContext(context.Background()).
		Where(i18nM.Key.Eq(key)).Where(i18nM.LangKey.Eq(s.lang)).Select(i18nM.Value).Limit(1).Row().Scan(&v)
	if err != nil && v != "" {
		s.logger.Error(err)
		return "NOT FOUND"
	}
	s.langValue[key] = v
	return v
}
