'use strict';

var dashboard = {
	title: 'Karhu',
	time: {
		from: 'now-24h',
		to: 'Now'
	},
	rows: []
};

var host = -1;
if (ARGS.host !== undefined) {
	host = ARGS.host;
}

dashboard.rows.push({
	title: 'creme',
	height: '300px',
	panels: get_panels()
});


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
				query: 'collectd_type="cpu"'
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
				query: 'collectd_type="memory"ANDtype_instance:"used"'
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
				query: 'collectd_type="swap"'
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
				query: 'collectd_type="disk_merged"'
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
				query: 'collectd_type="disk_octets"'
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
				query: 'collectd_type="disk_ops"'
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
				query: 'collectd_type="disk_time"'
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

console.log('DASHBOARD > ', dashboard);
return dashboard;