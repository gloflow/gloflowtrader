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
	"net/http"
	"github.com/globalsign/mgo/bson"
	"github.com/fatih/color"
	//finance "github.com/piquette/finance-go"
	"github.com/PuerkitoBio/goquery"
	"github.com/gloflow/gloflow/go/gf_core"

	finance "github.com/FlashBoys/go-finance"
	//"github.com/shopspring/decimal"
)

//-------------------------------------------------
type Gf_quote struct {
	Id                     bson.ObjectId `json:"-"                      bson:"_id,omitempty"`
	Id_str                 string        `json:"id_str"                 bson:"id_str"`
	T_str                  string        `json:"-"                      bson:"t"` //"quote"
	Creation_unix_time_f   float64       `json:"creation_unix_time_f"   bson:"creation_unix_time_f"`
	Symbol_str             string        `json:"symbol_str"             bson:"symbol_str"`
	Name_str               string        `json:"name_str"               bson:"name_str"`
	Trade_time_f           float64       `json:"trade_time_f"           bson:"trade_time_f"`
	Price_f                float64       `json:"price_f"                bson:"price_f"`
	Price_change_nominal_f float64       `json:"price_change_nominal_f" bson:"price_change_nominal_f"`
	Price_change_percent_f float64       `json:"price_change_percent_f" bson:"price_change_percent_f"`
}

//-------------------------------------------------
func test() {

	client := &http.Client{
		/*IMPORTANT!! - golang http lib does not copy user-set headers on redirects, so a manual
		                setting of these headers had to be added, via the CheckRedirect function
		                that gets called on every redirect, which gives us a chance to to re-set
		                user-agent headers again to the correct value*/
		/*CheckRedirect specifies the policy for handling redirects.
        If CheckRedirect is not nil, the client calls it before
        following an HTTP redirect. The arguments req and via are
        the upcoming request and the requests made already, oldest
        first. If CheckRedirect returns an error, the Client's Get
        method returns both the previous Response (with its Body
        closed) and CheckRedirect's error (wrapped in a url.Error)
        instead of issuing the Request req.
        As a special case, if CheckRedirect returns ErrUseLastResponse,
        then the most recent response is returned with its body
        unclosed, along with a nil error.
        If CheckRedirect is nil, the Client uses its default policy,
        which is to stop after 10 consecutive requests.*/
		CheckRedirect:func(req *http.Request, via []*http.Request) error {
			req.Header.Del("User-Agent")
			req.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.1")
			return nil
		},
	}

	url_str := "http://www.nasdaq.com/symbol/yhoo/real-time"
	req, err := http.NewRequest("GET", url_str, nil)
	if err != nil {
		//p_log_fun("ERROR",fmt.Sprint(err))
		return //nil,err
	}

	req.Header.Del("User-Agent")
	req.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.1")

	resp,err := client.Do(req)
	if err != nil {
		//p_log_fun("ERROR",fmt.Sprint(err))
		return //nil,err
	}

	//doc,err := goquery.NewDocument(p_url_str)od kad sam ustao 
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return //nil,err
	}

	fmt.Println(doc)

	text := doc.Find("#qwidget_lastsale").First().Text()
	fmt.Println(text)

	return //doc,nil
}

//-------------------------------------------------
func repeated__get_quotes(p_runtime *Runtime) {
	p_runtime.Runtime_sys.Log_fun("FUN_ENTER", "gf_quotes.repeated__get_quotes()")

	go func() {
		for ;; {

			//----------------------
			sleep_length := time.Second*time.Duration(60*15)
			time.Sleep(sleep_length)
			//----------------------

			symbols_lst := get_symbols()
			_, gf_err   := quotes__get(symbols_lst, p_runtime)
			if gf_err != nil {
				continue
			}
		}
	}()
}

//-------------------------------------------------
func quotes__get(p_stock_symbols_lst []string, p_runtime *Runtime) ([]*Gf_quote, *gf_core.Gf_error) {
	p_runtime.Runtime_sys.Log_fun("FUN_ENTER", "gf_quotes.quotes__get()")

	cyan   := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	//blue   := color.New(color.FgBlue).SprintFunc()
	
	market_open_bool := market__is_open(p_runtime)
	p_runtime.Runtime_sys.Log_fun("INFO","market is open - "+fmt.Sprint(market_open_bool))

	gf_quotes_lst := []*Gf_quote{}
    for _, symbol_str := range p_stock_symbols_lst {

    	p_runtime.Runtime_sys.Log_fun("INFO", yellow("symbol_str - ")+cyan(symbol_str)+" ------------------------------")
    	
		var gf_quote *Gf_quote

		//no quote was found in the DB.
		//this is the case if the stock symbol was not queried before
		//and has never been persisted in the DB
		if ok, gf_err := quote__exists_in_db(symbol_str, p_runtime); !ok {
			if gf_err != nil {
				return nil, gf_err
			}

			p_runtime.Runtime_sys.Log_fun("INFO", fmt.Sprintf("quote for symbol (%s) not in DB", symbol_str))
			new_gf_quote, gf_err := stock_quote__create_new(symbol_str, p_runtime)
			if gf_err != nil {
				return nil, gf_err
			}
			gf_quote = new_gf_quote
			
		//----------------
		//quote exists in DB
		} else {

			//----------------
			//MARKET_CLOSED - get DB quote value
			
			if !market_open_bool {
				p_runtime.Runtime_sys.Log_fun("INFO", "market closed - get from DB")
				db_gf_quote, gf_err := quote__get_from_db(symbol_str, p_runtime)
				if gf_err != nil {
					return nil, gf_err
				}
				gf_quote = db_gf_quote

			} else {
			//----------------
			//MARKET_OPEN
				

				//IMPORTANT!! - check in the DB if a certain threshold time has passed
				//              since the last time a new update of the quote was stored 
				//              in the DB
				current_time_f := float64(time.Now().UnixNano())/1000000000.0
				ok, gf_err     := quote__is_too_old(symbol_str, current_time_f, p_runtime)
				if gf_err != nil {
					return nil, gf_err
				}

				//--------------------
				//IMPORTANT!! - fetch/persist new quote information, if the currently latest one
				//              is too old (up to 15min old).
				//              once a realtime feed is integrated (from a third party) every new record
				//              streamed from the server will be persisted, and this function will only get the latest
				//              record from the DB.
				if ok {
					new_gf_quote, gf_err := stock_quote__create_new(symbol_str, p_runtime)
					if gf_err != nil {
						return nil, gf_err
					}
					gf_quote = new_gf_quote
				}
				//--------------------
			}
			//--------------------
		}

		p_runtime.Runtime_sys.Log_fun("INFO","last trade price - "+cyan(gf_quote.Price_f))
	    gf_quotes_lst = append(gf_quotes_lst, gf_quote)
	}

    return gf_quotes_lst, nil
}

