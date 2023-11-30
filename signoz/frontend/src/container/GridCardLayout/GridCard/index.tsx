import { PANEL_TYPES } from 'constants/queryBuilder';
import { useGetQueryRange } from 'hooks/queryBuilder/useGetQueryRange';
import { useStepInterval } from 'hooks/queryBuilder/useStepInterval';
import { useIsDarkMode } from 'hooks/useDarkMode';
import { useResizeObserver } from 'hooks/useDimensions';
import { useIntersectionObserver } from 'hooks/useIntersectionObserver';
import { getDashboardVariables } from 'lib/dashbaordVariables/getDashboardVariables';
import { getUPlotChartOptions } from 'lib/uPlotLib/getUplotChartOptions';
import { getUPlotChartData } from 'lib/uPlotLib/utils/getUplotChartData';
import isEmpty from 'lodash-es/isEmpty';
import _noop from 'lodash-es/noop';
import { memo, useCallback, useMemo, useRef, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { UpdateTimeInterval } from 'store/actions';
import { AppState } from 'store/reducers';
import { GlobalReducer } from 'types/reducer/globalTime';

import EmptyWidget from '../EmptyWidget';
import { MenuItemKeys } from '../WidgetHeader/contants';
import { GridCardGraphProps } from './types';
import WidgetGraphComponent from './WidgetGraphComponent';

function GridCardGraph({
	widget,
	name,
	onClickHandler = _noop,
	headerMenuList = [MenuItemKeys.View],
	isQueryEnabled,
	threshold,
	variables,
	fillSpans = false,
}: GridCardGraphProps): JSX.Element {
	const dispatch = useDispatch();
	const [errorMessage, setErrorMessage] = useState<string>();

	const onDragSelect = useCallback(
		(start: number, end: number): void => {
			const startTimestamp = Math.trunc(start);
			const endTimestamp = Math.trunc(end);

			if (startTimestamp !== endTimestamp) {
				dispatch(UpdateTimeInterval('custom', [startTimestamp, endTimestamp]));
			}
		},
		[dispatch],
	);

	const graphRef = useRef<HTMLDivElement>(null);

	const isVisible = useIntersectionObserver(graphRef, undefined, true);

	const { minTime, maxTime, selectedTime: globalSelectedInterval } = useSelector<
		AppState,
		GlobalReducer
	>((state) => state.globalTime);

	const updatedQuery = useStepInterval(widget?.query);

	const isEmptyWidget =
		widget?.id === PANEL_TYPES.EMPTY_WIDGET || isEmpty(widget);

	const queryResponse = useGetQueryRange(
		{
			selectedTime: widget?.timePreferance,
			graphType: widget?.panelTypes,
			query: updatedQuery,
			globalSelectedInterval,
			variables: getDashboardVariables(variables),
		},
		{
			queryKey: [
				maxTime,
				minTime,
				globalSelectedInterval,
				variables,
				widget?.query,
				widget?.panelTypes,
				widget.timePreferance,
			],
			keepPreviousData: true,
			enabled: isVisible && !isEmptyWidget && isQueryEnabled,
			refetchOnMount: false,
			onError: (error) => {
				setErrorMessage(error.message);
			},
		},
	);

	const isEmptyLayout = widget?.id === PANEL_TYPES.EMPTY_WIDGET;

	const containerDimensions = useResizeObserver(graphRef);

	const chartData = getUPlotChartData(queryResponse?.data?.payload, fillSpans);

	const isDarkMode = useIsDarkMode();

	const menuList =
		widget.panelTypes === PANEL_TYPES.TABLE
			? headerMenuList.filter((menu) => menu !== MenuItemKeys.CreateAlerts)
			: headerMenuList;

	const options = useMemo(
		() =>
			getUPlotChartOptions({
				id: widget?.id,
				apiResponse: queryResponse.data?.payload,
				dimensions: containerDimensions,
				isDarkMode,
				onDragSelect,
				yAxisUnit: widget?.yAxisUnit,
				onClickHandler,
				thresholds: widget.thresholds,
			}),
		[
			widget?.id,
			widget?.yAxisUnit,
			widget.thresholds,
			queryResponse.data?.payload,
			containerDimensions,
			isDarkMode,
			onDragSelect,
			onClickHandler,
		],
	);

	return (
		<div style={{ height: '100%', width: '100%' }} ref={graphRef}>
			{isEmptyLayout ? (
				<EmptyWidget />
			) : (
				<WidgetGraphComponent
					data={chartData}
					options={options}
					widget={widget}
					queryResponse={queryResponse}
					errorMessage={errorMessage}
					isWarning={false}
					name={name}
					onDragSelect={onDragSelect}
					threshold={threshold}
					headerMenuList={menuList}
					onClickHandler={onClickHandler}
				/>
			)}
		</div>
	);
}

GridCardGraph.defaultProps = {
	onDragSelect: undefined,
	onClickHandler: undefined,
	isQueryEnabled: true,
	threshold: undefined,
	headerMenuList: [MenuItemKeys.View],
};

export default memo(GridCardGraph);
