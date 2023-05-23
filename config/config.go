package config

import (
	"mocker/common"
	"mocker/core"
	"mocker/model"
	"sync"

	"fmt"
	"io"
	"io/ioutil"
	"mocker/util/database"
	"mocker/util/log"
	"os"
	"path/filepath"
	"time"

	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v2"
)

//AppContext holds config AppContext
type AppConfig struct {
	Config        *AppStartupConfigType
	db            *gorm.DB
	Apis          map[int64]common.IApi
	loadedFlows   []common.IFlow
	loadedObjects map[int64]common.IObject
	loadedObjLock sync.RWMutex
	//PdtProducer            *kafka.ProducerWrapper
	//MetricProducer         *kafka.ProducerWrapper
	loadStaticDataChan chan int
	LocalCache         map[LocalCacheKey]interface{}
	files              map[string]*os.File
	log                *log.Log
}

func (ac *AppConfig) Init(configFileReader io.Reader) error {
	ac.files = map[string]*os.File{}
	//Config Init
	file_content, err := ioutil.ReadAll(configFileReader)
	if err != nil {
		return err
	}
	ac.Config = &AppStartupConfigType{}
	err = yaml.Unmarshal(file_content, &ac.Config)
	if err != nil {
		return err
	}
	logLevel := log.LOG_LEVEL_MAP[ac.Config.LoggingConfig.Level]
	//Init Log
	if err := ac.openFile("log", ac.Config.LoggingConfig.Path); err != nil {
		return err
	}
	file := ac.files["log"]
	ac.log = log.New(file, int32(logLevel))
	logger := ac.log.WithFields(map[string]interface{}{"section": "Config Init"})
	if ac.Config.ServerConfig.StaticPath == "" {
		staticPath, _ := os.Executable()
		staticPath = filepath.Dir(staticPath)
		ac.Config.ServerConfig.StaticPath = staticPath
	}
	ac.loadStaticDataChan = make(chan int)
	//DB init
	DBConfig := ac.Config.DatabaseConfig
	if err := ac.openFile("db", DBConfig.LogFile); err != nil {
		return err
	}
	file = ac.files["db"]
	ac.files["db"] = file
	if DBConfig.Dialect == "sqlite" {
		ac.db, err = database.SqlLiteInit(DBConfig.Path)
	} else {
		ac.db, err = database.DatabaseInit(DBConfig.UserName, DBConfig.Password,
			DBConfig.Protocol, DBConfig.Host, DBConfig.Port, DBConfig.Database,
			DBConfig.Timeout, DBConfig.MaxIdleConnectionCount, DBConfig.MaxConnectionCount, logger, file, ac.Config.APIMode)
	}
	if err != nil {
		return err
	}
	apis, err := model.GetAllEnabledApis(ac.db)
	if err != nil {
		return err
	}
	ac.Apis = map[int64]common.IApi{}
	for _, apiM := range apis {
		api, err := core.GetApiFromModel(apiM)
		if err != nil {
			return err
		}
		ac.Apis[apiM.ID] = api
	}
	// Open other log files ;
	if err := ac.openFile("server.log.access", ac.Config.ServerConfig.AccessLog); err != nil {
		return err
	}
	if err := ac.openFile("server.log.error", ac.Config.ServerConfig.ErrorLog); err != nil {
		return err
	}
	return nil
}
func (ac *AppConfig) GetFile(key string) (*os.File, bool) {
	f, err := ac.files[key]
	return f, err
}

func (ac *AppConfig) openFile(key, filename string) error {
	if file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666); err != nil {
		return err
	} else {
		ac.files[key] = file
	}
	return nil
}

func (ac *AppConfig) Close() {
	//Close all files and connections
	for n, f := range ac.files {
		if err := f.Close(); err != nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("Error while closing file %s - %s", n, err))
		}
	}
	//Close db
	if err := ac.db.Close(); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("Error while closing db - %s", err))
	}
	//Close Cacher

	//Close Kafka producers and Consumers
	//ac.MetricProducer.Close()
	//ac.NavConsumer.Close()
	//ac.PdtProducer.Close()
	//Close Other things
}

func (ac *AppConfig) NewContext(requestId string) *AppContext {
	ctx := AppContext{
		Context:   nil,
		AppConfig: ac,
		Log:       ac.log.WithFields(map[string]interface{}{"request": requestId}),
		Now:       time.Now(),
		requestId: requestId,
		data:      map[string]interface{}{},
	}
	return &ctx
}

func (ac *AppConfig) GetFlows() []common.IFlow {
	return ac.loadedFlows
}
func (ac *AppConfig) LoadFlow(flow common.IFlow) {
	ac.UnloadFlow(flow)
	ac.loadedFlows = append(ac.loadedFlows, flow)
}

func (ac *AppConfig) UnloadFlow(rflow common.IFlow) {
	id := rflow.Id()
	for i, flow := range ac.loadedFlows {
		if flow.Id() == id {
			l := len(ac.loadedFlows) - 1
			ac.loadedFlows[i] = ac.loadedFlows[l]
			ac.loadedFlows[l] = nil
			ac.loadedFlows = ac.loadedFlows[:l]
			return
		}
	}
}
func (ac *AppConfig) EvictOldFlows() {
	ttl, ok := ac.Config.Constants["flowEvictionHours"]
	if !ok {
		ttl = "24"
	}
	ttlDur := time.Duration(common.StringToInt(ttl)) * time.Hour
	count := 0
	for _, flow := range ac.loadedFlows {
		if time.Since(flow.LastAccessedAt()) <= ttlDur {
			count++
		}
	}
	if count == len(ac.loadedFlows) {
		return
	}
	res := make([]common.IFlow, count)
	for _, flow := range ac.loadedFlows {
		if time.Since(flow.LastAccessedAt()) <= ttlDur {
			res = append(res, flow)
		}
	}
	ac.loadedFlows = res
}

func (ac *AppConfig) GetObject(id int64) (common.IObject, bool) {
	ac.loadedObjLock.RLock()
	defer ac.loadedObjLock.RUnlock()
	obj, ok := ac.loadedObjects[id]
	return obj, ok
}
func (ac *AppConfig) GetObjects() []common.IObject {
	ac.loadedObjLock.RLock()
	defer ac.loadedObjLock.RUnlock()
	res := []common.IObject{}
	for _, obj := range ac.loadedObjects {
		res = append(res, obj)
	}
	return res
}
func (ac *AppConfig) SetObject(id int64, obj common.IObject) {
	if _, ok := ac.loadedObjects[id]; ok {
		return
	}
	ac.loadedObjLock.Lock()
	defer ac.loadedObjLock.Unlock()
	ac.loadedObjects[id] = obj
}

func (ac *AppConfig) UnloadObject(id int64) {
	if _, ok := ac.loadedObjects[id]; !ok {
		return
	}
	ac.loadedObjLock.Lock()
	defer ac.loadedObjLock.Unlock()
	delete(ac.loadedObjects, id)
}

//LocalCacheKey is used to put and get data from LocalCache
type LocalCacheKey string

func New(configFileReader io.Reader) (*AppConfig, error) {

	appConfig := AppConfig{
		loadedFlows:   []common.IFlow{},
		loadedObjects: map[int64]common.IObject{},
	}
	err := appConfig.Init(configFileReader)
	return &appConfig, err
}
