import { Typography } from 'antd';
import axios from 'axios';
import { SOMETHING_WENT_WRONG } from 'constants/api';
import Graph from 'container/GridCardLayout/GridCard';
import { Card, GraphContainer } from 'container/MetricsApplication/styles';
import { OnClickPluginOpts } from 'lib/uPlotLib/plugins/onClickPlugin';
import { Widgets } from 'types/api/dashboard/getAll';

function TopLevelOperation({
	name,
	opName,
	topLevelOperationsIsError,
	topLevelOperationsError,
	onDragSelect,
	handleGraphClick,
	widget,
	topLevelOperationsIsLoading,
}: TopLevelOperationProps): JSX.Element {
	return (
		<Card data-testid={name}>
			{topLevelOperationsIsError ? (
				<Typography>
					{axios.isAxiosError(topLevelOperationsError)
						? topLevelOperationsError.response?.data
						: SOMETHING_WENT_WRONG}
				</Typography>
			) : (
				<GraphContainer>
					<Graph
						fillSpans={false}
						name={name}
						widget={widget}
						onClickHandler={handleGraphClick(opName)}
						onDragSelect={onDragSelect}
						isQueryEnabled={!topLevelOperationsIsLoading}
					/>
				</GraphContainer>
			)}
		</Card>
	);
}

interface TopLevelOperationProps {
	name: string;
	opName: string;
	topLevelOperationsIsError: boolean;
	topLevelOperationsError: unknown;
	onDragSelect: (start: number, end: number) => void;
	handleGraphClick: (type: string) => OnClickPluginOpts['onClick'];
	widget: Widgets;
	topLevelOperationsIsLoading: boolean;
}

export default TopLevelOperation;
