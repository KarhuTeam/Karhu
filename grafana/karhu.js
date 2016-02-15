'use strict';

var dashboard = {
	title: 'Karhu',
	time: {
		from: 'now-24h',
		to: 'Now'
	},
	rows: []
};

if (ARGS.hosts !== undefined) {
	var hosts = ARGS.hosts.split(",")

	for (var i = 0 ; i < hosts.length ; i++) {
		dashboard.rows.push({
			title: hosts[i],
			height: '300px',
			panels: get_panels(hosts[i])
		});
	}

} else {

	dashboard.rows.push({
		title: 'Global',
		height: '300px',
		panels: get_panels()
	});

}

function get_panels() {
	var panels = [];

	panels.push(build_panel('cpu'));
	panels.push(build_panel('memory'));
	panels.push(build_panel('disk'));
	panels.push(build_panel('others'));

	return panels;
}

function build_panel(v) {
	var panel = {
		title: v,
		type: 'graph',
		span: 6,
		fill: 1,
		linewidth: 2,
		targets: build_panel_targets(v)
	};
	return panel;
}

function build_panel_targets(v) {
	switch (v) {
		case 'cpu':
			return [{
				refId: 'A',
				metrics: [{
					type: 'avg',
					id: '1',
					field: 'value',
				}],
				dsType: "elasticsearch",
				bucketAggs: [{
					type: "date_histogram",
					id: "2",
					settings: {
						interval: "auto",
						min_doc_count: 0
					},
					field: "@timestamp"
				}],
				timeField: "@timestamp",
				query: 'collectd_type="cpu"',
				datasource: 'es'
			}];

		case 'memory':
			return [{
				refId: 'A',
				metrics: [{
					type: 'avg',
					id: '1',
					field: 'value',
				}],
				dsType: "elasticsearch",
				bucketAggs: [{
					type: "date_histogram",
					id: "2",
					settings: {
						interval: "auto",
						min_doc_count: 0
					},
					field: "@timestamp"
				}],
				timeField: "@timestamp",
				query: 'collectd_type="memory"ANDtype_instance:"used"',
				datasource: 'es'
			}, {
				refId: 'A',
				metrics: [{
					type: 'avg',
					id: '1',
					field: 'value',
				}],
				dsType: "elasticsearch",
				bucketAggs: [{
					type: "date_histogram",
					id: "2",
					settings: {
						interval: "auto",
						min_doc_count: 0
					},
					field: "@timestamp"
				}],
				timeField: "@timestamp",
				query: 'collectd_type="swap"',
				datasource: 'es'
			}];

		case 'disk':
			return [{
				refId: 'A',
				metrics: [{
					type: 'avg',
					id: '1',
					field: 'value',
				}],
				dsType: "elasticsearch",
				bucketAggs: [{
					type: "date_histogram",
					id: "2",
					settings: {
						interval: "auto",
						min_doc_count: 0
					},
					field: "@timestamp"
				}],
				timeField: "@timestamp",
				query: 'collectd_type="disk_merged"',
				datasource: 'es'
			}, {
				refId: 'A',
				metrics: [{
					type: 'avg',
					id: '1',
					field: 'value',
				}],
				dsType: "elasticsearch",
				bucketAggs: [{
					type: "date_histogram",
					id: "2",
					settings: {
						interval: "auto",
						min_doc_count: 0
					},
					field: "@timestamp"
				}],
				timeField: "@timestamp",
				query: 'collectd_type="disk_octets"',
				datasource: 'es'
			}, {
				refId: 'A',
				metrics: [{
					type: 'avg',
					id: '1',
					field: 'value',
				}],
				dsType: "elasticsearch",
				bucketAggs: [{
					type: "date_histogram",
					id: "2",
					settings: {
						interval: "auto",
						min_doc_count: 0
					},
					field: "@timestamp"
				}],
				timeField: "@timestamp",
				query: 'collectd_type="disk_ops"',
				datasource: 'es'
			}, {
				refId: 'A',
				metrics: [{
					type: 'avg',
					id: '1',
					field: 'value',
				}],
				dsType: "elasticsearch",
				bucketAggs: [{
					type: "date_histogram",
					id: "2",
					settings: {
						interval: "auto",
						min_doc_count: 0
					},
					field: "@timestamp"
				}],
				timeField: "@timestamp",
				query: 'collectd_type="disk_time"',
				datasource: 'es'
			}];

		case 'others':
			return [{
				refId: 'A',
				metrics: [{
					type: 'avg',
					id: '1',
					field: 'value',
				}],
				dsType: "elasticsearch",
				bucketAggs: [{
					type: "date_histogram",
					id: "2",
					settings: {
						interval: "auto",
						min_doc_count: 0
					},
					field: "@timestamp"
				}],
				timeField: "@timestamp",
				query: 'collectd_type="irq"'
			}];
	}
}

return dashboard;