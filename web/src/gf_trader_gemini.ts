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

///<reference path="./d/jquery.d.ts" />

namespace gf_trader_gemini {

declare var EventSource;
declare var p5;
//---------------------------------------------------
export function init(p_log_fun) {
    p_log_fun('FUN_ENTER','gf_trader_gemini.init()');

    const container = $(`
        <div id='gemini'>
            <div class='title'>gemini market</div>

            <div class='market_summary'>
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

            <div id='market_plot'>

                <!-- <div class='svg_plot' style='width: 1500px;height: 400px;'>
                    <svg width='1500' height='400'></svg>
                </div> -->
            </div>

            <div class='market_data'>
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
    

    const price_data_lst = [];

    init_p5(price_data_lst,
        p_log_fun);

    
    init_updates(price_data_lst,
            container,
            p_log_fun);

    return container;
}
//---------------------------------------------------
function draw_y_axis(p_max_price_f
                p_min_price_f,
                p_canvas_height_int,
                p_p5) {

    //------------
    //AXIS
    p_p5.fill(0);

    p_p5.stroke(0);
    p_p5.line(2,0,
        2,p_canvas_height_int);

    p_p5.rect(2,0,4,4);                     //max_price (p5 coord_sys origin is upper_left)
    p_p5.rect(2,p_canvas_height_int-4,4,4); //min_price
    //------------
    //TEXT
    p_p5.stroke(80);
    const text_size_int = 12;
    p_p5.textSize(text_size_int);

    //round dollar amounts, since only 2 decimal places make sense (100 cents in a dollar)
    const max_price__rounded_f = p_max_price_f.toFixed(2);
    const min_price__rounded_f = p_min_price_f.toFixed(2);

    p_p5.text("$"+max_price__rounded_f, 6, 0+4+text_size_int);
    p_p5.text("$"+min_price__rounded_f, 6, p_canvas_height_int-4-4);
}
//---------------------------------------------------
function draw_grid(p_max_price_f,
            p_min_price_f,
            p_canvas_height_int,
            p_canvas_width_int,
            p_p5) {





    p_p5.stroke(120);

    const price_range_f = p_max_price_f - p_min_price_f;

    var i=0;
    while (true) {


        if (p_min_price_f+i > p_max_price_f) {
            break;
        }

        const price_delta_f = Math.ceil(p_min_price_f+i) - p_min_price_f;
        const y_f        = (p_canvas_height_int * price_delta_f) / price_range_f;
        const y_inverted = p_canvas_height_int - y_f; //IMPORTANT!! - invert because p5 coord system is upper-left corner




        p_p5.line(0,y_inverted,
            p_canvas_width_int,y_inverted);


        i++;

    }
}
//---------------------------------------------------
function init_p5(p_price_data_lst,
            p_log_fun) {
    p_log_fun('FUN_ENTER','gf_trader_gemini.init_p5()');
        
    const canvas_width_int  = 1500;
    const canvas_height_int = 400;
    const price_point_dim_f = 2.0;
    const x_offset_f        = 10;

    const p5_env = function(p5) {

        p5.setup = function() {
            p5.createCanvas(canvas_width_int, canvas_height_int);
        };

        p5.draw = function() {
            p5.background(100);
            p5.fill(255);

            

            //IMPORTANT!! - plot is adjusted to only display the max/min price range. 
            //              so that the min_price_f is at Y=0, and max_price_f is at Y=canvas_height_int
            const max_price_f   = Math.max.apply(null,p_price_data_lst);
            const min_price_f   = Math.min.apply(null,p_price_data_lst);
            const price_range_f = max_price_f - min_price_f;


            draw_y_axis(max_price_f,
                min_price_f,
                canvas_height_int,
                p5);



            draw_grid(max_price_f,
                min_price_f,
                canvas_height_int,
                canvas_width_int,
                p5);


            const prev_x_f        = null;
            const prev_y_inverted = null;

            for (var i=0;i<p_price_data_lst.length;i++) {

                //-------------------
                const price_f    = p_price_data_lst[i];
                const x_f        = x_offset_f+i*3;

                //IMPORTANT!! - difference of the current price from the minimum price
                //              in the price_range_f
                const price_delta_f = price_f - min_price_f;
                //-------------------
                //IMPORTANT!! - convert price_delta_f to Y in pixels, so that it can 
                //              be plotted.
                //
                //1. canvas_height_int : price_range_f = y : price_delta_f
                //2. canvas_height_int * price_f = price_range_f * y
                //3. y = (canvas_height_int * price_f) / price_range_f
                
                const y_f        = (canvas_height_int * price_delta_f) / price_range_f;
                const y_inverted = canvas_height_int - y_f; //IMPORTANT!! - invert because p5 coord system is upper-left corner
                //-------------------






                //console.log('price_f - '+price_f)
                //console.log('>>>> -- '+price_range_f+' - '+y_f)

                p5.stroke(255);
                p5.rect(x_f-(price_point_dim_f/2),y_inverted-(price_point_dim_f/2),price_point_dim_f,price_point_dim_f);

                p5.stroke(153);
                if (prev_x_f != null) {
                    p5.line(prev_x_f,prev_y_inverted,
                        x_f,y_inverted);
                }

                prev_x_f        = x_f;
                prev_y_inverted = y_inverted;
            }
        };
    };

    const plot_containing_div_id_str = 'market_plot';
    const custom_p5                  = new p5(p5_env,plot_containing_div_id_str);
}
//---------------------------------------------------
export function init_updates(p_price_data_lst,
                        p_container,
                        p_log_fun) {
    p_log_fun('FUN_ENTER','gf_trader_gemini.init_updates()');

    console.log("REGISTER GEMINI EVENT_SOURCE");
    const events_id_str = "trader_gemini_events";
    const event_source  = new EventSource("/trader/events?events_id="+events_id_str)

    const initial_price_f = 760.5;   //initial price that the plot will start at
    
    var   i = 0;
    const market_summary_map = {
        'last_price_f'        :initial_price_f,
        'last_side_str'       :'bid', //'bid'|'ask'
        'bid_trades_count_int':0,
        'ask_trades_count_int':0
    };
    
    //const seconds_samples_num_int = 60*6; //number of seconds-resolution price datapoints
    p_price_data_lst.push(initial_price_f);

    setInterval(function() {

        //IMPORTANT!! - if there is a certain number of prices
        //              remove the first price in order to stay in the range and be able
        //              to display new prices
        if (p_price_data_lst.length >= 500) {
            p_price_data_lst.shift();
        }

        p_price_data_lst.push(market_summary_map['last_price_f']);

    },2000);

    event_source.onmessage = (p_e)=>{

        const event_data_map = JSON.parse(p_e.data);  
        console.log(event_data_map)
            
        const meta_map   = event_data_map['meta_map'];
        const symbol_str = meta_map['e__symbol_str'];

        if (symbol_str != "ETHUSD") {
            return
        }

        view__update(event_data_map,
                market_summary_map);

        update_market_symmary(event_data_map,
                        market_summary_map,
                        p_log_fun);
        i+=1;
    }

    const bid__trade_element            = $(p_container).find('.bid .trade__data');
    const bid__place_and_cancel_element = $(p_container).find('.bid .place_and_cancel__data');

    const ask__trade_element            = $(p_container).find('.ask .trade__data');
    const ask__place_and_cancel_element = $(p_container).find('.ask .place_and_cancel__data');

    //---------------------------------------------------
    function view__update(p_event_data_map,
                    p_market_summary_map) {

    	const meta_map    = p_event_data_map['meta_map'];
        const symbol_str  = meta_map['e__symbol_str'];
    	const price_f     = meta_map['e__price_f'];
        const side_str    = meta_map['e__side_str'];
        const remaining_f = meta_map['e__remaining_f'];

        console.log(' ----------     '+symbol_str)

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
            `<div class="price">`+
                `<span style="font-size:10px">$</span>`+
                price_f+' '+
                `<span style="font-size:10px">eth</span>`+
                `<span style="font-size:`+remaining__font_size_int+`px;font-weight:bold;background-color:`+remaining__color_str+`">`+remaining_f+`</span>`+
                `<span class="`+reason_class_str+`">`+reason_str+`</span>`+
            `</div>`);

        switch (side_str) {
            case 'bid':
                if (reason_str == 'trade') {
                    $(bid__trade_element).prepend(element);
                }
                else if (reason_str == 'place' || reason_str == 'cancel') {
                    $(bid__place_and_cancel_element).prepend(element);
                }
                break;

            case 'ask':
                if (reason_str == 'trade') {
                    $(ask__trade_element).prepend(element);
                }
                else if (reason_str == 'place' || reason_str == 'cancel') {
                    $(ask__place_and_cancel_element).prepend(element);
                }
                break;
        }

        return price_f;
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
function update_market_symmary(p_event_data_map,
                        p_market_summary_map,
                        p_log_fun) {

    const meta_map    = p_event_data_map['meta_map'];
    const symbol_str  = meta_map['e__symbol_str'];
    const price_f     = meta_map['e__price_f'];
    const side_str    = meta_map['e__side_str'];

    if (symbol_str != "ETHUSD") {
        return
    }

    const reason_str = meta_map['e__reason_str'];
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

        p_market_summary_map['last_price_f'] = price_f;
    }
}
//---------------------------------------------------
}