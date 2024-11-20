package app

import (
	"github.com/maddalax/htmgo/framework/service"
	"time"
)

func PatchResourceServer(locator *service.Locator, resourceId string, serverId string, cb func(server *ResourceServer) *ResourceServer) error {
	return ResourcePatch(locator, resourceId, func(r *Resource) *Resource {
		for i, s := range r.ServerDetails {
			if s.ServerId == serverId {
				r.ServerDetails[i] = *cb(&s)
				break
			}
		}
		return r
	})
}

func AttachServerToResource(locator *service.Locator, serverId string, resourceId string) error {
	server, err := ServerGet(locator, serverId)

	if err != nil {
		return err
	}

	return ResourcePatch(locator, resourceId, func(r *Resource) *Resource {

		has := false

		for _, s := range r.ServerDetails {
			if s.ServerId == server.Id {
				has = true
				break
			}
		}

		if has {
			return r
		}

		if r.ServerDetails == nil {
			r.ServerDetails = make([]ResourceServer, 0)
		}
		r.ServerDetails = append(r.ServerDetails, ResourceServer{
			ServerId:   server.Id,
			RunStatus:  RunStatusNotRunning,
			LastUpdate: time.Now(),
		})
		return r
	})
}
