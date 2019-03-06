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

export function get() {

    const initial_price_f               = 130.0;   //initial price that the plot will start at
    const max_num_of_prices_to_show_int = 100;     //number of items to display in the bid/ask order feed
    const max_num_of_prices_to_plot_int = 50;

    const plot__canvas_width_int  = 800; //$(window).width() * 0.8; //80% of window width
    const plot__canvas_height_int = 400;
    const plot__price_point_dim_f = 2.0;
    const plot__x_offset_f        = 10;
    const plot__x_delta_pixels_f  = plot__canvas_width_int / max_num_of_prices_to_plot_int;

    return {
        "initial_price_f":               initial_price_f, 
        "max_num_of_prices_to_show_int": max_num_of_prices_to_show_int,
        "max_num_of_prices_to_plot_int": max_num_of_prices_to_plot_int,

        "plot__canvas_width_int":  plot__canvas_width_int, 
        "plot__canvas_height_int": plot__canvas_height_int, 
        "plot__price_point_dim_f": plot__price_point_dim_f, 
        "plot__x_offset_f":        plot__x_offset_f,
        "plot__x_delta_pixels_f":  plot__x_delta_pixels_f
    }
}