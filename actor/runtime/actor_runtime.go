package runtime

import (
	"encoding/json"
	"sync"

	"github.com/dapr/go-sdk/actor"
	"github.com/dapr/go-sdk/actor/api"
	"github.com/dapr/go-sdk/actor/config"
	actorErr "github.com/dapr/go-sdk/actor/error"
	"github.com/dapr/go-sdk/actor/manager"
)

type ActorRunTime struct {
	config        api.ActorRuntimeConfig
	actorManagers sync.Map
}

var actorRuntimeInstance *ActorRunTime

// NewActorRuntime creates an empty ActorRuntime.
func NewActorRuntime() *ActorRunTime {
	return &ActorRunTime{}
}

// GetActorRuntimeInstance gets or create runtime instance.
func GetActorRuntimeInstance() *ActorRunTime {
	if actorRuntimeInstance == nil {
		actorRuntimeInstance = NewActorRuntime()
	}
	return actorRuntimeInstance
}

// RegisterActorFactory registers the given actor factory from user, and create new actor manager if not exists.
func (r *ActorRunTime) RegisterActorFactory(f actor.Factory, opt ...config.Option) {
	conf := config.GetConfigFromOptions(opt...)
	actType := f().Type()
	r.config.RegisteredActorTypes = append(r.config.RegisteredActorTypes, actType)
	mng, ok := r.actorManagers.Load(actType)
	if !ok {
		newMng, err := manager.NewDefaultActorManager(conf.SerializerType)
		if err != actorErr.Success {
			return
		}
		newMng.RegisterActorImplFactory(f)
		r.actorManagers.Store(actType, newMng)
		return
	}
	mng.(manager.ActorManager).RegisterActorImplFactory(f)
}

func (r *ActorRunTime) GetJSONSerializedConfig() ([]byte, error) {
	data, err := json.Marshal(&r.config)
	return data, err
}

func (r *ActorRunTime) InvokeActorMethod(actorTypeName, actorID, actorMethod string, payload []byte) ([]byte, actorErr.ActorErr) {
	mng, ok := r.actorManagers.Load(actorTypeName)
	if !ok {
		return nil, actorErr.ErrActorTypeNotFound
	}
	return mng.(manager.ActorManager).InvokeMethod(actorID, actorMethod, payload)
}

func (r *ActorRunTime) Deactivate(actorTypeName, actorID string) actorErr.ActorErr {
	targetManager, ok := r.actorManagers.Load(actorTypeName)
	if !ok {
		return actorErr.ErrActorTypeNotFound
	}
	return targetManager.(manager.ActorManager).DetectiveActor(actorID)
}

func (r *ActorRunTime) InvokeReminder(actorTypeName, actorID, reminderName string, params []byte) actorErr.ActorErr {
	targetManager, ok := r.actorManagers.Load(actorTypeName)
	if !ok {
		return actorErr.ErrActorTypeNotFound
	}
	mng := targetManager.(manager.ActorManager)
	return mng.InvokeReminder(actorID, reminderName, params)
}

func (r *ActorRunTime) InvokeTimer(actorTypeName, actorID, timerName string, params []byte) actorErr.ActorErr {
	targetManager, ok := r.actorManagers.Load(actorTypeName)
	if !ok {
		return actorErr.ErrActorTypeNotFound
	}
	mng := targetManager.(manager.ActorManager)
	return mng.InvokeTimer(actorID, timerName, params)
}
