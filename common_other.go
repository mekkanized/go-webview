// go:build !windows
package webview

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type Hint int

const (
	// HintNone specifies that width and height are default size
	HintNone Hint = iota

	// HintFixed specifies that window size can not be changed by a user
	HintFixed

	// HintMin specifies that width and height are minimum bounds
	HintMin

	// HintMax specifies that width and height are maximum bounds
	HintMax
)

func (w *webview) Bind(name string, f interface{}) error {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		return fmt.Errorf("only functions can be bound")
	}

	if n := v.Type().NumOut(); n > 2 {
		return fmt.Errorf("function may only return a value or a value+error")
	}

	w.mutex.Lock()
	w.bindings[name] = f
	w.mutex.Unlock()

	js := fmt.Sprintf(`(function() { var name = '%s';
		var RPC = window._rpc = (window._rpc || {nextSeq: 1});
		window[name] = function() {
			var seq = RPC.nextSeq++;
			var promise = new Promise(function(resolve, reject) {
				RPC[seq] = {
					resolve: resolve,
					reject: reject
				};
			});
			window.external.invoke(JSON.stringify({
				id: seq,
				method: name,
				params: Array.prototype.slice.call(arguments),
			}));
			return promise;
		};
	})();`, name)
	w.Init(js)
	w.Eval(js)

	fmt.Println("completed bind", name)

	return nil
}

type rpcMessage struct {
	ID     int               `json:"id"`
	Method string            `json:"method"`
	Params []json.RawMessage `json:"params"`
}

func (w *webview) onMessage(msg string) {
	var req rpcMessage
	if err := json.Unmarshal([]byte(msg), &req); err != nil {
		log.Printf("invalid RPC message: %v", err)
		return
	}

	defer w.Eval(fmt.Sprintf(`delete window._rpc[%d];`, req.ID))

	res, err := w.callBinding(req)
	if err != nil {
		w.Eval(fmt.Sprintf(`window._rpc[%d].reject(%s);`, req.ID, err.Error()))
		return
	}

	serRes, err := json.Marshal(res)
	if err != nil {
		w.Eval(fmt.Sprintf(`window._rpc[%d].reject(%s);`, req.ID, err.Error()))
		return
	}

	w.Eval(fmt.Sprintf(`window._rpc[%d].resolve(%s);`, req.ID, serRes))
}

func (w *webview) callBinding(req rpcMessage) (interface{}, error) {
	w.mutex.RLock()
	f, ok := w.bindings[req.Method]
	w.mutex.RUnlock()
	if !ok {
		return nil, nil
	}

	v := reflect.ValueOf(f)
	isVariadic := v.Type().IsVariadic()
	numIn := v.Type().NumIn()
	if (isVariadic && len(req.Params) < numIn-1) || (!isVariadic && len(req.Params) != numIn) {
		return nil, fmt.Errorf("function arguments mismatch: expected %d, got %d", numIn, len(req.Params))
	}

	args := []reflect.Value{}
	for i := range req.Params {
		var arg reflect.Value
		if isVariadic && i >= numIn-1 {
			arg = reflect.New(v.Type().In(numIn - 1).Elem())
		} else {
			arg = reflect.New(v.Type().In(i))
		}
		if err := json.Unmarshal(req.Params[i], arg.Interface()); err != nil {
			return nil, fmt.Errorf("failed to unmarshal argument: %w", err)
		}
		args = append(args, arg.Elem())
	}

	errorType := reflect.TypeOf((*error)(nil)).Elem()
	res := v.Call(args)
	switch len(res) {
	case 0:
		// No results from the function, just return nil
		return nil, nil
	case 1:
		// One result may be a value, or an error
		if res[0].Type().Implements(errorType) {
			if res[0].Interface() != nil {
				return nil, res[0].Interface().(error)
			}
			return nil, nil
		}
		return res[0].Interface(), nil
	case 2:
		// Two results: first one is value, second one is error
		if !res[1].Type().Implements(errorType) {
			return nil, fmt.Errorf("second return value must be an error, got %s", res[1].Type().String())
		}
		if res[1].Interface() != nil {
			return res[0].Interface(), nil
		}
		return res[0].Interface(), res[1].Interface().(error)
	default:
		return nil, fmt.Errorf("unexpected number of return values: %d", len(res))
	}
}
