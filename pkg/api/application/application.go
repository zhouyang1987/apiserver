package application

import (
	"net/http"

	"apiserver/pkg/storage/mysqld"
	"apiserver/pkg/util/jsonx"
	"apiserver/pkg/util/log"

	"github.com/emicklei/go-restful"
)

func Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/app").
		Doc("manage application").
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.POST("").To(new(App).Insert).
		Doc("create application").
		Produces("CreateApplication").
		Reads(App{}))

	ws.Route(ws.GET("/{app-id}").To(new(App).QueryOne).
		// docs
		Doc("get a app").
		Operation("findUser").
		Param(ws.PathParameter("app-id", "identifier of the app").DataType("int")).
		Writes(App{})) // on the response

	container.Add(ws)
}

type AppStatus int32
type UpdateStatus int32

const (
	AppBuilding  AppStatus = 0
	AppSuccessed AppStatus = 1
	AppFailed    AppStatus = 2
	AppRunning   AppStatus = 3
	AppStop      AppStatus = 4

	StartFailed    UpdateStatus = 10
	StartSuccessed UpdateStatus = 11

	StopFailed    UpdateStatus = 20
	StopSuccessed UpdateStatus = 21

	ScaleFailed    UpdateStatus = 30
	ScaleSuccessed UpdateStatus = 31

	UpdateConfigFailed    UpdateStatus = 40
	UpdateConfigSuccessed UpdateStatus = 41

	RedeploymentFailed    UpdateStatus = 50
	RedeploymentSuccessed UpdateStatus = 51
)

//App is struct of application
type App struct {
	Id            int       `json:"id" xorm:"pk not null autoincr int(11)"`
	Name          string    `json:"name" xorm:"varchar(256)"`
	Region        string    `json:"region" xorm:"varchar(256)"`
	Memory        string    `json:"memory" xorm:"varchar(11)"`
	Cpu           string    `json:"cpu" xorm:"varchar(11)"`
	InstanceCount int32     `json:"instanceCount" xorm:"int(11)"`
	Envs          string    `json:"envs" xorm:"varchar(1024)"`
	Ports         string    `json:"ports" xorm:"varchar(1024)"`
	Image         string    `json:"image" xorm:""`
	Status        AppStatus `json:"status" xorm:"int(1)"` //构建中 0 成功 1 失败 2 运行中 3 停止 4
	UserName      string    `json:"userName" xorm:"varchar(256)"`
	Remark        string    `json:"remark" xorm:"varchar(1024)"`
}

var (
	engine = mysqld.GetEngine()
)

func init() {
	engine.ShowSQL(true)
	if err := engine.Sync(new(App)); err != nil {
		log.Fatalf("Sync fail :%s", err.Error())
	}
}

func (app *App) String() string {
	appStr, err := jsonx.ToJson(app)
	if err != nil {
		log.Errorf("node to string err :%s", err.Error())
		return ""
	}
	return appStr
}

func (app *App) Insert(request *restful.Request, response *restful.Response) {
	aps := new(App)
	err := request.ReadEntity(aps)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	_, err = engine.Insert(aps)
	if err != nil {
		return
	}
	response.WriteHeaderAndEntity(http.StatusCreated, aps)
}

func (app *App) Delete() error {
	_, err := engine.Id(app.Id).Delete(app)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) Update() error {
	_, err := engine.Id(app.Id).Update(app)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) QueryOne(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("app-id")

	log.Debug(id)

	aps := &App{}
	// engine.Id(id).Get(&aps)
	_, err := engine.Id(1).Get(aps)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, err.Error())
		return
	}
	// if !has {
	// 	response.AddHeader("Content-Type", "text/plain")
	// 	response.WriteErrorString(http.StatusNotFound, "404: User could not be found.")
	// 	return
	// }

	log.Debugf("%#v", aps)
	response.WriteEntity(aps)
}

func (app *App) QuerySet() ([]*App, error) {
	appSet := []*App{}
	err := engine.Where("1 and 1 order by id desc").Find(&appSet)
	if err != nil {
		return nil, err
	}
	return appSet, nil
}