//-------------------------------------------------
func quote__is_too_old(p_symbol_str string,
	p_compare_to_time_f float64,
	p_runtime           *Runtime) (bool, *gf_core.Gf_error) {
	p_runtime.Runtime_sys.Log_fun("FUN_ENTER","gf_quotes.quote__is_too_old()")

	/*var quote *Quote
	err := p_runtime.Mongodb_coll.Find(bson.M{"t":"quote","symbol_str":p_symbol_str}).
								One(quote)
	if err != nil {
		return false,err
	}*/

	gf_quote, gf_err := quote__get_from_db(p_symbol_str, p_runtime)
	if gf_err != nil {
		return false, gf_err
	}

	delta_f      := p_compare_to_time_f - gf_quote.Creation_unix_time_f
	delta_mins_f := delta_f*60

	if delta_mins_f > 5 {
		return true, nil
	} else {
		return false, nil
	}

	return false, nil
}

//-------------------------------------------------
func stock_quote__create_new(p_symbol_str string,
	p_runtime *Runtime) (*Gf_quote, *gf_core.Gf_error) {
	p_runtime.Runtime_sys.Log_fun("FUN_ENTER","gf_quotes.stock_quote__create_new()")

	//q_lst,err := finance.GetQuotes([]string{p_symbol_str,})
	q, err := finance.GetQuote(p_symbol_str) //finance.quote.Get(p_symbol_str)
	if err != nil {
		gf_err := gf_core.Error__create("go-finance failed to get a quote via GetQuote() function",
			"library_error",
			&map[string]interface{}{"symbol_str": p_symbol_str,},
			err, "gf_trader", p_runtime.Runtime_sys)
    	return nil, gf_err
    }

    //q                         := q_lst[0]
    quote_name_str             := q.Name
    trade_time_f               := 0.0 //float64(q.LastTradeTime.t.Unix())
	price__f, _                := q.LastTradePrice.Float64()
	price__change_nominal_f, _ := q.ChangeNominal.Float64()
	price__change_percent_f, _ := q.ChangePercent.Float64()

	gf_quote, gf_err := quote__create(p_symbol_str,
		quote_name_str,
		trade_time_f,
		price__f,
		price__change_nominal_f,
		price__change_percent_f,
		p_runtime)
	if gf_err != nil {
		return nil, gf_err
	}


	//EVENT UPDATE
	if p_runtime.Events_ctx != nil {
		events_id_str  := "trader_events"
		event_type_str := "quote_update"
		event_msg_str  := "quote update for - "+p_symbol_str
		event_data_map := map[string]interface{}{
			"quote": gf_quote,
		}
		gf_core.Events__send_event(events_id_str,
			event_type_str,       //p_type_str
			event_msg_str,        //p_msg_str
			event_data_map,       //p_data_map
			p_runtime.Events_ctx,
			p_runtime.Runtime_sys)
	}
	//--------------

	return gf_quote, nil
}

//-------------------------------------------------
func quote__create(p_symbol_str string,
	p_quote_name_str         string,
	p_trade_time_f           float64,
	p_price_f                float64,
	p_price_change_nominal_f float64,
	p_price_change_percent_f float64,
	p_runtime                *Runtime) (*Gf_quote, *gf_core.Gf_error) {
	p_runtime.Runtime_sys.Log_fun("FUN_ENTER", "gf_quotes.quote__create()")

	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	id_str               := "gf_quote__"+fmt.Sprint(creation_unix_time_f)

	gf_quote := &Gf_quote{
		Id_str:                 id_str,
		T_str:                  "gf_quote",
		Creation_unix_time_f:   creation_unix_time_f,
		Symbol_str:             p_symbol_str,
		Name_str:               p_quote_name_str,
		Trade_time_f:           p_trade_time_f,
		Price_f:                p_price_f, 
		Price_change_nominal_f: p_price_change_nominal_f,
		Price_change_percent_f: p_price_change_percent_f,
	}

	p_runtime.Runtime_sys.Log_fun("INFO", "------- "+p_symbol_str+" - "+fmt.Sprint(gf_quote.Price_f))

	//--------------
	//DB PERSIST
	gf_err := quote__persist(gf_quote, p_runtime)
	if gf_err != nil {
		return nil, gf_err
	}
	//--------------
	
	return gf_quote, nil
}