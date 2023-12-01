import { Span } from 'types/api/trace/getTraceItem';

export const TraceData: Span[] = [
	[
		1657275433246,
		'span_1',
		'0000000000000000span_1',
		'frontend',
		'HTTP GET SPAN 1',
		'2',
		'683273000',
		['host.name.span1'],
		['span_1'],
		['{TraceId=0000000000000000span_1, SpanId=, RefType=CHILD_OF}'],
		[
			'{"timeUnixNano":1657275433246142000,"attributeMap":{"event":"HTTP request received S1","level":"info","method":"GET","url":"/dispatch?customer=392\\u0026nonse=0.015296363321630757"}}',
		],
		false,
	],
	[
		1657275433246,
		'span_2',
		'0000000000000000span_1',
		'frontend',
		'HTTP GET SPAN 2',
		'2',
		'683273000',
		['host.name.span2'],
		['span_2'],
		['{TraceId=0000000000000000span_1, SpanId=span_1, RefType=CHILD_OF}'],
		[
			'{"timeUnixNano":1657275433246142000,"attributeMap":{"event":"HTTP request received S2","level":"info","method":"GET","url":"/dispatch?customer=392\\u0026nonse=0.015296363321630757"}}',
		],
		false,
	],
	[
		1657275433246,
		'span_3',
		'0000000000000000span_1',
		'frontend',
		'HTTP GET SPAN 3',
		'2',
		'683273000',
		['host.name.span3'],
		['span_3'],
		['{TraceId=0000000000000000span_1, SpanId=span_2, RefType=CHILD_OF}'],
		[
			'{"timeUnixNano":1657275433246142000,"attributeMap":{"event":"HTTP request received S3","level":"info","method":"GET","url":"/dispatch?customer=392\\u0026nonse=0.015296363321630757"}}',
		],
		false,
	],
];
