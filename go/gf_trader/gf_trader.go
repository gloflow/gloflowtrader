/*
GloFlowTrader asset trading, management, and research platform
Copyright (C) 2019 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

package main

//-------------------------------------------------
import (
	"fmt"
	"time"
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type Runtime struct {
	Events_ctx  *gf_core.Events_ctx
	Runtime_sys *gf_core.Runtime_sys
}

//-------------------------------------------------
func main() {
	log_fun := gf_core.Init_log_fun()

	port_str := "4400"
	//-----------------
	//MONGODB

	mongodb_host_str    := "localhost"
	mongodb_db_name_str := "gf_trader"
	mongodb_db          := gf_core.Mongo__connect(mongodb_host_str, mongodb_db_name_str, log_fun)
	mongodb_coll        := mongodb_db.C("data")
	//-----------------

	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str: "gf_trader",
		Log_fun:          log_fun,
		Mongodb_coll:     mongodb_coll,
	}
	
	events_ctx := gf_core.Events__init("/trader/events", runtime_sys)
	runtime    := &Runtime{
		Events_ctx:  events_ctx,
		Runtime_sys: runtime_sys,
	}
	
	//STATIC CONFIG 
	//symbols_lst := get_symbols()

	init_handlers(runtime)

	/*//------------
	//RUN_QUOTES_QUERY
	//this is for debugging, to view quotes HTTP query working and its output in the console
	quotes__get(symbols_lst,
			runtime,
			log_fun)
	//------------*/

	/*quotes__get_history(symbols_lst[0],
					log_fun)*/

	/*//START QUERIES_HTTP_UPDATE LOOP
	repeated__get_quotes(runtime,
					log_fun)*/

	account__create_defaults(runtime)

	//------------------------
	//DASHBOARD SERVING
	static_files__url_base_str := "/trader"
	gf_core.HTTP__init_static_serving(static_files__url_base_str, runtime_sys)
	//------------------------
	//CRYPTO_EXCHANGES
	gemini__init(runtime)
	binance__init(runtime)
	//------------------------
	//test()

	log_fun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	log_fun("INFO", "STARTING HTTP SERVER - PORT - "+port_str)
	log_fun("INFO", "http://localhost:4400/trader/static/templates/gf_trader/gf_trader.html")
	log_fun("INFO", ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	
	http_err := http.ListenAndServe(":"+port_str, nil)
	if http_err != nil {
		msg_str := "cant start listening on port - "+port_str
		log_fun("ERROR", msg_str)
		log_fun("ERROR", fmt.Sprint(http_err))
		panic(fmt.Sprint(http_err))
	}
}

//-------------------------------------------------
func market__is_open(p_runtime *Runtime) bool {
	p_runtime.Runtime_sys.Log_fun("FUN_ENTER", "gf_trader.market__is_open()")

	t := time.Now()

	//market not open over the weekend
	if t.Weekday().String() == "Saturday" || t.Weekday().String() == "Sunday" {
		return false
	} else {

		//market hasnt open yet, or has closed, for the day
		if (t.Hour() <= 9 && t.Minute() < 30.0) || t.Hour() >= 16 {
			return false
		}
	}

	return true
}