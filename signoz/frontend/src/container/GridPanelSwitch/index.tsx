import { ToggleGraphProps } from 'components/Graph/types';
import { PANEL_TYPES_COMPONENT_MAP } from 'constants/panelTypes';
import { PANEL_TYPES } from 'constants/queryBuilder';
import { GRID_TABLE_CONFIG } from 'container/GridTableComponent/config';
import { FC, forwardRef, memo, useMemo } from 'react';

import { GridPanelSwitchProps, PropsTypePropsMap } from './types';

const GridPanelSwitch = forwardRef<
	ToggleGraphProps | undefined,
	GridPanelSwitchProps
>(
	(
		{ panelType, data, yAxisUnit, panelData, query, options, thresholds },
		ref,
	): JSX.Element | null => {
		const currentProps: PropsTypePropsMap = useMemo(() => {
			const result: PropsTypePropsMap = {
				[PANEL_TYPES.TIME_SERIES]: {
					data,
					options,
					ref,
				},
				[PANEL_TYPES.VALUE]: {
					data,
					yAxisUnit,
					thresholds,
				},
				[PANEL_TYPES.TABLE]: {
					...GRID_TABLE_CONFIG,
					data: panelData,
					query,
					thresholds,
				},
				[PANEL_TYPES.LIST]: null,
				[PANEL_TYPES.TRACE]: null,
				[PANEL_TYPES.EMPTY_WIDGET]: null,
			};

			return result;
		}, [data, options, ref, yAxisUnit, thresholds, panelData, query]);

		const Component = PANEL_TYPES_COMPONENT_MAP[panelType] as FC<
			PropsTypePropsMap[typeof panelType]
		>;
		const componentProps = useMemo(() => currentProps[panelType], [
			panelType,
			currentProps,
		]);

		if (!Component || !componentProps) return null;
		// eslint-disable-next-line react/jsx-props-no-spreading
		return <Component {...componentProps} />;
	},
);

GridPanelSwitch.displayName = 'GridPanelSwitch';

export default memo(GridPanelSwitch);
