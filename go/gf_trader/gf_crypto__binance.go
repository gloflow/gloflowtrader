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

		}
	}()

	return nil
}