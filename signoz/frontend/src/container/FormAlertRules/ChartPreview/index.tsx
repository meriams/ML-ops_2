import { InfoCircleOutlined } from '@ant-design/icons';
import Spinner from 'components/Spinner';
import { initialQueriesMap, PANEL_TYPES } from 'constants/queryBuilder';
import GridPanelSwitch from 'container/GridPanelSwitch';
import { timePreferenceType } from 'container/NewWidget/RightContainer/timeItems';
import { Time } from 'container/TopNav/DateTimeSelection/config';
import { useGetQueryRange } from 'hooks/queryBuilder/useGetQueryRange';
import { useIsDarkMode } from 'hooks/useDarkMode';
import { useResizeObserver } from 'hooks/useDimensions';
import { getUPlotChartOptions } from 'lib/uPlotLib/getUplotChartOptions';
import { getUPlotChartData } from 'lib/uPlotLib/utils/getUplotChartData';
import { useMemo, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import { useSelector } from 'react-redux';
import { AppState } from 'store/reducers';
import { AlertDef } from 'types/api/alerts/def';
import { Query } from 'types/api/queryBuilder/queryBuilderData';
import { EQueryType } from 'types/common/dashboard';
import { GlobalReducer } from 'types/reducer/globalTime';

import { ChartContainer, FailedMessageContainer } from './styles';
import { covertIntoDataFormats } from './utils';

export interface ChartPreviewProps {
	name: string;
	query: Query | null;
	graphType?: PANEL_TYPES;
	selectedTime?: timePreferenceType;
	selectedInterval?: Time;
	headline?: JSX.Element;
	alertDef?: AlertDef;
	userQueryKey?: string;
	allowSelectedIntervalForStepGen?: boolean;
}

function ChartPreview({
	name,
	query,
	graphType = PANEL_TYPES.TIME_SERIES,
	selectedTime = 'GLOBAL_TIME',
	selectedInterval = '5min',
	headline,
	userQueryKey,
	allowSelectedIntervalForStepGen = false,
	alertDef,
}: ChartPreviewProps): JSX.Element | null {
	const { t } = useTranslation('alerts');
	const threshold = alertDef?.condition.target || 0;
	const { minTime, maxTime } = useSelector<AppState, GlobalReducer>(
		(state) => state.globalTime,
	);

	const thresholdValue = covertIntoDataFormats({
		value: threshold,
		sourceUnit: alertDef?.condition.targetUnit,
		targetUnit: query?.unit,
	});

	const canQuery = useMemo((): boolean => {
		if (!query || query == null) {
			return false;
		}

		switch (query?.queryType) {
			case EQueryType.PROM:
				return query.promql?.length > 0 && query.promql[0].query !== '';
			case EQueryType.CLICKHOUSE:
				return (
					query.clickhouse_sql?.length > 0 &&
					query.clickhouse_sql[0].query?.length > 0
				);
			case EQueryType.QUERY_BUILDER:
				return (
					query.builder.queryData.length > 0 &&
					query.builder.queryData[0].queryName !== ''
				);
			default:
				return false;
		}
	}, [query]);

	const queryResponse = useGetQueryRange(
		{
			query: query || initialQueriesMap.metrics,
			globalSelectedInterval: selectedInterval,
			graphType,
			selectedTime,
			params: {
				allowSelectedIntervalForStepGen,
			},
		},
		{
			queryKey: [
				'chartPreview',
				userQueryKey || JSON.stringify(query),
				selectedInterval,
				minTime,
				maxTime,
			],
			retry: false,
			enabled: canQuery,
		},
	);

	const graphRef = useRef<HTMLDivElement>(null);

	const chartData = getUPlotChartData(queryResponse?.data?.payload);

	const containerDimensions = useResizeObserver(graphRef);

	const isDarkMode = useIsDarkMode();

	const options = useMemo(
		() =>
			getUPlotChartOptions({
				id: 'alert_legend_widget',
				yAxisUnit: query?.unit,
				apiResponse: queryResponse?.data?.payload,
				dimensions: containerDimensions,
				isDarkMode,
				thresholds: [
					{
						index: '0', // no impact
						keyIndex: 0,
						moveThreshold: (): void => {},
						selectedGraph: PANEL_TYPES.TIME_SERIES, // no impact
						thresholdValue,
						thresholdLabel: `${t(
							'preview_chart_threshold_label',
						)} (y=${thresholdValue} ${query?.unit || ''})`,
					},
				],
			}),
		[
			query?.unit,
			queryResponse?.data?.payload,
			containerDimensions,
			isDarkMode,
			t,
			thresholdValue,
		],
	);

	return (
		<ChartContainer>
			{headline}
			{(queryResponse?.isError || queryResponse?.error) && (
				<FailedMessageContainer color="red" title="Failed to refresh the chart">
					<InfoCircleOutlined />{' '}
					{queryResponse.error.message || t('preview_chart_unexpected_error')}
				</FailedMessageContainer>
			)}
			{queryResponse.isLoading && (
				<Spinner size="large" tip="Loading..." height="70vh" />
			)}
			{chartData && !queryResponse.isError && (
				<div ref={graphRef} style={{ height: '100%' }}>
					<GridPanelSwitch
						options={options}
						panelType={graphType}
						data={chartData}
						name={name || 'Chart Preview'}
						panelData={queryResponse.data?.payload.data.newResult.data.result || []}
						query={query || initialQueriesMap.metrics}
						yAxisUnit={query?.unit}
					/>
				</div>
			)}
		</ChartContainer>
	);
}

ChartPreview.defaultProps = {
	graphType: PANEL_TYPES.TIME_SERIES,
	selectedTime: 'GLOBAL_TIME',
	selectedInterval: '5min',
	headline: undefined,
	userQueryKey: '',
	allowSelectedIntervalForStepGen: false,
	alertDef: undefined,
};

export default ChartPreview;
