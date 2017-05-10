package breeze

import (
	"encoding/json"
	"fmt"
	"github.com/kr/pretty"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func ok(w http.ResponseWriter, c interface{}) {
	w.Header().Set("Content-Type", "application/json")
	buf, _ := json.Marshal(struct {
		Status  string
		Content interface{}
	}{"ok", c})
	w.Write(buf)
}

func ucHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		uc, err := NewUserContextFromJson(decoder)
		GetContextManager().AddUC(uc)
		pretty.Println("New Context: ", uc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ok(w, struct {
			ID string
		}{uc.ID})
		return
	} else if r.Method == "GET" {
		ok(w, struct {
			Contexts map[string]*UserContext
		}{GetContextManager().contexts})
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

}


func straDepoyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		buf, err := ioutil.ReadAll(r.Body)
		// fmt.Println("Buf: ", buf)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s, err := DeployStrategy(buf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ok(w, s)
		return
	} else if r.Method == "GET" {
		values := r.URL.Query()
		s := GetStrategyManager().GetStrategy(values.Get("id"))
		w.Header().Set("Content-Type", "application/zip")
		path := s.path()
		fmt.Println("Path: ", path)
		pkg := &StrategyPkg{}
		GetStore().ReadObject(STRATEGY_PKG_TBL_NAME, s.name(), pkg)
		pretty.Println(pkg)
		w.Write(pkg.Content)
		return
	}
}

func straHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		stra, _ := NewStrategyFromToml(string(buf))
		GetStrategyManager().AddStrategy(stra)
		ok(w, struct {
			ID string
		}{stra.ID})
		return
	} else if r.Method == "DELETE" {
		r.ParseForm()
		sid := r.Form.Get("id")
		n := GetStrategyManager().Remove(sid)
		ok(w, struct {
			Removed int
		}{n})
		return
	} else if r.Method == "GET" {
		values := r.URL.Query()
		fmt.Println(values)
		if values.Get("id") != "" {
			id := values.Get("id")
			s := GetStrategyManager().GetStrategy(id)
			if s != nil {
				ok(w, struct {
					Strategies map[string]*Strategy
				}{map[string]*Strategy{s.ID: s}})
				return
			}
		}
		ret := make(map[string]*Strategy)
		for k, v := range GetStrategyManager().strategies {
			if v.Visiblity {
				ret[k] = v
			}
		}
		ok(w, struct {
			Strategies map[string]*Strategy
		}{ret})
		return
	}
	http.Error(w, "", http.StatusMethodNotAllowed)
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		e := &struct {
			Ucid  string
			Event Event
		}{}
		decoder.Decode(e)
		err := GetContextManager().SendEvent(e.Ucid, &e.Event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ok(w, nil)
		return
	}
	http.Error(w, "", http.StatusMethodNotAllowed)
}

func executorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		executor := &Executor{}
		decoder.Decode(executor)
		GetExecutorManager().AddExecutor(executor.Type, executor)
		ok(w, nil)
		return
	} else if r.Method == "GET" {
		ok(w, struct {
			Executors map[string]IExecutor
		}{GetExecutorManager().Executors})
		return
	}
	http.Error(w, "", http.StatusMethodNotAllowed)
}

func tradeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		values := r.URL.Query()
		uid := values.Get("uid")
		if uid == "" {
			http.Error(w, "Need `uid`", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/csv")
		trades := make([]map[string]interface{}, 0)
		GetStore().ReadBy(TRADE_TBL_NAME, &trades, map[string]string{"uid": uid})

		var keys []string
		for k, _ := range trades[0] {
			keys = append(keys, k)
		}
		buf, _ := encodeCSV(keys, trades)
		w.Write(buf)
		return
	}
	http.Error(w, "", http.StatusMethodNotAllowed)
}

func brokeTradeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		values := r.URL.Query()
		pipe := make(chan interface{})
		start, _ :=  strconv.ParseInt(values.Get("start"), 10, 32)
		stop, _ := strconv.ParseInt(values.Get("stop"), 10, 32)
		event := &Event{
			Type: "get_trades",
			Params: map[string]interface{}{
				"start": start,
				"stop": stop,
			},
			callback: func(elist interface{}) {
				content := elist.(EventList)[0].Params
				pipe <- content
		 	},
		}
		uid := values.Get("client_id")
		GetContextManager().GetUC(uid).Send(event)
		ret := <- pipe
		w.Header().Set("Content-Type", "application/json")
		buf, _ := json.Marshal(map[string]interface{}{
			"trades": ret,
		})
		w.Write(buf)
	}
}

func smartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		e := &struct {
			Ucid  string
			Event Event
		}{}
		decoder.Decode(e)
		pipe := make(chan interface{}, 1)
		e.Event.callback = func(elist interface{}) {
			_elist := elist.(EventList)
			if len(_elist) < 1 {
				pipe <- nil
			} else {
				content := elist.(EventList)[0].Params
				pipe <- content
			}
		}
		GetContextManager().SendEvent(e.Ucid, &e.Event)
		ret := <- pipe
		w.Header().Set("Content-Type", "application/json")
		buf, _ := json.Marshal(ret)
		w.Write(buf)
		return
	}
	http.Error(w, "", http.StatusMethodNotAllowed)
}

func StartServer(addr string) {
	http.HandleFunc("/api/v1/context", ucHandler)
	http.HandleFunc("/api/v1/strategy", straHandler)
	http.HandleFunc("/api/v1/event", eventHandler)
	http.HandleFunc("/api/v1/executor", executorHandler)
	http.HandleFunc("/api/v1/deploy", straDepoyHandler)
	http.HandleFunc("/api/v1/flows", tradeHandler)
	http.HandleFunc("/litscope/trades", brokeTradeHandler)
	http.HandleFunc("/api/v1/smart", smartHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}
