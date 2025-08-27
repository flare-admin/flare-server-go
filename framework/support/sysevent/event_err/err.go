package event_err

import "github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"

var (
	AddEventFail        = herrors.NewServerError("AddEventFail")             //新增事件失败
	EditEventFail       = herrors.NewServerError("EditEventFail")            //修改事件失败
	DeleteEventFail     = herrors.NewServerError("DeleteEventFail")          //删除事件失败
	GetEventFail        = herrors.NewServerError("GetEventFail")             //获取事件是失败
	TopicIsExistFail    = herrors.NewBusinessServerError("TopicIsExistFail") //主题已经存在
	AddSubscribeFail    = herrors.NewServerError("AddSubscribeFail")         //新增订阅失败
	EditSubscribeFail   = herrors.NewServerError("EditSubscribeFail")        //修改订阅失败
	DeleteSubscribeFail = herrors.NewServerError("DeleteSubscribeFail")      //删除订阅失败
	GetSubscribeFail    = herrors.NewServerError("GetSubscribeFail")         //获取订阅是失败

	EventNotExistFail                    = herrors.NewBusinessServerError("EventNotExistFail")                    //事件不存在
	EventPublishFail                     = herrors.NewServerError("EventPublishFail")                             //事件发布失败
	EventNotEnable                       = herrors.NewBusinessServerError("EventNotEnable")                       //事件未开启
	SubscriptionIsAlreadyEnabled         = herrors.NewBusinessServerError("SubscriptionIsAlreadyEnabled")         //订阅已经开始
	SubscriptionEventFail                = herrors.NewServerError("SubscriptionEventFail")                        //订阅事件是比啊
	SubscriptionNoCorrespondingProcessor = herrors.NewBusinessServerError("SubscriptionNoCorrespondingProcessor") //订阅没有对应的处理器
	TheSameSubscriptionAlreadyExists     = herrors.NewBusinessServerError("TheSameSubscriptionAlreadyExists")     //已经存在相同订阅
)
