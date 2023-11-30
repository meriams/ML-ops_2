import { ILog } from 'types/api/logs/log';
import { IBuilderQuery } from 'types/api/queryBuilder/queryBuilderData';

export type LogsExplorerListProps = {
	isLoading: boolean;
	currentStagedQueryData: IBuilderQuery | null;
	logs: ILog[];
	onEndReached: (index: number) => void;
};
