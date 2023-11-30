import { EQueryType } from 'types/common/dashboard';

import { Tag } from '../styles';

function QueryTypeTag({ queryType }: IQueryTypeTagProps): JSX.Element {
	switch (queryType) {
		case EQueryType.QUERY_BUILDER:
			return (
				<span>
					<Tag color="geekblue">Query Builder</Tag>
				</span>
			);

		case EQueryType.CLICKHOUSE:
			return (
				<span>
					<Tag color="orange">ClickHouse Query</Tag>
				</span>
			);
		case EQueryType.PROM:
			return (
				<span>
					<Tag color="green">PromQL</Tag>
				</span>
			);
		default:
			return <span />;
	}
}

interface IQueryTypeTagProps {
	queryType?: EQueryType;
}

QueryTypeTag.defaultProps = {
	queryType: EQueryType.QUERY_BUILDER,
};

export default QueryTypeTag;
