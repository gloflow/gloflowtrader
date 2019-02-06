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
	"github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)
//-------------------------------------------------
func quote__get_from_db(p_symbol_str string, p_runtime *Runtime) (*Gf_quote, *gf_core.Gf_error) {

	var gf_quote *Gf_quote
	err := p_runtime.Runtime_sys.Mongodb_coll.Find(bson.M{"t":"gf_quote", "symbol_str":p_symbol_str}).One(&gf_quote)
	if err != nil {
		gf_err := gf_core.Error__create("failed to get a quote in the DB",
			"mongodb_find_error",
			&map[string]interface{}{"symbol_str": p_symbol_str,},
			err, "gf_trader", p_runtime.Runtime_sys)
		return nil, gf_err
	}
	return gf_quote, nil
}
//-------------------------------------------------
func quote__persist(p_gf_quote *Gf_quote, p_runtime *Runtime) *gf_core.Gf_error {

	err := p_runtime.Runtime_sys.Mongodb_coll.Insert(p_gf_quote)
	if err != nil {
		gf_err := gf_core.Error__create("failed to insert an gf_quote into the DB",
			"mongodb_insert_error",
			&map[string]interface{}{"quote_symbol_str":p_gf_quote.Symbol_str,},
			err, "gf_trader", p_runtime.Runtime_sys)
		return gf_err
	}
	return nil
}
//-------------------------------------------------
func quote__exists_in_db(p_symbol_str string, p_runtime *Runtime) (bool, *gf_core.Gf_error) {

	c, err := p_runtime.Runtime_sys.Mongodb_coll.Find(bson.M{"t":"gf_quote", "symbol_str":p_symbol_str}).Count()
	if err != nil {
		gf_err := gf_core.Error__create("failed to get a quote from DB to check if it exists",
			"mongodb_find_error",
			&map[string]interface{}{"symbol_str": p_symbol_str,},
			err, "gf_trader", p_runtime.Runtime_sys)
		return false, gf_err
	}

	if c == 0 {
		return false, nil
	}

	return true, nil
}
//-------------------------------------------------
/*func (q Quote) GetBSON() (interface{},error) {

	return struct {

		
	}{

	},nil
}
func (q *Quote) SetBSON(p_raw bson.Raw) error {

	decoded := new(struct {
		Id                   bson.ObjectId   `bson:"_id,omitempty"        json:"-"`
		Id_str               string          `bson:"id_str"               json:"id_str"`
		T_str                string          `bson:"t"                    json:"-"` //"quote"
		Creation_unix_time_f float64         `bson:"creation_unix_time_f" json:"creation_unix_time_f"`
		Symbol_str           string          `bson:"symbol_str"           json:"symbol_str"`
		Date                 time.Time       `bson:"date"                 json:"date"`

		//IMPORTANT!! - these are serialized as float64 (decimal.Decimal is a runtime type,
		//              not persistable) from there a "_f" postfix is used
        Last_trade_price_f   float64 `bson:"last_trade_price_f"`
        Change_nominal_f     float64 `bson:"change_nominal_f"`
        Change_percent_f     float64 `bson:"change_percent_f"`
    })

    bson_err := p_raw.Unmarshal(decoded)

    if bson_err == nil {
    	q.Id_str               = decoded.Id_str
    	q.T_str                = decoded.T_str
    	q.Creation_unix_time_f = decoded.Creation_unix_time_f
    	q.Symbol_str           = decoded.Symbol_str
    	q.Date                 = decoded.Date
        q.Last_trade_price_d   = decimal.NewFromFloat(decoded.Last_trade_price_f)
        q.Change_nominal_d     = decimal.NewFromFloat(decoded.Change_nominal_f)
        q.Change_percent_d     = decimal.NewFromFloat(decoded.Change_percent_f)
        return nil
    } else {
        return bson_err
    }
	return nil
}
//-------------------------------------------------
func (q Quote__day_historical) GetBSON() (interface{},error) {

	return struct {

		
	}{

	},nil 
}
func (q *Quote__day_historical) SetBSON(p_raw bson.Raw) error {
	return nil
}*/