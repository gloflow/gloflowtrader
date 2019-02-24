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

//---------------------------------------------------
export function view__buy_dialog(p_stock_symbol_str, p_stock_price_f, p_log_fun) {

	const b = $(`
		<div id='buy'>
			<div id='close_btn'>x</div>
			<div id='buy_quantity'>
				<span class='label'>buying quantity</span> <input id='number' type='number'></input>
			</div>

			<div id='available_funds'>
				<span class='label'>available funds</span>    $<span class='amount'></span>
			</div>
			<div id='buying_capacity'>
				<span class='label'>buying capacity</span>  <span class='number'></span>
			</div>
			<div id='debit'>
				<div id='stock_debit'>
					<span class='label'>stock debit</span>    $<span class='amount'></span>
				</div>
				<div id='comission'>
					<span class='label'>comission</span>    $<span class='amount'></span>
				</div>
				<div id='total_debit'>
					<span class='label'>total debit</span>    $<span class='amount'></span>
				</div>
			</div>
			<div id='execute_btn'>execute</div>
		</div>`);

	$(b).find('#close_btn').on('click', ()=>{
		$(b).remove();
	})
	
	$(b).find('input').keyup((e)=>{
		const stocks_number_int = $(b).find('input').val();
		const stock_debit_f     = stocks_number_int * p_stock_price_f;
		const comission_f       = 10.0;
		const total_debit_f     = stock_debit_f + comission_f;


		$(b).find('#stock_debit .amount').text(stock_debit_f.toFixed(2)); //2 decimal places
		$(b).find('#comission .amount').text(comission_f.toFixed(2));
		$(b).find('#total_debit .amount').text(total_debit_f.toFixed(2));
	});

	$('#execute_btn').on('click',()=>{

	});

	return b;
}
//---------------------------------------------------
export function init__import(p_log_fun) {

	$('#transactions #import_btn').on('click', ()=>{

		const import_dialog = $(`
			<div id='import_dialog'>
				<div id='close_btn'>x</div>
				<div><span>symbol</span><input id='symbol_input'></input></div>
				<div><span>date</span><input id='date_input' type="date"></input>

				<div id='comission'>
					<span>comission $</span>
					<select id='comission_select' type='number'>
						<option value='10'>10</option>
						<option value='5'>5</option>
					</select>
				</div>

				
				<div><span>shares #</span><input id='shares_num' type='number'></input></div>
				<div><span>shares cost</span> $<input id='shares_cost' type='number'></input></div>
				
				<div>
					<span>type</span>
					<select id='type_select'>
						<option value='buy'>buy</option>
						<option value='sell'>sell</option>
					</select>
				</div>

				<div>
					<span>executor type</span>
					<select id='executor_type_select'>
						<option value='stocktrainer'>stocktrainer</option>
						<option value='etrade'>etrade</option>
						<option value='robinhood'>robinhood</option>
					</select>
				</div>

				<div id='create_btn'>create</div>
			</div>`);
		$('#transactions').append(import_dialog);

		$(import_dialog).find('#close_btn').on('click', ()=>{
			$(import_dialog).remove();
		});

		$(import_dialog).find('#create_btn').on('click', ()=>{

			const transaction_map = {
				'symbol_str':        $(import_dialog).find('#symbol_input').val(),
				'date':              $(import_dialog).find('#date_input').val(),
				'comission_f':       $(import_dialog).find('#comission_select option:selected').text(),
				'shares_num_int':    $(import_dialog).find('#shares_num').val(),
				'share_cost_f':      $(import_dialog).find('#shares_cost').val(),
				'type_str':          $(import_dialog).find('#type_select option:selected').text(),
				'executor_type_str': $(import_dialog).find('#executor_type_select option:selected').text(),
				'origin_type_str':   'manual_import'
			};
			http__transaction_import(transaction_map, ()=>{}, ()=>{}, p_log_fun);
		});
	});
}
//---------------------------------------------------
export function http__transaction_execute(p_on_complete_fun, p_on_error_fun, p_log_fun) {

	const url_str = '/trader/transaction/import';
    p_log_fun('INFO','url_str - '+url_str);   
}
//---------------------------------------------------
export function http__transaction_import(p_transaction_map :Object, p_on_complete_fun, p_on_error_fun, p_log_fun) {

    const url_str = '/trader/transaction/import';
    p_log_fun('INFO','url_str - '+url_str);

    //-------------------------
    //HTTP AJAX
    $.post(url_str,
    	p_transaction_map,
        function(p_data) {
            console.log('response received');
            //console.log('p_data - '+p_data);
            const data_map = JSON.parse(p_data);

            console.log('data_map["status_str"] - '+data_map["status_str"]);
            
            if (data_map["status_str"] == 'OK') {
                const quotes_lst = data_map['quotes_lst'];
                p_on_complete_fun(quotes_lst);
            }
            else {
                p_on_error_fun(data_map["data"]);
            }
        });
    //------------------------- 
}