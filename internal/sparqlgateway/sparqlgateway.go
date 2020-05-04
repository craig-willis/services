package sparqlgateway

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	restful "github.com/emicklei/go-restful"
)

type URLSet []string

// New fires up the services to query the graph
func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/api/v1/graph").
		Doc("Graph query services").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) //Consumes(restful.M).

	service.Route(service.POST("/ressetdetails").To(ResourceSetCall).
		Doc("Call for details on an array of resources from the triplestore (graph)").
		Param(service.BodyParameter("body", "The body containing an array of URIs to obtain parameter values from")).
		Writes([]ResourceResults{}).
		Consumes("application/x-www-form-urlencoded").
		Operation("ResourceSetCall"))

	service.Route(service.POST("/ressetpeople").To(ResourceSetPeopleCall).
		Doc("Call for people associated with an array of resources from the triplestore (graph)").
		Param(service.BodyParameter("body", "The body containing an array of URIs to obtain people relation values from")).
		Writes([]ResourceSetPeople{}).
		Consumes("application/x-www-form-urlencoded").
		Operation("ResourceSetPeopleCall"))

	service.Route(service.GET("/details").To(Details).
		Doc("Call for details on a resource from the triplestore (graph)").
		Param(service.QueryParameter("r", "Resource ID").DataType("string")).
		Writes([]LogoResults{}).
		Operation("Details"))

	// add in start point and length cursors
	service.Route(service.GET("/resdetails").To(ResourceCall).
		Doc("Call for measurement and prameter details").
		Param(service.QueryParameter("r", "Resource ID").DataType("string")).
		Writes([]ResourceResults{}).
		Operation("ResourceCall"))

	return service
}

// Dev fires up the DEVELOPMENT services
func Dev() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/api/dev/graph").
		Doc("DEV: Graph query services").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON) //Consumes(restful.M).

	service.Route(service.GET("/temporal").To(Temporal).
		Doc("Dev call for temporal data in the graph ").
		Param(service.QueryParameter("b", "Resource ID").DataType("string")).
		Param(service.QueryParameter("e", "Resource ID").DataType("string")).
		Writes([]LogoResults{}).
		Operation("Temporal"))

	service.Route(service.GET("/logo").To(Logo).
		Doc("Call for logo URL on a resource from the triplestore (graph)").
		Param(service.QueryParameter("r", "Resource ID").DataType("string")).
		Writes([]LogoResults{}).
		Operation("Logo"))

	service.Route(service.GET("/describe").To(Describe).
		Doc("Describe a triple store resource").
		Param(service.QueryParameter("r", "Resource ID").DataType("string")).
		Writes([]LogoResults{}).
		Operation("Describe"))

	service.Route(service.GET("/orgsearch").To(Organizations).
		Doc("Search an organization based on its description").
		Param(service.QueryParameter("r", "Resource ID").DataType("string")).
		Writes([]LogoResults{}).
		Operation("Organizations"))

	return service
}

func Organizations(request *restful.Request, response *restful.Response) {
	resource := request.QueryParameter("r")

	sr := OrgCall(resource, getHost(request))
	// fmt.Println(sr)
	// response.WriteJson(string(sr), " ")
	response.AddHeader("Content-Type", "application/json")
	response.Write(sr)
}

func Describe(request *restful.Request, response *restful.Response) {
	resource := request.QueryParameter("r")

	sr := DescribeCall(resource, getHost(request))
	response.WriteEntity(sr)
}

// Details call for details on the resource from the graph
func Details(request *restful.Request, response *restful.Response) {
	resource := request.QueryParameter("r")

	sr := DetailsCall(resource, getHost(request))
	response.WriteEntity(sr)
}

// ResourceCall call for details on the resource from the graph
func ResourceCall(request *restful.Request, response *restful.Response) {
	resource := request.QueryParameter("r")

	sr := ResCall(resource, getHost(request))
	response.WriteEntity(sr)
}

// Logo call for details on the resource from the graph
func Logo(request *restful.Request, response *restful.Response) {
	resource := request.QueryParameter("r")

	sr := LogoCall(resource, getHost(request))
	response.WriteEntity(sr)
}

// ResourceSetCall call for details on the resource array from the graph
func ResourceSetCall(request *restful.Request, response *restful.Response) {

	body, err := request.BodyParameter("body")
	if err != nil {
		log.Printf("Error on body parameter read %v \n", err)
		log.Println(err)
		response.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	ja := []byte(body)
	var jas URLSet
	err = json.Unmarshal(ja, &jas)
	if err != nil {
		log.Println("error with unmarshal..   return http error")
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	sr := ResSetCall(jas, getHost(request))
	response.WriteEntity(sr)
}

// ResourceSetPeopleCall call for details on the resource array from the graph
func ResourceSetPeopleCall(request *restful.Request, response *restful.Response) {

	body, err := request.BodyParameter("body")
	if err != nil {
		log.Printf("Error on body parameter read %v \n", err)
		log.Println(err)
		response.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	ja := []byte(body)
	var jas URLSet
	err = json.Unmarshal(ja, &jas)
	if err != nil {
		log.Println("error with unmarshal..   return http error")
		log.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	sr := ResSetPeople(jas, getHost(request))
	response.WriteEntity(sr)
}


func getHost(req *restful.Request) (string) {
	var host string
	hostvalues, ok := req.Request.Header["Host"]
	if !ok || len(hostvalues) == 0 {
		host = req.Request.Host
	} else	{
		host = hostvalues[0]
	}
	return host
}
