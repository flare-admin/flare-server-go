package base

import (
	"context"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/dtm-labs/rockscache"
	"github.com/flare-admin/flare-server-go/framework/pkg/actx"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/mqevent"
	"github.com/flare-admin/flare-server-go/framework/pkg/mqevent/manager"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/biz"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/model"
	"time"
)

const expirationTime = time.Second * 60 * 60

type SubscribeManagerUseCase struct {
	sr      biz.ISubscribeRepo
	par     biz.ISubscribeParameterRepo
	deadPar biz.IDeadLetterSubscribeRepo
	rc      *rockscache.Client
}

func NewSubscribeManagerUseCase(sr biz.ISubscribeRepo, par biz.ISubscribeParameterRepo, deadPar biz.IDeadLetterSubscribeRepo, rc *rockscache.Client) manager.ISubscribeSmServerApi {
	return &SubscribeManagerUseCase{
		sr:      sr,
		rc:      rc,
		par:     par,
		deadPar: deadPar,
	}
}
func (s SubscribeManagerUseCase) GetByStatus(ctx context.Context, status int32) ([]*manager.Subscribe, error) {
	subscribes, err := s.sr.GetByStatus(ctx, int8(status))
	if err != nil {
		if database.IfErrorNotFound(err) {
			return make([]*manager.Subscribe, 0), nil
		}
		hlog.CtxErrorf(ctx, "get subscribe list failed: %v", err)
		return nil, err
	}
	evs := make([]*manager.Subscribe, 0, len(subscribes))
	for _, subscribe := range subscribes {
		evs = append(evs, &manager.Subscribe{
			Id:     subscribe.Id,
			Topic:  subscribe.Topic,
			Group:  subscribe.Group,
			Name:   subscribe.Name,
			Status: subscribe.Status,
		})
	}
	return evs, nil
}

func (s SubscribeManagerUseCase) GetParameters(ctx context.Context, topic, channel string) (map[string]interface{}, error) {
	fetch, err := s.rc.Fetch(model.GetSubordinatesParameterCacheKey(topic, channel, actx.GetTenantId(ctx)), expirationTime, func() (string, error) {
		data, err := s.GetSubordinatesPar(ctx, topic, channel)
		if err != nil {
			return "", err
		}
		marshal, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		return string(marshal), nil
	})
	if err != nil {
		return nil, err
	}
	newInfo := make(map[string]interface{})
	err = json.Unmarshal([]byte(fetch), &newInfo)
	if err != nil {
		return nil, err
	}
	return newInfo, nil
}

func (s SubscribeManagerUseCase) DadEventSave(ctx context.Context, event *mqevent.DeadLetterEvent) error {
	subscribe := model.DeadLetterSubscribe{}
	err := subscribe.FromDeadLetterEvent(event)
	if err != nil {
		hlog.CtxErrorf(ctx, "convert dead letter event to subscribe failed: %v", err)
		return err
	}
	_, err = s.deadPar.Add(ctx, &subscribe)
	return err
}

func (s SubscribeManagerUseCase) GetSubordinatesParameterCache(ctx context.Context, topic, group string) (map[string]interface{}, error) {
	fetch, err := s.rc.Fetch(model.GetSubordinatesParameterCacheKey(topic, group, actx.GetTenantId(ctx)), expirationTime, func() (string, error) {
		data, err := s.GetSubordinatesPar(ctx, topic, group)
		if err != nil {
			return "", err
		}
		marshal, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		return string(marshal), nil
	})
	if err != nil {
		return nil, err
	}
	newInfo := make(map[string]interface{})
	err = json.Unmarshal([]byte(fetch), &newInfo)
	if err != nil {
		return nil, err
	}
	return newInfo, nil
}

func (s SubscribeManagerUseCase) GetSubordinatesPar(ctx context.Context, topic, group string) (map[string]interface{}, error) {
	subscribe, err := s.sr.GetByTopicAndGroup(ctx, topic, group)
	if err != nil {
		hlog.CtxErrorf(ctx, "get subscribe list failed: %v", err)
		return nil, err
	}
	if subscribe == nil || subscribe.Id == "" {
		return make(map[string]interface{}), nil
	}
	parameters, err := s.par.GetBySubscribe(ctx, subscribe.Id)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return make(map[string]interface{}), nil
		}
		hlog.CtxErrorf(ctx, "get subscribe list failed: %v", err)
		return nil, err
	}
	if parameters == nil || len(parameters) == 0 {
		return make(map[string]interface{}), nil
	}
	parametersToMap, err := model.ConvertParametersToMap(parameters)
	if err != nil {
		hlog.CtxErrorf(ctx, "convert parameters to map failed: %v", err)
		return nil, err
	}
	return parametersToMap, nil
}
