import os
cwd_str = os.path.abspath(os.path.dirname(__file__))
#-------------------------------------------------------------
def get():

	apps_map = {
		#-----------------------------
		'gf_trader':{
			'pages_map':{
				#-------------
				#IMAGES_FLOWS_BROWSER
				'gf_trader':{
					'build_dir_str':      '%s/../web/build'%(cwd_str),
					'main_html_path_str': '%s/../web/src/gf_trader/templates/gf_trader.html'%(cwd_str),
					'url_base_str':       '/trader/static',

					# 'type_str':      'ts',
					# 'build_dir_str': '%s/../web/build/gf_apps/gf_images'%(cwd_str),
					# 'ts':{
					# 	'out_file_str':      '%s/../web/build/gf_apps/gf_images/js/gf_images_flows_browser.js'%(cwd_str),
					# 	'minified_file_str': '%s/../web/build/gf_apps/gf_images/js/gf_images_flows_browser.min.js'%(cwd_str),
					# 	'files_lst':[
					# 		'%s/../web/src/gf_apps/gf_images/ts/flows_browser/gf_images_flows_browser.ts'%(cwd_str),
					# 		'%s/../web/src/gf_core/ts/gf_gifs.ts'%(cwd_str),
					# 		'%s/../web/src/gf_core/ts/gf_gifs_viewer.ts'%(cwd_str),
					# 		'%s/../web/src/gf_core/ts/gf_image_viewer.ts'%(cwd_str),
					# 		'%s/../web/src/gf_core/ts/gf_sys_panel.ts'%(cwd_str),
					# 	],
					# 	#-------------
					# 	#LIBS
					# 	'libs_files_lst':[
					# 		'%s/../web/libs/js/masonry.pkgd.min.js'%(cwd_str),
					# 		'%s/../web/libs/js/jquery.timeago.js'%(cwd_str),
					# 	]
					# 	#-------------
					# },
					# 'css':{
					# 	'files_lst':[
					# 		('%s/../web/src/gf_apps/gf_images/css/gf_images_flows_browser.css'%(cwd_str), '%s/../web/build/gf_apps/gf_images/css/flows_browser'%(cwd_str)),
					# 		('%s/../web/src/gf_core/css/gf_gifs_viewer.css'%(cwd_str),                    '%s/../web/build/gf_apps/gf_images/css/flows_browser'%(cwd_str)),
					# 		('%s/../web/src/gf_core/css/gf_image_viewer.css'%(cwd_str),                   '%s/../web/build/gf_apps/gf_images/css/flows_browser'%(cwd_str)),
					# 		('%s/../web/src/gf_core/css/gf_sys_panel.css'%(cwd_str),                      '%s/../web/build/gf_apps/gf_images/css/flows_browser'%(cwd_str)),
					# 	]
					# },
					# 'templates':{
					# 	'files_lst':[
					# 		('%s/../web/src/gf_apps/gf_images/templates/flows_browser/gf_images_flows_browser.html'%(cwd_str), '%s/../web/build/gf_apps/gf_images/templates/flows_browser'%(cwd_str)),
					# 	]
					# }
				},
				#-------------
			}
		}
	}
	return apps_map