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

import * as gf_trader              from "./gf_trader";
import * as gf_trader_transactions from "./gf_trader_transactions";

declare var EventSource;
declare var c3;

//---------------------------------------------------
export function init_updates(p_log_fun) {
    p_log_fun('FUN_ENTER', 'gf_trader_quotes.init_updates()');

    const events_id_str = "trader_events";
    const event_source  = new EventSource("/trader/events?events_id="+events_id_str)

    console.log("REGISTER EVENT_SOURCE");

    var i=0;
    event_source.onmessage = (p_e) => {
        console.log('>>>>> MESSAGE');
        const event_data_map = JSON.parse(p_e.data);
        
        console.log(event_data_map)        
        view__update(event_data_map, p_log_fun);

        i+=1;
    }

    event_source.onerror = (p_e)=>{
        console.log('EventSource >> ERROR - '+event_source.readyState);
        console.log(EventSource.CLOSED)
        console.log(p_e);
          
        //connection was closed
        if (event_source.readyState == EventSource.CLOSED) {
            console.log("EVENT_SOURCE CLOSED")
        }
    }

    event_source.onopen = (p_e)=>{
        console.log('EventSource >> OPEN CONN');
    }
}
//---------------------------------------------------
export function view__update(p_quote_update_map, p_log_fun) {




	
}
//---------------------------------------------------
export function view_quote(p_q_map, p_log_fun) {

	const symbol_str             = p_q_map['symbol_str'];
	const name_str               = p_q_map['name_str'];
	const price_f                = p_q_map['price_f'];
	const price_change_nominal_f = p_q_map['price_change_nominal_f'];
	const price_change_percent_f = p_q_map['price_change_percent_f'];

    const q = $(`
    	<div class='quote'>
    		<div class='info_display'>
    			<div class='symbol'>`+symbol_str+`</div>
    			<div class='name'>`+name_str+`</div>
    			<div class='last_trade_price'>$`+price_f+`</div>
    			<div class='change_nominal'>$`+price_change_nominal_f+`</div>
    			<div class='change_percent'>%`+price_change_percent_f+`</div>
    		</div>
    		<div class='view_plot_btn'>plot</div>
    		<div class='plots'></div>

    		<div class='transactions'>
    			<div class='buy_btn'>buy</div>
    			<div class='sell_btn'>sell</div>
    		</div>
    	</div>`);

    if (price_change_nominal_f > 0) {
    	$(q).find('.change_nominal').css('background-color', '#80ff80');
    } else {
    	$(q).find('.change_nominal').css('background-color', 'red');
    }

    if (price_change_percent_f > 0) {
    	$(q).find('.change_percent').css('background-color', '#80ff80');
    } else {
    	$(q).find('.change_percent').css('background-color', 'red');
    }

	var last_trade_price_f = 0.0; //FIX!!
    $(q).find('.buy_btn').on('click',()=>{
    	const b = gf_trader_transactions.view__buy_dialog(symbol_str, last_trade_price_f, p_log_fun);
    	$(q).find('.transactions').append(b);
    });

    $(q).find('.view_plot_btn').on('click', (p_e)=>{

    	//------------------------
    	//CLOSE_BTN
    	const close_btn = $(`<div class='plot_close_btn'>x</div>`);
		$(q).append(close_btn);
		$(close_btn).on('click', ()=>{
			$(close_btn).remove();
			$(q).find('.svg_plot').remove();
			$(q).find('.plots').remove();
		});
		//------------------------

    	//IMPORTANT!! - google API requires upper case letters for symbol
    	$(q).find('.plots').append(`<img class='google_plot_6m plot' src='https://www.google.com/finance/getchart?q=`+symbol_str.toUpperCase()+`&p=6M'></img>`);
    	$(q).find('.plots').append(`<img class='google_plot_5y plot' src='https://www.google.com/finance/getchart?q=`+symbol_str.toUpperCase()+`&p=5Y'></img>`);
    	$(q).find('.plots').append(`<img class='google_plot_20y plot' src='https://www.google.com/finance/getchart?q=`+symbol_str.toUpperCase()+`&p=20Y'></img>`);
    	
	    /*//-------------------------
		//$('body').append(`<script type="text/javascript" src="https://s3.tradingview.com/tv.js"></script>`);

		//const s = $(`<script type="text/javascript" src="https://s3.tradingview.com/tv.js"></script>`);
		//$(s).on('load',()=>{

			console.log('aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa')
			$('body').append(`
				<script type="text/javascript">
				new TradingView.widget({
				  "width": 980,
				  "height": 610,
				  "symbol": "NASDAQ:`+symbol_str+`",
				  "interval": "D",
				  "timezone": "Etc/UTC",
				  "theme": "Black",
				  "style": "1",
				  "locale": "en",
				  "toolbar_bg": "#f1f3f6",
				  "enable_publishing": false,
				  "allow_symbol_change": true,
				  "hideideas": true
				});
				</script>`);
		//});

		console.log('zzzzzzzzzzzzzzzzzzzzzzzzzzz')
		//console.log(s)
		//$('body').append(s);	
		//-------------------------*/

	    gf_trader.http__get_symbols_daily_historic([symbol_str],
			(p_quotes_lst)=>{

				//console.log(p_quotes_lst)

				const plot_id_str = symbol_str+'_plot';
				const plot        = $(`
					<div id='`+plot_id_str+`' class='svg_plot'>
						<svg width='800' height='600'></svg>
					</div>`);

				$(q).append(plot);

				const data_lst = [];
				for (var q_map of p_quotes_lst) {
					data_lst.push(q_map['close_price_f']);
				}

				data_lst.unshift(symbol_str+' close prices');

				const chart = c3.generate({
					bindto: '#'+plot_id_str,
					data: {
						columns: [
							data_lst,
							//['data1', 30, 200, 100, 400, 150, 250],
							//['data2', 50, 20, 10, 40, 15, 25]
						]
					}
				});				
			},
			()=>{},
			p_log_fun);
	});

    return q;
}