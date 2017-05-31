package buffalo

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"sync"

	"github.com/monoculum/formam"
	"github.com/pkg/errors"
)

var binderLock = &sync.Mutex{}
var binders = map[string]BinderFunc{}
var decoder *formam.Decoder

// BinderFunc takes a request and binds it to an interface.
// If there is a problem it should return an error.
type BinderFunc func(*http.Request, interface{}) error

// RegisterBinder maps a request Content-Type (application/json)
// to a BinderFunc.
func RegisterBinder(contentType string, fn BinderFunc) {
	binderLock.Lock()
	defer binderLock.Unlock()
	binders[contentType] = fn
}

func init() {
	decoder = formam.NewDecoder(&formam.DecoderOptions{TagName: "schema"})
	sb := func(req *http.Request, i interface{}) error {
		err := req.ParseForm()
		if err != nil {
			return errors.WithStack(err)
		}

		if err := decoder.Decode(req.Form, i); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}

	binders["application/html"] = sb
	binders["text/html"] = sb
	binders["application/x-www-form-urlencoded"] = sb
	binders["multipart/form-data"] = sb
}

func init() {
	jb := func(req *http.Request, value interface{}) error {
		return json.NewDecoder(req.Body).Decode(value)
	}
	binders["application/json"] = jb
	binders["text/json"] = jb
	binders["json"] = jb
}

func init() {
	xb := func(req *http.Request, value interface{}) error {
		return xml.NewDecoder(req.Body).Decode(value)
	}
	binders["application/xml"] = xb
	binders["text/xml"] = xb
	binders["xml"] = xb
}
