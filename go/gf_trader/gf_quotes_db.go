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
)
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
//-------------------------------------------------
func quote__get_from_db(p_symbol_str string, p_runtime *Runtime) (*Quote,error) {

	var quote *Quote
	err := p_runtime.Runtime_sys.Mongodb_coll.Find(bson.M{"t":"quote", "symbol_str":p_symbol_str}).One(&quote)
	if err != nil {
		return nil, err
	}

	return quote, nil
}
//-------------------------------------------------
func quote__persist(p_quote *Quote, p_runtime *Runtime) error {

	err := p_runtime.Runtime_sys.Mongodb_coll.Insert(p_quote)
	if err != nil {
		return err
	}

	return nil
}
//-------------------------------------------------
func quote__exists_in_db(p_symbol_str string, p_runtime *Runtime) (bool,error) {

	c,err := p_runtime.Runtime_sys.Mongodb_coll.Find(bson.M{"t":"quote", "symbol_str":p_symbol_str}).Count()
	if err != nil {
		return false, err
	}

	if c == 0 {
		return false, nil
	}

	return true, nil
}