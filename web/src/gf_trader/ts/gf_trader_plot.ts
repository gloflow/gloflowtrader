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

declare var p5;
//---------------------------------------------------
export function init_p5(p_canvas_parent_div_id_str, p_price_data_lst, p_log_fun) {
    p_log_fun('FUN_ENTER', 'gf_trader_plot.init_p5()');
        
    const canvas_width_int  = $(window).width() * 0.8; //80% of window width
    const canvas_height_int = 400;
    const price_point_dim_f = 2.0;
    const x_offset_f        = 10;


    $(`#${p_canvas_parent_div_id_str}`).css("width", `${canvas_width_int}px`);            //make parent div the same width as canvas
    $(`#${p_canvas_parent_div_id_str}`).parent().css("width", `${canvas_width_int}px`);   //make parent/parent div the same width/height as canvas
    $(`#${p_canvas_parent_div_id_str}`).parent().css("height", `${canvas_height_int}px`);

    const p5_env = function(p5) {

        p5.setup = function() {
            p5.createCanvas(canvas_width_int, canvas_height_int);
        };

        p5.draw = function() {
            p5.background(100);
            p5.fill(255);

            //IMPORTANT!! - plot is adjusted to only display the max/min price range. 
            //              so that the min_price_f is at Y=0, and max_price_f is at Y=canvas_height_int
            const max_price_f   = Math.max.apply(null, p_price_data_lst);
            const min_price_f   = Math.min.apply(null, p_price_data_lst);
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

            var prev_x_f        = null;
            var prev_y_inverted = null;

            for (var i=0; i < p_price_data_lst.length; i++) {

                //-------------------
                const price_f = p_price_data_lst[i];
                const x_f     = x_offset_f+i*3;

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

                p5.stroke(255);
                p5.rect(x_f-(price_point_dim_f/2), y_inverted-(price_point_dim_f/2), price_point_dim_f, price_point_dim_f);

                p5.stroke(153);
                if (prev_x_f != null) {
                    p5.line(prev_x_f, prev_y_inverted, x_f, y_inverted);
                }

                prev_x_f        = x_f;
                prev_y_inverted = y_inverted;
            }
        };
    };

    const custom_p5 = new p5(p5_env, p_canvas_parent_div_id_str);
}
//---------------------------------------------------
function draw_y_axis(p_max_price_f, p_min_price_f, p_canvas_height_int, p_p5) {

    //------------
    //AXIS
    p_p5.fill(0);

    p_p5.stroke(0);
    p_p5.line(2, 0, 2, p_canvas_height_int);

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

        p_p5.line(0, y_inverted, p_canvas_width_int, y_inverted);

        i++;
    }
}