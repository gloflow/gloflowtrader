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

///<reference path="./../../d/jquery.d.ts" />

import * as gf_trader        from "./gf_trader";
import * as gf_trader_quotes from "./gf_trader_quotes";

//---------------------------------------------------
export function view(p_log_fun) {

	const container = $(`
		<div id='watchlist'>
			<div class='title'>watchlist</div>
			<div class='quotes'></div>
		</div>`);

	const watchlist_lst = [
		{'s':'nvda',  'title':'nvidia corporation'},
		{'s':'aapl',  'title':'apple'},
		{'s':'msft',  'title':'microsoft'},
		{'s':'googl', 'title':'Alphabet Inc.'},
		{'s':'amzn',  'title':'Amazon.com Inc'},
		{'s':'fb',    'title':'Facebook Inc'},
		{'s':'vod',   'title':'Vodafone Group Plc'},
		{'s':'ibm',   'title':'Internation Business Machines Corporation'},
		{'s':'amd',   'title':'Advanced Micro Devices'},
		{'s':'sne',   'title':'Sony Corporation'},
		{'s':'etrm',  'title':'EnteroMedics Inc'},
		{'s':'dis',   'title':'The Walt Desney Company'},
		{'s':'csco',  'title':'cisco systems'},
		{'s':'xom',   'title':'Exxon Mobile Coporation'},
		{'s':'sgnl',  'title':'Signal Genetics Inc'},
		{'s':'cur',   'title':'Neuralstem'},
		{'s':'qrvo',  'title':'Qorvo Inc'},
		{'s':'urre',  'title':'uranium resources'},
		{'s':'lmt',   'title':'lockheed martin corporation'},
		{'s':'ba',    'title':'boeing company'},
		{'s':'acn',   'title':'accenture plc'},
		{'s':'iots',  'title':'adesto technologies corporation'},
		{'s':'nok',   'title':'Nokia Corp'},
		{'s':'twtr',  'title':'Twitter Inc'},
		{'s':'znga',  'title':'Zynga Inc'},
		{'s':'txn',   'title':'Texas Instruments Inc.'},
		{'s':'lnkd',  'title':'LinkedIn Corp'},
		{'s':'etsy',  'title':'Etsy Inc'},
		{'s':'adbe',  'title':'Adobe Systems Incorporated'},
		{'s':'intc',  'title':'Intel Corporation'},
	];

	const watchlist_s_lst = watchlist_lst.map((w)=>{return w['s']});
	gf_trader.http__get_symbols_latest(watchlist_s_lst,
		(p_quotes_lst)=>{

			//-----------------------
			//IMPORTANT!! - sort stocks by change in dollar amount of price
			//              sort(()=>{}) - return -1, 0, or 1
			p_quotes_lst.sort((a, b)=>{
				return b['price_change_nominal_f'] - a['price_change_nominal_f'];
			});
			//-----------------------

			p_quotes_lst.forEach((q_map)=>{
				
				const quote = gf_trader_quotes.view_quote(q_map, p_log_fun);
				$(container).find('.quotes').append(quote);
			});
		},
		()=>{},
		p_log_fun);

	return container;
}