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

import "./gf_trader_watchlist";
import "./gf_trader_gemini";

namespace gf_trader {

$(document).ready(()=>{
    //-------------------------------------------------
    function log_fun(p_g,p_m) {
        var msg_str = p_g+':'+p_m
        //chrome.extension.getBackgroundPage().console.log(msg_str);

        switch (p_g) {
            case "INFO":
                console.log("%cINFO"+":"+"%c"+p_m,"color:green; background-color:#ACCFAC;","background-color:#ACCFAC;");
                break;
            case "FUN_ENTER":
                console.log("%cFUN_ENTER"+":"+"%c"+p_m,"color:yellow; background-color:lightgray","background-color:lightgray");
                break;
        }
    }
    //-------------------------------------------------
    gf_trader.init(log_fun);
});
//---------------------------------------------------
export function init(p_log_fun) {

    var stocks__open_bool = false;
    var crypto__open_bool = false;

    var stocks_container;
    var crypto_container;
    //----------------
    //STOCKS_WATCHLIST
    
    $('#markets #market__stocks_btn').on('click', ()=>{

        if (stocks__open_bool == false) {
            
            if (crypto__open_bool) {
                $(crypto_container).remove();
                crypto__open_bool = false;
            }

            stocks_container = gf_trader_watchlist.view(p_log_fun)
            $('body').append(stocks_container);
            stocks__open_bool = true;
        }
    });
    //----------------

    //----------------
    //CRYPTO_CURRENCY
    
    $('#markets #market__crypto_btn').on('click', ()=>{

        $(stocks_container).remove();
        if (crypto__open_bool == false) {

            if (stocks__open_bool) {
                $(stocks_container).remove();
                stocks__open_bool = false;
            }

            crypto_container = gf_trader_gemini.init(p_log_fun);
            
            crypto__open_bool = true;
        }
    });
    //----------------

    gf_trader_transactions.init__import(p_log_fun);
}
//---------------------------------------------------
export function http__get_symbols_daily_historic(p_symbols_lst :string[], p_onComplete_fun, p_onError_fun, p_log_fun) {

    const url_str = '/trader/quotes/daily_historic?symbols='+p_symbols_lst.join();
    p_log_fun('INFO','url_str - '+url_str);

    //-------------------------
    //HTTP AJAX
    $.get(url_str,
        function(p_data) {
            console.log('response received');
            //console.log('p_data - '+p_data);
            const data_map = JSON.parse(p_data);

            console.log('data_map["status_str"] - '+data_map["status_str"]);
            
            if (data_map["status_str"] == 'OK') {
                const quotes_lst = data_map['quotes_lst'];
                p_onComplete_fun(quotes_lst);
            }
            else {
                p_onError_fun(data_map["data"]);
            }
        });
    //------------------------- 
}
//---------------------------------------------------
export function http__get_symbols_latest(p_symbols_lst :string[], p_onComplete_fun, p_onError_fun, p_log_fun) {

    const url_str = '/trader/quotes/latest?symbols='+p_symbols_lst.join();
    p_log_fun('INFO','url_str - '+url_str);

    //-------------------------
    //HTTP AJAX
    $.get(url_str,
        function(p_data) {
            console.log('response received');
            //console.log('p_data - '+p_data);
            const data_map = JSON.parse(p_data);

            console.log('data_map["status_str"] - '+data_map["status_str"]);
            
            if (data_map["status_str"] == 'OK') {
                const quotes_lst = data_map['quotes_lst'];
                p_onComplete_fun(quotes_lst);
            }
            else {
                p_onError_fun(data_map["data"]);
            }
        });
    //------------------------- 
}
//---------------------------------------------------
}