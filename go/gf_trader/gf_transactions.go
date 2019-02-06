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
	"time"
	"github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)
//-------------------------------------------------
type Gf_transaction struct {
	Id                   bson.ObjectId   `bson:"_id,omitempty"        json:"-"`
	Id_str               string          `bson:"id_str"               json:"id_str"`
	T_str                string          `bson:"t"                    json:"-"` //"transaction"
	Creation_unix_time_f float64         `bson:"creation_unix_time_f" json:"creation_unix_time_f"`
	Symbol_str           string
	Date_f               float64

	Comission_f          float64       `bson:"Comission_f" json:"Comission_f"` 
	Shares_num_int       int           `bson:"Shares_num_int" json:"Shares_num_int"`
	Share_cost_f         float64
	

	Type_str             string  //"buy"|"sell"

	//IMPORTANT!! - who executed the order
	//              "stocktrainer" - order placed in the stocktrainer mobile app,
	//                               and run there in a simulated environment.
	//                               it is then manually entered into gf_trader
	//              "robinhood" - placed in robinhood mobile app
	Executor_type_str string  //"stocktrainer"|"etrade"|"robinhood"

	//IMPORTANT!! - this indicates where the transaction came from.
	//              "gf_trader_ui"  - it was entered using the gf_trader UI
	//              "manual_import" - it was manually imported via the gf_trader_ui
 	Origin_type_str   string  //"gf_trader_ui"|"manual_import"
}

type Gf_transaction__extern__import_input struct {
	Symbol_str        string    `json:"symbol_str"`
	Date              time.Time `json:"date"`
	Comission_f       float64   `json:"comission_f"`
	Shares_num_int    int       `json:"shares_num_int"`
	Share_cost_f      float64   `json:"share_cost_f"`
	Type_str          string    `json:"type_str"`
	Executor_type_str string    `json:"executor_type_str"`
	Origin_type_str   string    `json:"origin_type_str"`
}

type Gf_transaction__extern__execute_input struct {
	Symbol_str        string    `json:"symbol_str"`
	Shares_num_int    int       `json:"shares_num_int"`
	Share_cost_f      float64   `json:"share_cost_f"`
	Type_str          string    `json:"type_str"`       //"buy"|"sell"
}
//-------------------------------------------------
func transact__execute(p_extern_transaction *Gf_transaction__extern__execute_input,
	p_gf_account *Gf_account,
	p_runtime    *Runtime) (*Gf_transaction, *gf_core.Gf_error) {
	p_runtime.Runtime_sys.Log_fun("FUN_ENTER", "gf_transactions.transact__execute()")

	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	id_str               := "gf_transaction__"+fmt.Sprint(creation_unix_time_f)

	gf_transaction := &Gf_transaction{
		Id_str:               id_str,
		T_str:                "gf_transaction",
		Creation_unix_time_f: creation_unix_time_f,
		Symbol_str:           p_extern_transaction.Symbol_str,
		Date_f:               creation_unix_time_f,
		Comission_f:          10.0, //p_extern_transaction.Comission_f,
		Shares_num_int:       p_extern_transaction.Shares_num_int,
		Share_cost_f:         p_extern_transaction.Share_cost_f,
		Type_str:             p_extern_transaction.Type_str,
		Executor_type_str:    "etrade",       //p_extern_transaction.Executor_type_str,
		Origin_type_str:      "gf_trader_ui", //p_extern_transaction.Origin_type_str,
	}

	gf_err := transact__persist(gf_transaction, p_runtime)
	if gf_err != nil {
		return nil, gf_err
	}

	gf_err = account__update(gf_transaction, p_gf_account, p_runtime)
	if gf_err != nil {
		return nil, gf_err
	}

	return gf_transaction, nil
}
//-------------------------------------------------
func transact__import(p_extern_transaction *Gf_transaction__extern__import_input, p_runtime *Runtime) (*Gf_transaction, *gf_core.Gf_error) {
	p_runtime.Runtime_sys.Log_fun("FUN_ENTER", "gf_transactions.transact__import()")

	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	id_str               := "transaction__"+fmt.Sprint(creation_unix_time_f)

	gf_transaction := &Gf_transaction{
		Id_str:               id_str,
		T_str:                "gf_transaction",
		Creation_unix_time_f: creation_unix_time_f,
		Symbol_str:           p_extern_transaction.Symbol_str,
		Date_f:               float64(p_extern_transaction.Date.Unix()),
		Comission_f:          p_extern_transaction.Comission_f,
		Shares_num_int:       p_extern_transaction.Shares_num_int,
		Share_cost_f:         p_extern_transaction.Share_cost_f,
		Type_str:             p_extern_transaction.Type_str,
		Executor_type_str:    p_extern_transaction.Executor_type_str,
		Origin_type_str:      p_extern_transaction.Origin_type_str,
	}

	return gf_transaction, nil
}
//-------------------------------------------------
func transact__persist(p_gf_transaction *Gf_transaction, p_runtime *Runtime) *gf_core.Gf_error {
	p_runtime.Runtime_sys.Log_fun("FUN_ENTER", "gf_transactions.transact__persist()")
	
	err := p_runtime.Runtime_sys.Mongodb_coll.Insert(p_gf_transaction)
	if err != nil {
		gf_err := gf_core.Error__create("failed to insert a gf_transaction into the DB",
			"mongodb_insert_error",
			&map[string]interface{}{
				"account_name_str":p_gf_transaction.Symbol_str,
				"shares_num_int":  p_gf_transaction.Shares_num_int,
				"shares_num_int":  p_gf_transaction.Share_cost_f,
			},
			err, "gf_trader", p_runtime.Runtime_sys)
		return gf_err
	}

	return nil
}