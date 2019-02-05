# GloFlowTrader asset trading, management, and research platform
# Copyright (C) 2019 Ivan Trajkovic
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; either version 2 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program; if not, write to the Free Software
# Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA

import os, sys
cwd_str = os.path.abspath(os.path.dirname(__file__))
#-------------------------------------------------------------
def get():

    meta_map = {
        'build_info_map':{
            #------------------------
            #GF_TRADER
            'gf_trader':{
                'version_str':        '0.8.0.0',
                'go_path_str':        '%s/../go/gf_trader'%(cwd_str),
                'go_output_path_str': '%s/../bin/gf_trader/gf_trader'%(cwd_str),
                'copy_to_dir_lst':    []
            },
            #-------------
        }
    }
    return meta_map