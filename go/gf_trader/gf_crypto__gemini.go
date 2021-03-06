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
	"strconv"
	"github.com/gorilla/websocket"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type gf_market_data_parsed_event struct {
	events_id_str   string
	type_str        string
	remote_type_str string //type of event in the remote market system
	msg_str         string
	data_map        map[string]interface{}
}

//-------------------------------------------------
func quote__init_persist_stream(p_symbol_str string,
	p_symbol_name_str string,
	p_runtime         *Runtime) chan float64 {
	p_runtime.Runtime_sys.Log_fun("FUN_ENTER", "gf_crypto__gemini.quote__init_persist_stream()")

	price_updates__ch := make(chan float64, 100)
	go func() {

		price__last_f := 0.0
		for {
			select {
				case new_price_f := <-price_updates__ch:

					trade_time_f            := float64(time.Now().UnixNano())/1000000000.0
					price_f                 := new_price_f
					price__change_nominal_f := price_f - price__last_f
					price__change_percent_f := (100*price__change_nominal_f)/price__last_f

					_, gf_err := quote__create(p_symbol_str,
						p_symbol_name_str,
						trade_time_f,
						price_f,
						price__change_nominal_f,
						price__change_percent_f,
						p_runtime)
					if gf_err != nil {
						panic("cant create quote")
					}

					price__last_f = price_f
			}
		}
	}()
	return price_updates__ch
}

//-------------------------------------------------
func gemini__init(p_runtime *Runtime) {

	gemini__init_symbol("ETHUSD", "etherium_dollar", p_runtime)
	gemini__init_symbol("BTCUSD", "bitcoin_dollar",  p_runtime)
}

//-------------------------------------------------
func gemini__init_symbol(p_symbol_str string,
	p_symbol_name_str string,
	p_runtime         *Runtime) *gf_core.Gf_error {
	p_runtime.Runtime_sys.Log_fun("FUN_ENTER", "gf_crypto__gemini.gemini__init()")

	//--------------------
	//url_str := "wss://api.gemini.com/v1/marketdata/BTCUSD"
	url_str := fmt.Sprintf("wss://api.gemini.com/v1/marketdata/%s", p_symbol_str)

	var ws_dialer *websocket.Dialer
	c, _, err := ws_dialer.Dial(url_str, nil)
	if err != nil {
		gf_err := gf_core.Error__create("failed to connect to Gemini marketdata websocket url",
			"ws_connection_init_error",
			&map[string]interface{}{"url_str": url_str,},
			err, "gf_trader", p_runtime.Runtime_sys)
		return gf_err
	}
	//--------------------
	//QUOTE PERSISTING
	/*price_updates__ch := make(chan float64,100)
	go func() {

		price__last_f := 0.0
		for {
			select {
				case new_price_f := <-price_updates__ch:

					trade_time_f            := float64(time.Now().UnixNano())/1000000000.0
					price_f                 := new_price_f
					price__change_nominal_f := price_f - price__last_f
					price__change_percent_f := (100*price__change_nominal_f)/price__last_f

					_, gf_err := quote__create(p_symbol_str,
						p_symbol_name_str,
						trade_time_f,
						price_f,
						price__change_nominal_f,
						price__change_percent_f,
						p_runtime)
					if gf_err != nil {
						panic("cant create quote")
					}

					price__last_f = price_f
			}
		}
	}()*/

	price_updates__ch := quote__init_persist_stream(p_symbol_str, p_symbol_name_str, p_runtime)
	//--------------------
	//REMOTE_SYSTEM_STREAM_PROCESSING

	go func() {
		defer c.Close()

		for {

			//--------------------
			//READ_GEMINI_MESSAGE
			message_map := map[string]interface{}{}
			err := c.ReadJSON(&message_map)

			if err != nil {
				fmt.Println("read:", err)
				return
			}
			//fmt.Println("---- message - "+fmt.Sprint(message_map))
			//type_str        := message_map["type"].(string) //"update"
			//event_id_int    := message_map["eventId"].(int)
			//timestamp_int   := message_map["timestamp"].(int)
			//timestampms_int := message_map["timestampms"].(int)
			//--------------------

			market_events_lst := message_map["events"].([]interface{})
			parsed_events_lst := gemini__parse_message(p_symbol_str, market_events_lst)

			for _, parsed_event := range parsed_events_lst {
				
				//if the type of event is a "trade" in the remote market system,
				//then just send the clients of the price_feed the price update
				if parsed_event.remote_type_str == "trade" {
					price_updates__ch <- parsed_event.data_map["e__price_f"].(float64) //e__price_f
				}
				
				//-----------------------
				//SEND_EVENT
				gf_core.Events__send_event(parsed_event.events_id_str,
					parsed_event.type_str, //p_type_str
					parsed_event.msg_str,  //p_msg_str
					parsed_event.data_map, //p_data_map
					p_runtime.Events_ctx,
					p_runtime.Runtime_sys)
				//-----------------------
			}
		}
	}()
	//--------------------
	return nil
}

//-------------------------------------------------
func gemini__parse_message(p_symbol_str string,
	p_market_events_lst []interface{}) []*gf_market_data_parsed_event {

	parsed_events_lst := []*gf_market_data_parsed_event{}
	for _, market_event := range p_market_events_lst {

		market_event_map := market_event.(map[string]interface{})

		//the price of this order book entry.
		e__price_f, _ := strconv.ParseFloat(market_event_map["price"].(string),32)
		e__type_str   := market_event_map["type"].(string) //"change" - its always that value
		
		events_id_str  := "trader_gemini_events"
		event_type_str := "gemini_market_update"
		event_msg_str  := fmt.Sprintf("%s market update", p_symbol_str)
		event_data_map := map[string]interface{}{
			"e__symbol_str": p_symbol_str,
			"e__price_f":    e__price_f,
			"e__type_str":   e__type_str,
		}

		parsed_event := &gf_market_data_parsed_event{
			events_id_str: events_id_str,
			type_str:      event_type_str,
			msg_str:       event_msg_str,
			data_map:      event_data_map, 
		}
		parsed_events_lst = append(parsed_events_lst, parsed_event)
		

		//e__reason_str - "place"|"trade"|"cancel"|"initial"
		//                indicates why the "change" (e__type_str) has occurred.
		if e__reason_str, ok := market_event_map["reason"].(string); ok {

			parsed_event.remote_type_str = e__reason_str

			//ADD!! - handle the initial data, that represents the market orders
			//        that are active before the WebSockets clients connected 
			//        to Gemini servers.
			if e__reason_str == "initial" {

			}

			event_data_map["e__reason_str"] = e__reason_str;
		}

		//"bid"|"ask"
		if e__side_str, ok := market_event_map["side"].(string); ok {
			event_data_map["e__side_str"] = e__side_str
		}

		//REMAINING_ORDER_ETH/BTC - amount still remaining of the original market order?
		if remaining_str, ok := market_event_map["remaining"].(string); ok {
			e__remaining_f, _ := strconv.ParseFloat(remaining_str, 32)
			event_data_map["e__remaining_f"] = e__remaining_f
		}

		if delta_str, ok := market_event_map["delta"].(string); ok {
			e__delta_f, _ := strconv.ParseFloat(delta_str, 32)
			event_data_map["e__delta_f"] = e__delta_f
		}
	}
	return parsed_events_lst
}