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
	"sort"
	"github.com/globalsign/mgo/bson"
	"github.com/FlashBoys/go-finance"
	"github.com/gloflow/gloflow/go/gf_core"
)
//-------------------------------------------------
type Gf_quote__day_historical struct {
	Id                   bson.ObjectId `bson:"_id,omitempty"`
	Id_str               string        `bson:"id_str"               json:"-"`              
	T_str                string        `bson:"t"                    json:"-"` //"gf_quote__day_historical"
	Creation_unix_time_f float64       `bson:"creation_unix_time_f" json:"creation_unix_time_f"`
	Symbol_str           string        `bson:"symbol_str"           json:"symbol_str"`
	Date_f               float64       `bson:"date_f"               json:"date_f"` //date for which the quote is actually for (day)
	Open_price_f         float64       `bson:"open_price_f"         json:"open_price_f"`
	High_price_f         float64       `bson:"high_price_f"         json:"high_price_f"`
	Low_price_f          float64       `bson:"low_price_f"          json:"low_price_f"`
	Close_price_f        float64       `bson:"close_price_f"        json:"close_price_f"`
	Volume_int           int           `bson:"volume_int"           json:"volume_int"`
}
//-------------------------------------------------
type quotes__day_historical []*Gf_quote__day_historical
func (d_lst quotes__day_historical) Len() int {
    return len(d_lst)
}
func (d_lst quotes__day_historical) Swap(i, j int) {
    d_lst[i],d_lst[j] = d_lst[j],d_lst[i]
}
func (d_lst quotes__day_historical) Less(i, j int) bool {
    return d_lst[i].Date_f > d_lst[j].Date_f
}

func quotes_historical__get(p_symbol_str string, p_runtime *Runtime) ([]*Gf_quote__day_historical, *gf_core.Gf_error) {

	//--------------------
	//HACK!! - for some reason when going 1 month in the past GetQuoteHistory() will returns
	//         all stock prices in the history of the stock. 
	//         instead 2 months are used (end.AddDate(0,-2,0)) which gives a proper 
	//         1 month of stock quotes. 

	//start,_ := time.Parse(time.RFC3339,"2016-12-01T16:00:00+00:00")
	end   := time.Now() //start.AddDate(0, 1, 0) //1 month period
	start := end.AddDate(0,-2,0) //time.Date(2016,time.November,1,16,0,0,0,time.UTC)//
    //--------------------

	fmt.Println(">>>>>> HISTORICAL >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println("start - "+fmt.Sprint(start))
	fmt.Println("end   - "+fmt.Sprint(end))


    // Request daily history for TWTR.
    // IntervalDaily OR IntervalWeekly OR IntervalMonthly are supported.
    bars_lst, err := finance.GetQuoteHistory(p_symbol_str, start, end, finance.IntervalDaily)
    if err != nil {
        return nil, err
    }

    gf_quotes_lst := []*Gf_quote__day_historical{}
    for _, b := range bars_lst {

    	fmt.Println("---")
    	fmt.Println("b.Date   - "+fmt.Sprint(b.Date))
    	fmt.Println("b.Open   - "+fmt.Sprint(b.Open))
    	fmt.Println("b.High   - "+fmt.Sprint(b.High))
    	fmt.Println("b.Low    - "+fmt.Sprint(b.Low))
    	fmt.Println("b.Close  - "+fmt.Sprint(b.Close))
    	fmt.Println("b.Volume - "+fmt.Sprint(b.Volume))


    	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
		id_str               := "gf_quote_day_historical__"+fmt.Sprint(creation_unix_time_f)
		date_f               := float64(b.Date.Unix())
		open_price_f,_       := b.Open.Float64()
		high_price_f,_       := b.High.Float64()
		low_price_f,_        := b.Low.Float64()
		close_price_f,_      := b.Close.Float64()

    	gf_quote := &Gf_quote__day_historical{
    		Id_str:               id_str,
    		T_str:                "gf_quote__day_historical",
    		Creation_unix_time_f: creation_unix_time_f,
    		Symbol_str:           p_symbol_str,
    		Date_f:               date_f,        //b.Date,
    		Open_price_f:         open_price_f,  //b.Open,
    		High_price_f:         high_price_f,  //b.High,
    		Low_price_f:          low_price_f,   //b.Low,
    		Close_price_f:        close_price_f, //b.Close,
    		Volume_int:           b.Volume,
    	}

    	//-------------
    	//DB
    	gf_err := quote_historical__persist(gf_quote, p_runtime)
    	if gf_err != nil {
    		return nil, gf_err
    	}
    	//-------------

    	gf_quotes_lst = append(gf_quotes_lst, gf_quote)
    }

    //--------------------------------
    //GetQuoteHistory() - returns quotes where the latest is first in the list. 
    //                    this has to be reversed, where the oldest quote is first in the list
    sort.Sort(sort.Reverse(quotes__day_historical(gf_quotes_lst)))
    //--------------------------------

    return gf_quotes_lst, nil
}
//-------------------------------------------------
func quote_historical__persist(p_gf_quote *Gf_quote__day_historical, p_runtime *Runtime) *gf_core.Gf_error {

	//create new historical record if one for this date doesnt already 
	//exist in the DB
	c, err := p_runtime.Runtime_sys.Mongodb_coll.Find(bson.M{
			"t":          "gf_quote__day_historical",
			"symbol_str": p_gf_quote.Symbol_str,
			"date_f":     p_gf_quote.Date_f,
		}).Count()

	if err != nil {
		gf_err := gf_core.Error__create("failed to get a gf_quote__day_historical in the DB",
			"mongodb_find_error",
			&map[string]interface{}{"symbol_str": p_gf_quote.Symbol_str,},
			err, "gf_trader", p_runtime.Runtime_sys)
		return gf_err
	}
	
	if c == 0 {
		err := p_runtime.Runtime_sys.Mongodb_coll.Insert(p_gf_quote)
		if err != nil {
			gf_err := gf_core.Error__create("failed to insert an gf_quote into the DB",
				"mongodb_insert_error",
				&map[string]interface{}{"quote_symbol_str":p_gf_quote.Symbol_str,},
				err, "gf_trader", p_runtime.Runtime_sys)
			return gf_err
		}
	}

	return nil
}