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

func DetachServerFromResource(locator *service.Locator, serverId string, resourceId string) error {
	return ResourcePatch(locator, resourceId, func(r *Resource) *Resource {
		newDetails := make([]ResourceServer, 0)
		for _, s := range r.ServerDetails {
			if s.ServerId != serverId {
				newDetails = append(newDetails, s)
			}
		}
		r.ServerDetails = newDetails
		return r
	})
}

func ResourceGetServerIds(locator *service.Locator, resourceId string) ([]string, error) {
	resource, err := ResourceGet(locator, resourceId)

	if err != nil {
		return nil, err
	}

	ids := make([]string, 0)

	for _, s := range resource.ServerDetails {
		// ensure the server exists
		_, err := ServerGet(locator, s.ServerId)
		if err == nil {
			ids = append(ids, s.ServerId)
		}
	}

	return ids, nil
}
