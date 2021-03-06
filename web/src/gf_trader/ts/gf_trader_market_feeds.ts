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

import * as gf_trader_plot from "./gf_trader_plot";

declare var EventSource;
//---------------------------------------------------
export function init(p_config_map, p_log_fun) {
    p_log_fun('FUN_ENTER', 'gf_trader_market_feeds.init()');

    //<div class='market_data'>

    const container = $(`
        <div>
            <div id='binance_feed' class='market_feed'>
                <div class='market_summary'>
                    <div class='title'>binance</div>
                    <div class='data_points'>
                        <div class='current_price'>
                            <span>price - </span><span>$</span><span class='price'>0.0</span>
                        </div>
                        
                        <div class='last_bid_price'>
                            <span>price bid - </span><span>$</span><span class='price'>0.0</span>
                        </div>
                        <div class='last_ask_price'>
                            <span>price ask - </span><span>$</span><span class='price'>0.0</span>
                        </div>

                        <div class='bid_trades_count'><span>bid/buy trades count - </span><span class='count'>0</span></div>
                        <div class='ask_trades_count'><span>ask/sell trades count - </span><span class='count'>0</span></div>
                    </div>
                </div>

                <div id='binance_market_plot'>

                    <!-- <div class='svg_plot' style='width: 1500px;height: 400px;'>
                        <svg width='1500' height='400'></svg>
                    </div> -->
                </div>

                
                    <div class='bid'>
                        <div class='title'>bid/buy</div>
                        <div class='trade__data'></div>
                        <div class='place_and_cancel__data'></div>
                    </div>

                    <div class='ask'>
                        <div class='title'>ask/sell</div>
                        <div class='trade__data'></div>
                        <div class='place_and_cancel__data'></div>
                    </div>
 
            </div>
            <div id='gemini_feed' class='market_feed'>
                <div class='market_summary'>
                    <div class='title'>gemini</div>
                    <div class='data_points'>
                        <div class='current_price'>
                            <span>price - </span><span>$</span><span class='price'>0.0</span>
                        </div>
                        
                        <div class='last_bid_price'>
                            <span>price bid - </span><span>$</span><span class='price'>0.0</span>
                        </div>
                        <div class='last_ask_price'>
                            <span>price ask - </span><span>$</span><span class='price'>0.0</span>
                        </div>

                        <div class='bid_trades_count'><span>bid/buy trades count - </span><span class='count'>0</span></div>
                        <div class='ask_trades_count'><span>ask/sell trades count - </span><span class='count'>0</span></div>
                    </div>
                </div>

                <div id='gemini_market_plot'>

                    <!-- <div class='svg_plot' style='width: 1500px;height: 400px;'>
                        <svg width='1500' height='400'></svg>
                    </div> -->
                </div>


                    <div class='bid'>
                        <div class='title'>bid/buy</div>
                        <div class='trade__data'></div>
                        <div class='place_and_cancel__data'></div>
                    </div>

                    <div class='ask'>
                        <div class='title'>ask/sell</div>
                        <div class='trade__data'></div>
                        <div class='place_and_cancel__data'></div>
                    </div>
           
            </div>
        </div>`);

    $('body').append(container);
    
    const binance__price_data_lst = [];
    const gemini__price_data_lst  = []; 

    gf_trader_plot.init_p5('binance_market_plot', binance__price_data_lst, p_config_map, p_log_fun);
    gf_trader_plot.init_p5('gemini_market_plot', gemini__price_data_lst, p_config_map, p_log_fun);

    init_updates("trader_binance_events", binance__price_data_lst, $(container).find(`#binance_feed`), p_config_map, p_log_fun);
    init_updates("trader_gemini_events", gemini__price_data_lst, $(container).find(`#gemini_feed`), p_config_map, p_log_fun);

    return container;
}
//---------------------------------------------------
export function init_updates(p_events_id_str, p_price_data_lst, p_container, p_config_map, p_log_fun) {
    p_log_fun('FUN_ENTER', 'gf_trader_market_feeds.init_updates()');

    //const initial_price_f               = 130.0;   //initial price that the plot will start at
    //const max_num_of_prices_to_show_int = 100;
    //const max_num_of_prices_to_plot_int = 10;

    console.log("REGISTER MARKET_FEED EVENT_SOURCE");
    const event_source  = new EventSource("/trader/events?events_id="+p_events_id_str)

    var   i = 0;
    const market_summary_map = {
        'last_price_f':         p_config_map['initial_price_f'], //initial_price_f,
        'last_side_str':        'bid', //'bid'|'ask'
        'bid_trades_count_int': 0,
        'ask_trades_count_int': 0
    };
    
    //const seconds_samples_num_int = 60*6; //number of seconds-resolution price datapoints
    //p_price_data_lst.push(initial_price_f);

    //IMPORTANT!! - this loop is happening on a regular interval, as oppose to event_source.onmessage
    //              which is irregular and depends on when the data comes from the server.
    setInterval(()=>{

        //IMPORTANT!! - if there is a certain number of prices
        //              remove the first price in order to stay in the range and be able
        //              to display new prices
        if (p_price_data_lst.length >= p_config_map['max_num_of_prices_to_plot_int']) { //max_num_of_prices_to_plot_int) {
            p_price_data_lst.shift();
        }

        const last_price_f = market_summary_map['last_price_f'];
        p_price_data_lst.push(last_price_f);

    }, 2000);

    event_source.onmessage = (p_e)=>{

        const event_data_map = JSON.parse(p_e.data);  
        //console.log(event_data_map)
            
        const meta_map   = event_data_map['meta_map'];
        const symbol_str = meta_map['e__symbol_str'];

        console.log(p_events_id_str+" - "+symbol_str)

        /*if (symbol_str != "ETHUSD") {
            return
        }*/

        view__update(event_data_map, market_summary_map);

        update_market_symmary(event_data_map, market_summary_map, p_log_fun);
        i+=1;
    }

    const bid__trade_element            = $(p_container).find('.bid .trade__data');
    const bid__place_and_cancel_element = $(p_container).find('.bid .place_and_cancel__data');

    const ask__trade_element            = $(p_container).find('.ask .trade__data');
    const ask__place_and_cancel_element = $(p_container).find('.ask .place_and_cancel__data');

    //---------------------------------------------------
    function init_visibility_onmouseover(p_element) {
        $(p_element).on('mouseover', ()=>{$(p_element).css({'overflow':'visible', 'z-index':1})});
        $(p_element).on('mouseout', ()=>{$(p_element).css({'overflow': 'hidden', 'z-index':0})});
    }
    //---------------------------------------------------
    init_visibility_onmouseover(bid__trade_element);
    init_visibility_onmouseover(bid__place_and_cancel_element);
    init_visibility_onmouseover(ask__trade_element);
    init_visibility_onmouseover(ask__place_and_cancel_element);


    //---------------------------------------------------
    function view__update(p_event_data_map, p_market_summary_map) {

    	const meta_map        = p_event_data_map['meta_map'];
        const symbol_str      = meta_map['e__symbol_str'];
        const price_f         = meta_map['e__price_f'];
        const price_rounded_f = Math.round(price_f * 1000) / 1000; 
        const side_str        = meta_map['e__side_str'];
        const remaining_f     = meta_map['e__remaining_f'];

        //console.log(' ----------     '+symbol_str)

        //-----------
        //REASON
        const reason_str = meta_map['e__reason_str'];

        var reason_class_str = 'reason';
        switch (reason_str) {

            case 'cancel':
                reason_class_str = 'reason__cancel';
                break;
            case 'place':
                reason_class_str = 'reason__place';
                break;
            case 'trade':
                reason_class_str = 'reason__trade';
                break;
        }
        //-----------

        var remaining__font_size_int = 10;
        var remaining__color_str     = '#bac2d4';
        if (remaining_f>1 && remaining_f<3) {
            remaining__font_size_int = 12;
            remaining__color_str     = '#e6e6e6';
        }
        else if (remaining_f>3 && remaining_f<10) {
            remaining__font_size_int = 14;
            remaining__color_str     = '#c3c3c3';
        }
        else if (remaining_f>10 && remaining_f<20) {
            remaining__font_size_int = 16;
            remaining__color_str     = '#909090';
        }
        else if (remaining_f>20) {
            remaining__font_size_int = 18;
            remaining__color_str     = '#7b7b7b';
        }

        const element = $(
            `<div class="price" id="${i}">`+
                `<span class="${reason_class_str}">${reason_str}</span>`+    
                `<span style="font-size:12px">$</span>${price_rounded_f}<span style="font-size:10px">eth</span>`+
                `<span style="font-size:${remaining__font_size_int}px;font-weight:bold;background-color:${remaining__color_str}">${remaining_f}</span>`+
            `</div>`);
        //---------------------------------------------------
        function remove_last_price_if_too_many(p_container) {
            const prices_num_int = $(p_container).find(`.price`).length;
            if (prices_num_int > p_config_map['max_num_of_prices_to_show_int']) { //max_num_of_prices_to_show_int) {
                $($(p_container).find(`.price`)[prices_num_int-1]).remove()
            }
        }
        //---------------------------------------------------
        switch (side_str) {
            case 'bid':
                if (reason_str == 'trade') {
                    $(bid__trade_element).prepend(element);
                    remove_last_price_if_too_many(bid__trade_element);
                }
                else if (reason_str == 'place' || reason_str == 'cancel') {
                    $(bid__place_and_cancel_element).prepend(element);
                    remove_last_price_if_too_many(bid__place_and_cancel_element);
                }
                break;

            case 'ask':
                if (reason_str == 'trade') {
                    $(ask__trade_element).prepend(element);
                    remove_last_price_if_too_many(ask__trade_element);
                }
                else if (reason_str == 'place' || reason_str == 'cancel') {
                    $(ask__place_and_cancel_element).prepend(element);
                    remove_last_price_if_too_many(ask__place_and_cancel_element);
                }
                break;
        }

        return price_rounded_f;
    }
    //---------------------------------------------------

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
function update_market_symmary(p_event_data_map, p_market_summary_map, p_log_fun) {

    const meta_map   = p_event_data_map['meta_map'];
    const symbol_str = meta_map['e__symbol_str'];
    const price_f    = meta_map['e__price_f'];
    const side_str   = meta_map['e__side_str'];
    const reason_str = meta_map['e__reason_str'];


    console.log(p_event_data_map)
    
    if (reason_str == 'trade') {

        $('.market_summary .current_price .price').text(price_f);

        switch (side_str) {
            case 'bid':
                $('.market_summary .last_bid_price .price').text(price_f);

                $('.market_summary .bid_trades_count .count').text(p_market_summary_map['bid_trades_count_int']+1);
                p_market_summary_map['bid_trades_count_int'] += 1;
                p_market_summary_map['last_side_str']         = 'bid';
                break;

            case 'ask':
                $('.market_summary .last_ask_price .price').text(price_f);

                $('.market_summary .ask_trades_count .count').text(p_market_summary_map['ask_trades_count_int']+1);
                p_market_summary_map['ask_trades_count_int'] += 1;
                p_market_summary_map['last_side_str']         = 'ask';
                break;
        }

        //IMPORTANT!! - record onlyt trade prices, not order place/cancel prices
        console.log(price_f);
        p_market_summary_map['last_price_f'] = price_f;
    }
}