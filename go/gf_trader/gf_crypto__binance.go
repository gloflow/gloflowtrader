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
	"strconv"
	"github.com/gorilla/websocket"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
func binance__init(p_runtime *Runtime) {

	binance__init_symbol("ETHUSDT", "etherium_dollar", p_runtime)
	binance__init_symbol("BTCUSDT", "bitcoin_dollar",  p_runtime)
}

//-------------------------------------------------
func binance__init_symbol(p_symbol_str string,
	p_symbol_name_str string,
	p_runtime         *Runtime) *gf_core.Gf_error {
	p_runtime.Runtime_sys.Log_fun("FUN_ENTER", "gf_crypto__binance.binance__init_symbol()")

	//--------------------
	url_str := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@aggTrade", strings.ToLower(p_symbol_str))
	fmt.Println(fmt.Sprintf("url_str - %s", url_str))

	var ws_dialer *websocket.Dialer
	c, _, err := ws_dialer.Dial(url_str, nil)
	if err != nil {
		gf_err := gf_core.Error__create("failed to connect to Binance marketdata websocket url",
			"ws_connection_init_error",
			&map[string]interface{}{"url_str": url_str,},
			err, "gf_trader", p_runtime.Runtime_sys)
		return gf_err
	}
	//--------------------


	price_updates__ch := quote__init_persist_stream(p_symbol_str, p_symbol_name_str, p_runtime)

	go func() {
		defer c.Close()

		for {

			fmt.Println("----")
			//--------------------
			//READ_MESSAGE
			message_map := map[string]interface{}{}
			err := c.ReadJSON(&message_map)

			if err != nil {
				fmt.Println("read:", err)
				return
			}
			fmt.Println(message_map)
			//--------------------
			market_event_map := message_map
			parsed_event     := binance__parse_message(p_symbol_str, market_event_map)

			price_updates__ch <- parsed_event.data_map["e__price_f"].(float64) //e__price_f

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
	}()

	return nil
}

//-------------------------------------------------
func binance__parse_message(p_symbol_str string,
	p_market_event map[string]interface{}) *gf_market_data_parsed_event {
	//p_runtime.Runtime_sys.Log_fun("FUN_ENTER", "gf_crypto__binance.binance__parse_message()")


	/*{
		"e": "aggTrade",  // Event type
		"E": 123456789,   // Event time
		"s": "BNBBTC",    // Symbol
		"a": 12345,       // Aggregate trade ID
		"p": "0.001",     // Price
		"q": "100",       // Quantity
		"f": 100,         // First trade ID
		"l": 105,         // Last trade ID
		"T": 123456785,   // Trade time
		"m": true,        // Is the buyer the market maker?
		"M": true         // Ignore
	}*/

	event_remote_time_f        := p_market_event["E"].(float64)
	trade_remote_time_f        := p_market_event["T"].(float64)
	trade_is_market_maker_bool := p_market_event["m"].(bool)
	e__price_str               := p_market_event["p"].(string)
	e__price_f, _              := strconv.ParseFloat(e__price_str, 64)


	events_id_str  := "trader_binance_events"
	event_type_str := "binance_market_update"
	event_msg_str  := fmt.Sprintf("%s market update", p_symbol_str)
	event_data_map := map[string]interface{}{
		"e__symbol_str":                 p_symbol_str,
		"e__price_f":                    e__price_f,
		"e__event_remote_time_f":        event_remote_time_f,
		"e__trade_remote_time_f":        trade_remote_time_f,
		"e__trade_is_market_maker_bool": trade_is_market_maker_bool,
	}

	parsed_event := &gf_market_data_parsed_event{
		events_id_str: events_id_str,
		type_str:      event_type_str,
		msg_str:       event_msg_str,
		data_map:      event_data_map, 
	}

	return parsed_event
}