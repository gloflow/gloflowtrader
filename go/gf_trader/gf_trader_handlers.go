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

import (
	"fmt"
	"strings"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
)
//-------------------------------------------------
func init_handlers(p_runtime *Runtime) error {

	//-------------------
	//QUOTES
	http.HandleFunc("/trader/quotes/latest", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime.Runtime_sys.Log_fun("INFO", "INCOMING HTTP REQUEST -- /trader/quotes/latest ----------")

		if p_req.Method == "GET" {

			qs_map            := p_req.URL.Query()
			stock_symbols_lst := strings.Split(qs_map["symbols"][0],",")

			quotes_lst, gf_err := quotes__get(stock_symbols_lst, p_runtime)
			if gf_err != nil {

			}
			//------------------
			//OUTPUT
			
			r_map := map[string]interface{}{
				"status_str": "OK",
				"quotes_lst": quotes_lst,
			}

			r_lst,_ := json.Marshal(r_map)
			r_str   := string(r_lst)
			fmt.Fprintf(p_resp, r_str)
			//------------------
		}
	})

	http.HandleFunc("/trader/quotes/daily_historic", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime.Runtime_sys.Log_fun("INFO", "INCOMING HTTP REQUEST -- /trader/quotes/daily_historic ----------")

		if p_req.Method == "GET" {

			qs_map           := p_req.URL.Query()
			stock_symbol_str := qs_map["symbols"][0]

			quotes_lst, gf_err := quotes_historical__get(stock_symbol_str, p_runtime)
			if gf_err != nil {

			}
			//------------------
			//OUTPUT
			
			r_map := map[string]interface{}{
				"status_str": "OK",
				"quotes_lst": quotes_lst,
			}

			r_lst,_ := json.Marshal(r_map)
			r_str   := string(r_lst)
			fmt.Fprintf(p_resp, r_str)
			//------------------
		}
	})

	//-------------------
	//TRANSACTIONS
	
	http.HandleFunc("/trader/transaction/execute", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime.Runtime_sys.Log_fun("INFO", "INCOMING HTTP REQUEST -- /trader/transaction/execute ----------")

		if p_req.Method == "POST" {

			//------------
			//INPUT
			var input Gf_transaction__extern__execute_input
			body_bytes_lst,_ := ioutil.ReadAll(p_req.Body)
		    err              := json.Unmarshal(body_bytes_lst, &input)

			if err != nil {
				gf_err := gf_core.Error__create("failed to parse json http input",
					"json_unmarshal_error",
					&map[string]interface{}{},
					err, "gf_trader", p_runtime.Runtime_sys)

				gf_rpc_lib.Error__in_handler("/trader/transaction/execute",
					"transaction_execute received bad transaction__extern_execute input",
					gf_err, p_resp, p_runtime.Runtime_sys)
				return
			}
			//------------

			account_name_str   := "practice_trading"
			gf_account, gf_err := account__get(account_name_str, p_runtime)
			if gf_err != nil {
				return
			}

			_, gf_err = transact__execute(&input, gf_account, p_runtime)
			if gf_err != nil {
				return
			}
		}
	})

	http.HandleFunc("/trader/transaction/import", func(p_resp http.ResponseWriter, p_req *http.Request) {
		p_runtime.Runtime_sys.Log_fun("INFO", "INCOMING HTTP REQUEST -- /trader/transaction/import ----------")

		if p_req.Method == "POST" {

			//------------
			//INPUT
			var input Gf_transaction__extern__import_input
			body_bytes_lst, _ := ioutil.ReadAll(p_req.Body)
		    err               := json.Unmarshal(body_bytes_lst, &input)

			if err != nil {
				gf_err := gf_core.Error__create("failed to parse json http input",
					"json_unmarshal_error",
					&map[string]interface{}{},
					err, "gf_trader", p_runtime.Runtime_sys)

				gf_rpc_lib.Error__in_handler("/trader/transaction/import",
					"transaction_import received bad transaction__extern_input input",
					gf_err, p_resp, p_runtime.Runtime_sys)
				return
			}
			//------------

			_, gf_err := transact__import(&input, p_runtime)
			if gf_err != nil {
				return
			}
		}
	})
	//-------------------

	return nil
}