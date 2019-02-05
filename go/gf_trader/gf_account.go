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
	"time"
	"fmt"
	"github.com/globalsign/mgo/bson"
)
//-------------------------------------------------
type Account struct {
	Id                     bson.ObjectId `bson:"_id,omitempty"          json:"-"`
	Id_str                 string        `bson:"id_str"                 json:"id_str"`
	T_str                  string        `bson:"t"                      json:"-"` //"account"
	Creation_unix_time_f   float64       `bson:"creation_unix_time_f"   json:"creation_unix_time_f"`
	Name_str               string        `bson:"name_str"               json:"name_str"`
	Starting_investment_f  float64       `bson:"starting_investment_f"  json:"starting_investment_f"`
	Comission_per_trade_f  float64       `bson:"comission_per_trade_f"  json:"comission_per_trade_f"`     //different accounts (brokers) have different comissions
	Current_balance_id_str string        `bson:"current_balance_id_str" json:"current_balance_id_str"` //id_str of Account__balance that is the latest one
}

//IMPORTANT!! - each account state update is persisted, so that account changes can be easily logged
//              and processed/monitored. Append-only nature of this data set, also allows for easy reverses
//              of account state
type Account__balance struct {
	Id                        bson.ObjectId `bson:"_id,omitempty"             json:"-"`
	Id_str                    string        `bson:"id_str"                    json:"id_str"`
	T_str                     string        `bson:"t"                         json:"-"` //"account balance"
	Account_name_str          string        `bson:"account_name_str"          json:"account_name_str"`
	Creation_unix_time_f      float64       `bson:"creation_unix_time_f"      json:"creation_unix_time_f"`
	Stocks_value_f            float64       `bson:"stocks_value_f"            json:"stocks_value_f"`       //value of stocks currently owned
	Available_funds_f         float64       `bson:"available_funds_f"         json:"available_funds_f"`    //money available for investment
	Total_value_f             float64       `bson:"total_value_f"             json:"total_value_f"`        //available cash plus stocks value
	Comissions_total_f        float64       `bson:"comissions_total_f"        json:"comissions_total_f"`   //amount of all comissions paid out
	Total_transactions_int    int           `bson:"total_transactions_int"    json:"total_transactions_int"`
	Positive_transactions_int int           `bson:"positive_transactions_int" json:"positive_transactions_int"`
	Negative_transactions_int int           `bson:"negative_transactions_int" json:"negative_transactions_int"`
}
//-------------------------------------------------
func account__create_defaults(p_runtime *Runtime) {
	p_runtime.Runtime_sys.Log_fun("FUN_ENTER","gf_account.account__create_defaults()")

	account__create("practice_trading",
				10000.0, //starting investment
				10.0,    //comission
				p_runtime)
}
//-------------------------------------------------
func account__create(p_name_str string,
				p_starting_investment_f float64,
				p_comission_per_trade_f float64,
				p_runtime               *Runtime) error {
	p_runtime.Runtime_sys.Log_fun("FUN_ENTER","gf_account.account__create()")

	c,err := p_runtime.Runtime_sys.Mongodb_coll.Find(bson.M{"t":"account","name_str":p_name_str}).Count()
	if err != nil {
		return err
	}

	//only create an account if one by this name doesnt already exist
	if c == 0 {

		//---------------
		//ACCOUNT_BALANCE - initial state
		creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
		balance_id_str       := "account_balance__"+fmt.Sprint(creation_unix_time_f)

		balance := &Account__balance{
			Id_str:                   balance_id_str,
			T_str:                    "account balance",
			Account_name_str:         p_name_str,
			Creation_unix_time_f:     creation_unix_time_f,
			Stocks_value_f:           0.0,                     //value of stocks currently owned
			Available_funds_f:        p_starting_investment_f, //money available for investment
			Total_value_f:            p_starting_investment_f, //available cash plus stocks value
			Comissions_total_f:       0.0,                     //amount of all comissions paid out
			Total_transactions_int:   0,
			Positive_transactions_int:0,
			Negative_transactions_int:0,
		}

		//DB
		err := p_runtime.Runtime_sys.Mongodb_coll.Insert(balance)
		if err != nil {
			return err
		}
		//---------------
		//ACCOUNT

		creation_unix_time_f = float64(time.Now().UnixNano())/1000000000.0
		id_str              := "account__"+fmt.Sprint(creation_unix_time_f)

		account := &Account{
			Id_str:                id_str,
			T_str:                 "account",
			Creation_unix_time_f:  creation_unix_time_f,
			Name_str:              p_name_str,
			Starting_investment_f: p_starting_investment_f,
			Comission_per_trade_f: p_comission_per_trade_f,
			Current_balance_id_str:balance.Id_str,
		}

		//DB
		err = p_runtime.Runtime_sys.Mongodb_coll.Insert(account)
		if err != nil {
			return err
		}
		//---------------
	}

	return nil
}
//-------------------------------------------------
func account__get(p_account_name_str string, p_runtime *Runtime) (*Account,error) {
	p_runtime.Runtime_sys.Log_fun("FUN_ENTER","gf_account.account__get()")


	var account *Account
	err := p_runtime.Runtime_sys.Mongodb_coll.Find(bson.M{
						"t":       "account",
						"name_str":p_account_name_str,
					}).One(account)
	if err != nil {
		return nil,err
	}

	return account,nil
}
//-------------------------------------------------
func account__update(p_transaction *Transaction,
				p_account *Account,
				p_runtime *Runtime) error {
	p_runtime.Runtime_sys.Log_fun("FUN_ENTER","gf_account.account__update()")




	var account_balance *Account__balance
	err := p_runtime.Runtime_sys.Mongodb_coll.Find(bson.M{
							"t"     :"account_balance",
							"id_str":p_account.Current_balance_id_str,
						}).One(account_balance)
	if err != nil {
		return err
	}

	


	new_stocks_value_f  := float64(p_transaction.Shares_num_int) * p_transaction.Share_cost_f
	total_stock_value_f := account_balance.Stocks_value_f + new_stocks_value_f



	new_available_funds_f := account_balance.Available_funds_f - new_stocks_value_f


	//FIX!! - this is a varying figure, and changes with the change in the quote/stock price.
	//        so that value, when used by other displays/calculations should recalculated each time, 
	//        since the quotes will change continuously.
	new_total_value_f          := account_balance.Total_value_f + new_stocks_value_f
	new_total_transactions_int := account_balance.Total_transactions_int
	
	creation_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	balance_id_str       := "account_balance__"+fmt.Sprint(creation_unix_time_f)

	new_account_balance := &Account__balance{
		Id_str:                   balance_id_str,
		T_str:                    "account balance",
		Account_name_str:         p_account.Name_str,
		Creation_unix_time_f:     creation_unix_time_f,
		Stocks_value_f:           total_stock_value_f,   //value of stocks currently owned
		Available_funds_f:        new_available_funds_f, //money available for investment
		Total_value_f:            new_total_value_f,     //available cash plus stocks value
		Comissions_total_f:       0.0,                   //amount of all comissions paid out
		Total_transactions_int:   new_total_transactions_int,
		Positive_transactions_int:0,
		Negative_transactions_int:0,
	}
	fmt.Println(new_account_balance)


	return nil
}
//-------------------------------------------------
func account__get_available_funds(p_account_name_str string, p_runtime *Runtime) (*float64,error) {
	p_runtime.Runtime_sys.Log_fun("FUN_ENTER","gf_account.account__get_available_funds()")

	var account Account
	err := p_runtime.Runtime_sys.Mongodb_coll.Find(bson.M{
			"t"       :"account",
			"name_str":p_account_name_str,
		}).One(&account)
	if err != nil {
		return nil,err
	}

	return nil,nil //&account.Available_funds_f,nil
}