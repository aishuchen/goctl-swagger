package v2

import (
	"net/http"
	"strings"

	"github.com/aishuchen/goctl-swagger/render/types"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func renderPaths(svc spec.Service, opt types.Option) map[string]*Path {
	paths := make(map[string]*Path)
	tags := svc.Name
	for _, grp := range svc.Groups {
		if value := grp.GetAnnotation("group"); len(value) > 0 {
			tags = value
		}
		if len(opt.TagPrefix) > 0 {
			tags = opt.TagPrefix + tags
		}
		for _, route := range grp.Routes {
			uri := grp.GetAnnotation("prefix") + route.Path
			if uri[0] != '/' {
				uri = "/" + uri
			}
			path, ok := paths[uri]
			if !ok {
				path = new(Path)
			}
			op := Operation{
				Summary: strings.Trim(route.AtDoc.Text, `"`),
				Tags:    []string{tags},
			}
			if obj, ok := route.RequestType.(spec.DefineStruct); ok {
				op.Parameters = renderParameters(obj, strings.ToUpper(route.Method))
			}

			// Just support json response yet.
			if obj, ok := route.ResponseType.(spec.DefineStruct); ok {
				// root schema
				root := &Schema{
					Description: strings.Join(obj.Docs, ","),
				}
				op.Responses = map[string]*Response{
					"200": {
						Description: "OK",
						Schema:      root,
					},
				}
				_, resp := renderSchema(obj)
				ref := registerModel(obj.Name(), resp)
				root.Ref = ref

			} else {
				op.Responses = map[string]*Response{
					"200": {
						Description: "OK",
					},
				}
			}
			switch strings.ToUpper(route.Method) {
			case http.MethodGet:
				path.Get = &op
			case http.MethodPost:
				path.Post = &op
			case http.MethodPut:
				path.Put = &op
			case http.MethodDelete:
				path.Delete = &op
			case http.MethodPatch:
				path.Patch = &op
			case http.MethodHead:
				path.Head = &op
			case http.MethodOptions:
				path.Options = &op
			}
			paths[uri] = path
		}
	}
	return paths
}
