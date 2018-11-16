package geoserver

import (
	"encoding/json"
	"fmt"
)

//GroupPublishableItem geoserver Group
type GroupPublishableItem struct {
	Type string `json:"@type,omitempty" xml:"type"`
	Name string `json:"name,omitempty" xml:"name"`
	Href string `json:"href,omitempty" xml:"href"`
}

//Publishables Geoserver Published Layers
type Publishables struct {
	Published []*GroupPublishableItem `json:"published" xml:"published"`
}

//UnmarshalJSON custom deserialization
func (u *Publishables) UnmarshalJSON(data []byte) error {
	var raw interface{}
	json.Unmarshal(data, &raw)
	switch raw := raw.(type) {
	case interface{}:
		published := raw.(map[string]interface{})["published"].(map[string]interface{})
		*u = Publishables{Published: []*GroupPublishableItem{{Name: published["name"].(string), Href: published["href"].(string)}}}
	case map[string]interface{}:
		var publishables Publishables
		err := json.Unmarshal(data, &publishables)
		if err != nil {
			return err
		}
		*u = publishables
	}
	return nil
}

//LayerGroupStyles geoserver layergroup styles
type LayerGroupStyles struct {
	Style []*Resource `json:"style,omitempty" xml:"style"`
}

//LayerGroup geoserver layergroup details
type LayerGroup struct {
	Name         string            `json:"name,omitempty" xml:"name"`
	Mode         string            `json:"mode,omitempty" xml:"mode"`
	Title        string            `json:"title,omitempty" xml:"title"`
	Workspace    *Resource         `json:"workspace,omitempty" xml:"workspace"`
	Publishables Publishables      `json:"publishables" xml:"publishables"`
	Styles       LayerGroupStyles  `json:"styles,omitempty" xml:"styles"`
	Bounds       NativeBoundingBox `json:"bounds" xml:"bounds"`
}

type layerGroupResponse struct {
	LayerGroups struct {
		LayerGroup []*Resource `json:"layerGroup,omitempty"`
	} `json:"layerGroups,omitempty"`
}
type layerGroupDetailsResponse struct {
	LayerGroup *LayerGroup `json:"layerGroup,omitempty"`
}

// LayerGroupService define  geoserver layergroup operations
type LayerGroupService interface {
	GetLayerGroups(workspaceName string) (layerGroups []*Resource, err error)
	GetLayerGroup(workspaceName string, layerGroupName string) (layer *LayerGroup, err error)
	CreateLayerGroup(workspaceName string, layerGroup *LayerGroup) (created bool, err error)
}

//GetLayerGroups  get all layergroups from workspace in geoserver else return error,
//if workspace is "" the it will return all public layers in geoserver
func (g *GeoServer) GetLayerGroups(workspaceName string) (layerGroups []*Resource, err error) {
	if workspaceName != "" {
		workspaceName = fmt.Sprintf("workspaces/%s/", workspaceName)
	}
	targetURL := g.ParseURL("rest", workspaceName, "layergroups")
	httpRequest := HTTPRequest{
		Method: getMethod,
		Accept: jsonType,
		URL:    targetURL,
		Query:  nil,
	}
	response, responseCode := g.DoRequest(httpRequest)
	if responseCode != statusOk {
		g.logger.Error(string(response))
		layerGroups = nil
		err = g.GetError(responseCode, response)
		return
	}
	var layerGroupList layerGroupResponse
	g.DeSerializeJSON(response, &layerGroupList)
	layerGroups = layerGroupList.LayerGroups.LayerGroup
	return
}

//GetLayerGroup get specific LayerGroup in a workspace from geoserver else return error,
//if workspace is "" the it will return geoserver public layer with ${layerName}
func (g *GeoServer) GetLayerGroup(workspaceName string, layerGroupName string) (layerGroup *LayerGroup, err error) {
	if workspaceName != "" {
		workspaceName = fmt.Sprintf("workspaces/%s/", workspaceName)
	}
	targetURL := g.ParseURL("rest", workspaceName, "layergroups", layerGroupName)
	httpRequest := HTTPRequest{
		Method: getMethod,
		Accept: jsonType,
		URL:    targetURL,
		Query:  nil,
	}
	response, responseCode := g.DoRequest(httpRequest)
	if responseCode != statusOk {
		g.logger.Error(string(response))
		layerGroup = &LayerGroup{}
		err = g.GetError(responseCode, response)
		return
	}
	var layerGroupResponse layerGroupDetailsResponse
	g.DeSerializeJSON(response, &layerGroupResponse)
	layerGroup = layerGroupResponse.LayerGroup
	return
}

//CreateLayerGroup create specific LayerGroup in geoserver return created=true else created=false and the error,
//if workspace is "" the it will return geoserver public layer with ${layerName}
// func (g *GeoServer) CreateLayerGroup(workspaceName string, layerGroup *LayerGroup) (created bool, err error) {
// 	if workspaceName != "" {
// 		workspaceName = fmt.Sprintf("workspaces/%s/", workspaceName)
// 	}
// 	group := layerGroupDetailsResponse{LayerGroup: layerGroup}
// 	serializedGroup, _ := g.SerializeStruct(group)
// 	targetURL := g.ParseURL("rest", workspaceName, "layergroups")
// 	data := bytes.NewBuffer(serializedGroup)
// 	fmt.Printf("\n\n\n\n")
// 	fmt.Println(data)
// 	httpRequest := HTTPRequest{
// 		Method:   postMethod,
// 		Accept:   jsonType,
// 		Data:     data,
// 		DataType: jsonType,
// 		URL:      targetURL,
// 		Query:    nil,
// 	}
// 	response, responseCode := g.DoRequest(httpRequest)
// 	if responseCode != statusCreated {
// 		g.logger.Error(string(response))
// 		created = false
// 		err = g.GetError(responseCode, response)
// 		return
// 	}
// 	created = true
// 	return
// }
